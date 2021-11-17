package server

import (
	"reflect"

	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/log"
)

type Server struct {
	worker *Worker
}

var svr *Server

func NewServer(opts ...WorkerOpt) *Server {
	svr = &Server{}
	svr.worker = NewWorker()
	for _, opt := range opts {
		opt(svr.worker)
	}

	return svr
}

func (s *Server) Reg(funcName string, wFunc interface{}) {
	log.Info("reg func: %s", funcName)
	t := reflect.TypeOf(wFunc).Kind().String()
	if t == "func" {
		s.worker.addFuncWorker(funcName, wFunc)
	} else {
		panic("must be func")
	}
}

//启动
func (s *Server) Run(numWorkers int) (err error) {
	if err = s.worker.broker.Activate(); err != nil {
		return
	}
	if err = s.worker.backend.Activate(); err != nil {
		return
	}

	s.worker.Run(numWorkers)

	return err
}

func (s *Server) ShutDown() {
	s.worker.Stop()
}

func (s *Server) CloneBroker() broker.Broker {

	return s.worker.broker.Clone()

}

func (s *Server) CloneBackend() backend.Backend {

	return s.worker.backend.Clone()
}

func GetServer() *Server {
	if svr != nil {
		return svr
	}

	return nil
}
