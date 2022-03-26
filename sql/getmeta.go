/**
 * @Author: guobob
 * @Description:
 * @File:  getmeta.go
 * @Version: 1.0.0
 * @Date: 2022/3/24 22:57
 */

package sql

//get meta data stmt id set 0
var (
	getMetaSql = "select table_name,column_name,ORDINAL_POSITION,DATA_TYPE," +
		"CHARACTER_MAXIMUM_LENGTH ,CHARACTER_OCTET_LENGTH ,NUMERIC_PRECISION,NUMERIC_SCALE," +
		"DATETIME_PRECISION from information_schema.COLUMNS where   TABLE_SCHEMA =? and table_name =? ;"
	getMetaStmtID uint64 = 0
)

func GetColumnInfo(mysql *SQLHandle, databaseName, tableName string) error {
	if _, ok := mysql.stmts[getMetaStmtID]; !ok {
		err := mysql.StmtPrepare(getMetaStmtID, getMetaSql)
		if err != nil {
			mysql.Log.Error("prepare sql fail ," + err.Error())
			return err
		}
	}
	params := make([]interface{}, 0)
	params = append(params, databaseName, tableName)
	err := mysql.StmtExecute(getMetaStmtID, params)
	if err != nil {
		return err
	}

	return nil
}
