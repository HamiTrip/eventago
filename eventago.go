package eventago

import (
	"reflect"
	"strings"
	"fmt"
)


//Eventa interface
type Eventa interface {
	On(string, interface{})
	Fire(string, ...interface{}) [][]interface{}
	GoFire(string, ...interface{})
	HaveListener(string) bool
	CountListeners(string) int
}

type tEventa struct {
	handlers map[string][]interface{}
}

var container *tEventa

func init() {
	if container == nil {
		container = new(tEventa)
		container.handlers = make(map[string][]interface{})
	}
}

/*
On gives a Handler for specific Event and adds it to listening bus.
 Set multiple handlers for one event is supported.
 */
func On(event string, handler interface{}) {
	event = strings.ToLower(event)
	if reflect.ValueOf(handler).Kind() != reflect.Func {
		panic("handler is not func!")
	}
	if _, ok := container.handlers[event]; ok {
		container.handlers[event] = append(container.handlers[event], handler)
	} else {
		container.handlers[event] = []interface{}{handler}
	}
	fmt.Println("EVENTA: A handler added for:", event)
}

/*
Fire will run all of given event's handlers with parameters concurrent and collect results of them and return as an [][]interface{}
 */
func Fire(event string, params ...interface{}) (res [][]interface{}) {
	event = strings.ToLower(event)
	if handlers, ok := container.handlers[event]; ok {
		var resultsChan = make(chan []interface{}, len(container.handlers))
		var rParams = toR(params)
		for _, handler := range handlers {
			go runner(resultsChan, handler, rParams...)
		}
		for i := 0; i < len(container.handlers[event]); i++ {
			select {
			case result := <-resultsChan:
				res = append(res, result)
				if len(res) >= len(container.handlers[event]) {
					break
				}
			}
		}
	}
	return
}

/*
GoFire will fire events with full concurrency, But with this method can't receive handler's result!
 */
func GoFire(event string, params ...interface{}) {
	event = strings.ToLower(event)
	go Fire(event, params...)
}

func runner(result chan <- []interface{}, handler interface{}, rParams ...reflect.Value) {
	var res []interface{}
	var rRes = reflect.ValueOf(handler).Call(rParams)
	for _, value := range rRes {
		res = append(res, value.Interface())
	}
	result <- res
}

/*
HaveListener checks if given event have any listener
 */
func HaveListener(event string) bool {
	event = strings.ToLower(event)
	if handlers, ok := container.handlers[event]; ok {
		return len(handlers) > 0
	}
	return false
}

/*
CountListeners returns count of handlers are listening for given event
 */
func CountListeners(event string) int {
	event = strings.ToLower(event)
	if handlers, ok := container.handlers[event]; ok {
		return len(handlers)
	}
	return 0
}

func toR(i []interface{}) (r []reflect.Value) {
	for _, item := range i {
		r = append(r, reflect.ValueOf(item))
	}
	return
}