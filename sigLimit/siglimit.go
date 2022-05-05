package sigLimit

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
)

type sigLimit struct {
	wait     *sync.WaitGroup
	weighted *semaphore.Weighted
}

func NewSigLimit(max int) *sigLimit {
	return &sigLimit{
		wait:     &sync.WaitGroup{},
		weighted: semaphore.NewWeighted(int64(max)),
	}
}

func (s *sigLimit) Add() {
	err := s.weighted.Acquire(context.Background(), 1)
	if err != nil {
		fmt.Println("添加任务错误: ", err.Error())
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
