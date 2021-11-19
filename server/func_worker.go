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

func (f *FuncWorker) toValue(val interface{}, Type int) reflect.Value {
	var inter interface{}
	switch Type {
	case message.Bool:
		inter = cast.ToBool(val)
	case message.Int:
		inter = cast.ToInt(val)
	case message.Int8:
		inter = cast.ToInt8(val)
	case message.Int16:
		inter = cast.ToInt16(val)
	case message.Int32:
		inter = cast.ToInt32(val)
	case message.Int64:
		inter = cast.ToInt64(val)
	case message.Uint:
		inter = cast.ToUint(val)
	case message.Uint8:
		inter = cast.ToUint8(val)
	case message.Uint16:
		inter = cast.ToUint16(val)
	case message.Uint32:
		inter = cast.ToUint32(val)
	case message.Uint64:
		inter = cast.ToUint64(val)
	case message.Float32:
		inter = cast.ToFloat32(val)
	case message.Float64:
		inter = cast.ToFloat32(val)
	case message.String:
		inter = cast.ToString(val)
	}

	return reflect.ValueOf(inter)
}

func (f *FuncWorker) execFunc(args ...message.Signature) (vals []string, err error) {
	fc := reflect.ValueOf(f.wFunc)

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = f.toValue(arg.Val, arg.Type)
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
