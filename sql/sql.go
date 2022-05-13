/**
 * @Author: guobob
 * @Description:
 * @File:  sql.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:49
 */

package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/generate_data/util"
	"github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
	//"os"
)

//Store prepare statement and handle
type statement struct {
	query  string
	handle *sql.Stmt
}

type SQLHandle struct {
	Dsn    string
	cfg    *mysql.Config
	schema string
	pool   *sql.DB
	conn   *sql.Conn
	stmts  map[uint64]statement
	Ctx    context.Context
	Sqlch  chan string
	SqlRes [][]string
	//Log    *zap.Logger
}

func NewSQLHandle(dsn string, cfg *mysql.Config) *SQLHandle {
	return &SQLHandle{
		Ctx: context.Background(),
		Dsn: dsn,
		//Log:    log,
		cfg:    cfg,
		SqlRes: make([][]string, 0),
		stmts:  make(map[uint64]statement),
	}
}

func (s *SQLHandle) open(schema string) (*sql.DB, error) {
	cfg := s.cfg
	if len(schema) > 0 && cfg.DBName != schema {
		cfg = cfg.Clone()
		cfg.DBName = schema
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

//Handle Handshake messages, similar to Use Database
func (s *SQLHandle) handshake(ctx context.Context, schema string) error {
	pool, err := s.open(schema)
	if err != nil {
		return err
	}
	s.pool = pool
	s.schema = schema
	_, err = s.getConn(ctx)
	return err
}

// Conn returns a single connection by either opening a new connection
// or returning an existing connection from the connection pool. Conn will
// block until either a connection is returned or ctx is canceled.
// Queries run on the same Conn will be run in the same database session.
//
// Every Conn must be returned to the database pool after use by
// calling Conn.Close.
func (s *SQLHandle) getConn(ctx context.Context) (*sql.Conn, error) {
	var err error
	if s.pool == nil {
		s.pool, err = s.open(s.schema)
		if err != nil {
			return nil, err
		}
	}
	if s.conn == nil {
		s.conn, err = s.pool.Conn(ctx)
		if err != nil {
			return nil, err
		}
	}
	return s.conn, nil
}

//Disconnect from replay server
func (s *SQLHandle) quit(reconnect bool) {
	for id, stmt := range s.stmts {
		if stmt.handle != nil {
			if err := stmt.handle.Close(); err != nil {
				fmt.Println("close stmt.handle fail ," + err.Error())
			}
			stmt.handle = nil
		}
		if reconnect {
			s.stmts[id] = stmt
		} else {
			delete(s.stmts, id)
		}
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			fmt.Println("close conn fail ," + err.Error())
		}
		s.conn = nil
	}
	if s.pool != nil {
		if err := s.pool.Close(); err != nil {
			fmt.Println("close pool fail ," + err.Error())
		}
		s.pool = nil
	}
}

//Execute SQL on replay Server
func (s *SQLHandle) execute(ctx context.Context, query string) error {
	conn, err := s.getConn(ctx)
	if err != nil {
		return err
	}

	rows, err := conn.QueryContext(ctx, query)
	defer func() {
		if rows != nil {
			if rs := rows.Close(); rs != nil {
				fmt.Println("close row fail," + rs.Error())
			}
		}
	}()
	if err != nil {
		return err
	}

	for rows.Next() {
		s.ReadRowValues(rows)
	}
	return nil
}

//Exec prepare statment on replay sql
func (s *SQLHandle) stmtPrepare(ctx context.Context, id uint64, query string) error {
	stmt := s.stmts[id]
	stmt.query = query
	if stmt.handle != nil {
		if err := stmt.handle.Close(); err != nil {
			fmt.Println("close stmt handle fail ," + err.Error())
		}
		stmt.handle = nil
	}
	delete(s.stmts, id)
	conn, err := s.getConn(ctx)
	if err != nil {
		return err
	}
	//stats.Add(stats.StmtPrepares, 1)
	stmt.handle, err = conn.PrepareContext(ctx, stmt.query)
	if err != nil {
		//stats.Add(stats.FailedStmtPrepares, 1)
		return err
	}
	s.stmts[id] = stmt
	fmt.Println(fmt.Sprintf("%v id is %v", query, id))
	return nil
}

//Exec prepare on replay server
func (s *SQLHandle) stmtExecute(ctx context.Context, id uint64, params []interface{}) error {
	stmt, err := s.getStmt(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println(params)
	rows, err := stmt.QueryContext(ctx, params...)
	defer func() {
		if rows != nil {
			if err := rows.Close(); err != nil {
				fmt.Println("close rows fail," + err.Error())
			}
		}
	}()

	if err != nil {
		return err
	}

	for rows.Next() {
		s.ReadRowValues(rows)
	}
	return nil
}

//Close prepare handle
func (s *SQLHandle) stmtClose(id uint64) {
	stmt, ok := s.stmts[id]
	if !ok {
		return
	}
	if stmt.handle != nil {
		if err := stmt.handle.Close(); err != nil {
			fmt.Println("close stmt handle fail," + err.Error())
		}
		stmt.handle = nil
	}
	delete(s.stmts, id)
}

//Get prepare handle ID
func (s *SQLHandle) getStmt(ctx context.Context, id uint64) (*sql.Stmt, error) {
	stmt, ok := s.stmts[id]
	if ok && stmt.handle != nil {
		return stmt.handle, nil
	} else if !ok {
		return nil, errors.Errorf("no such statement #%d", id)
	}
	conn, err := s.getConn(ctx)
	if err != nil {
		return nil, err
	}
	stmt.handle, err = conn.PrepareContext(ctx, stmt.query)
	if err != nil {
		return nil, err
	}
	s.stmts[id] = stmt
	return stmt.handle, nil
}

//var iiiii =0
func (s *SQLHandle) ReadRowValues(f *sql.Rows) {
	//Get the lastcols value from the sql.Rows
	//structure using unsafe and reflection mechanisms
	//and load it into the cache

	rs := reflect.ValueOf(f)
	foo := rs.Elem().FieldByName("lastcols")
	rf := foo
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	z := rf.Interface().([]driver.Value)
	rr := make([]string, 0, len(z))
	for i := range z {
		if z[i] == nil {
			rr = append(rr, "")
			continue
		}
		var v string
		err := util.ConvertAssign(&v, z[i])
		if err != nil {
			fmt.Println("get result fail ", err)
			os.Exit(1)
		}
		rr = append(rr, v)
	}
	s.SqlRes = append(s.SqlRes, rr)
}

func (s *SQLHandle) HandShake(schema string) error {
	return s.handshake(s.Ctx, schema)
}

func (s *SQLHandle) StmtPrepare(id uint64, query string) error {
	return s.stmtPrepare(s.Ctx, id, query)
}

func (s *SQLHandle) StmtExecute(id uint64, params []interface{}) error {
	return s.stmtExecute(s.Ctx, id, params)
}

func (s *SQLHandle) Quit(reconnect bool) {
	s.quit(reconnect)
}

func (s *SQLHandle) Execute(query string) error {
	return s.execute(s.Ctx, query)
}
