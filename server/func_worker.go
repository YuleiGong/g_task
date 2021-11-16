package server

import (
	"reflect"

	"github.com/YuleiGong/g_task/message"
	werrors "github.com/pkg/errors"
	"github.com/spf13/cast"
)

type FuncWorker struct {
	name  string
	wFunc interface{}
}

func NewFuncWorker(name string, wFunc interface{}) FuncWorker {
	return FuncWorker{name: name, wFunc: wFunc}

}

func (f *FuncWorker) execFunc(args ...message.Arg) (vals []string, err error) {
	fc := reflect.ValueOf(f.wFunc)

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		var inter interface{}
		switch arg.Type {
		case "int":
			inter = cast.ToInt(arg.Val)
		case "int64":
			inter = cast.ToInt64(arg.Val)
		case "string":
			inter = cast.ToString(arg.Val)
		}
		in[i] = reflect.ValueOf(inter)
	}

	val := fc.Call(in)
	length := len(val)

	if length == 0 { //无返回值
		err = werrors.WithMessagef(ErrFunc, "func %s must return error", f.name)
		return
	}

	errInter := val[length-1].Interface()
	if errInter != nil {
		var ok bool
		err, ok = errInter.(error)
		if !ok { //至少需要一个error返回
			err = werrors.WithMessagef(ErrFunc, "func %s must return include error", f.name)
		}
		return
	}

	if length == 1 {
		return
	}

	for i := 0; i < length-1; i++ {
		vals = append(vals, cast.ToString(val[i].Interface()))
	}

	return vals, err
}
