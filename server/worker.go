package server

import (
	"context"
	"errors"

	"time"

	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/message"
	"github.com/go-redis/redis"
)

type Worker struct {
	funcWorkerMap map[string]FuncWorker
	readyChan     chan struct{}
	resultChan    chan *message.MessageResult
	stopChan      chan error

	broker  broker.Broker
	backend backend.Backend
}

type WorkerOpt func(w *Worker)

func WithBroker(broker broker.Broker) WorkerOpt {
	return func(s *Worker) {
		s.broker = broker
	}
}

func WithBackend(backend backend.Backend) WorkerOpt {
	return func(s *Worker) {
		s.backend = backend
	}
}

func NewWorker() *Worker {
	return &Worker{
		funcWorkerMap: make(map[string]FuncWorker),
	}
}

func (w *Worker) addFuncWorker(name string, wf interface{}) {
	w.funcWorkerMap[name] = NewFuncWorker(name, wf)
}

func (w *Worker) Run(numWorker int) error {
	w.stopChan = make(chan error)
	w.initReadyWoker(numWorker)
	w.resultChan = make(chan *message.MessageResult, numWorker)
	go w.resultSchedule()
	go w.wokerSchedule()

	return <-w.stopChan
}

func (w *Worker) Stop() {
	close(w.readyChan)
	for len(w.resultChan) > 0 {
		time.Sleep(time.Millisecond)
	}
	close(w.resultChan)
	w.stopChan <- errors.New("stop")
}

func (w *Worker) initReadyWoker(numWorker int) {
	w.readyChan = make(chan struct{}, numWorker)
	for numWorker > 0 {
		w.readyChan <- struct{}{}
		numWorker--
	}
}

func (w *Worker) wokerSchedule() {
	for range w.readyChan {
		taskID, msg, err := w.broker.Pop()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				w.resultChan <- &message.MessageResult{}
				continue
			}
			w.Stop()
		}
		w.execFuncWorker(taskID, msg)
	}
}

func (w *Worker) execFuncWorker(taskID string, msg *message.Message) {
	var (
		err    error
		result []string
	)
	msgRes := &message.MessageResult{ErrCode: message.SUCCESS, TaskID: taskID}
	defer func() {
		if e := recover(); e != nil {
			w.Stop()
		}
		if err != nil {
			msgRes.ErrCode = message.FAILURE
			msgRes.ErrMsg = err.Error()
			msg.Failure()
		} else {
			msgRes.Val = result
			msg.Success()
		}
		w.updateBroker(msg, taskID)
		w.resultChan <- msgRes
	}()
	msg.Started()
	w.updateBroker(msg, taskID) //任务开始

	if msg.IsTimeoutOpt() {
		result, err = w.execFuncWithTimeout(msg)
	} else {
		result, err = w.execFunc(msg)
	}
	if errors.Is(err, ErrTimeout) && msg.IsRetry() {
		w.RetryTask(msg)
	}

}

type funcResp struct {
	result []string
	err    error
}

func (w *Worker) execFuncWithTimeout(msg *message.Message) (result []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), msg.Timeout)
	defer cancel()

	resp := make(chan funcResp, 1)
	go func(m *message.Message) {
		fr := funcResp{}
		fr.result, fr.err = w.execFunc(m)
		resp <- fr
	}(msg)
	defer close(resp)

	select {
	case <-ctx.Done():
		err = ErrTimeout
	case res := <-resp:
		result = res.result
		err = res.err
	}

	return result, err
}

func (w *Worker) execFunc(msg *message.Message) (result []string, err error) {
	funcWorker := w.funcWorkerMap[msg.FuncName]
	return funcWorker.execFunc(msg.Args...)
}

func (w *Worker) resultSchedule() {
	for m := range w.resultChan {
		if m.TaskID != "" {
			if err := w.backend.SetResult(m.TaskID, m); err != nil {
				w.Stop()
			}
		}
		w.readyChan <- struct{}{}
	}
}

func (w *Worker) updateBroker(msg *message.Message, taskID string) (err error) {
	if taskID == "" {
		return
	}
	if err = w.broker.Set(taskID, msg); err != nil {
		w.Stop()
	}
	return err
}

func (w *Worker) RetryTask(msg *message.Message) (err error) {
	msg.Retry()
	msg.AddRetry()
	if err = w.broker.Push(msg.TaskID, msg); err != nil {
		w.Stop()
		return
	}

	return nil

}
