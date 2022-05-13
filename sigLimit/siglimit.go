package sigLimit

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"gotest.tools/gotestsum/log"
	//"gotest.tools/gotestsum/log"
	//"gotest.tools/gotestsum/log"
)

type sigLimit struct {
	wait     *sync.WaitGroup
	weighted *semaphore.Weighted
	log      *zap.Logger
}

func NewSigLimit(max int, log *zap.Logger) *sigLimit {
	return &sigLimit{
		wait:     &sync.WaitGroup{},
		weighted: semaphore.NewWeighted(int64(max)),
		log:      log,
	}
}

func (s *sigLimit) Add() {
	err := s.weighted.Acquire(context.Background(), 1)
	if err != nil {
		log.Errorf("添加任务错误: ", err.Error())
	}
	s.wait.Add(1)
}

func (s *sigLimit) Done() {
	s.weighted.Release(1)
	s.wait.Done()
}

func (s *sigLimit) Wait() {
	s.wait.Wait()
}
