package dispatch

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	cdp "github.com/chromedp/cdproto"
	"github.com/mailru/easyjson"

	"github.com/argots/debugger"
)

func New() debugger.Dispatcher {
	return dispatcher{}
}

type dispatcher struct{}

func (dispatcher) Dispatch(d debugger.Driver, msg *cdp.Message) {
	parts := strings.Split(string(msg.Method), ".") // domain.method
	v, err := call(d, parts[0], nil)
	if err == nil {
		// log.Println("Got receiver ", reflect.TypeOf(v), parts[1])
		v, err = call(v, parts[1], msg.Params)
	}
	if err == nil {
		if v == nil {
			v = map[string]string{}
		}
		b, err := json.Marshal(v)
		if err == nil {
			msg.Result = easyjson.RawMessage(b)
		}
	}
	if err != nil {
		msg.Error = &cdp.Error{Code: -1, Message: err.Error()}
	}
	log.Println("Responding", msg.ID, string(msg.Result), msg.Error)
}

func call(receiver interface{}, method string, arg easyjson.RawMessage) (interface{}, error) {
	if receiver == nil {
		return nil, nil
	}
	r := reflect.ValueOf(receiver)
	rType := r.Type()
	for kk := 0; kk < rType.NumMethod(); kk++ {
		m := rType.Method(kk)
		if strings.EqualFold(method, m.Name) {
			// log.Println("Found method", m.Name)
			return callMethod(m.Type, r.Method(kk), arg)
		}
	}
	return nil, nil
}

func callMethod(t reflect.Type, m reflect.Value, arg easyjson.RawMessage) (interface{}, error) {
	if len(arg) == 0 || t.NumIn() <= 1 {
		return callResult(m.Call(nil))
	}

	argVal := reflect.New(t.In(1).Elem())
	if err := json.Unmarshal([]byte(arg), argVal.Interface()); err != nil {
		log.Println("Unexpected unmarshal error", err)
		return nil, err
	}

	// log.Println("Unmarshaled", argVal.Type(), argVal.Interface())
	return callResult(m.Call([]reflect.Value{argVal}))
}

func callResult(v []reflect.Value) (interface{}, error) {
	switch len(v) {
	case 1:
		result := v[0].Interface()
		if err, ok := result.(error); ok {
			return nil, err
		}
		return result, nil
	case 2: //nolint: gomnd
		err, _ := v[1].Interface().(error)
		return v[0].Interface(), err
	case 0:
	default:
		log.Println("Unexpected number of results", v)
	}

	return nil, nil
}
