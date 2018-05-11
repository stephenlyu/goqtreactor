package reactor

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt"
	"fmt"
	"sync/atomic"
)

var event_type = core.QEvent_RegisterEventType(-1)
var __invoker = createInvoker()
var _pool Pool
var initialized uint32

type Invoker struct {
	core.QObject
}

func createInvoker() *Invoker {
	ret := NewInvoker(nil)
	ret.ConnectCustomEvent(ret.onCustomEvent)
	return ret
}

func (this *Invoker) onCustomEvent(e *core.QEvent) {
	ptr, ok := qt.Receive(e.Pointer())
	if !ok {
		fmt.Println("bad")
		return
	}

	ce, ok := ptr.(*CallbackEvent)
	if !ok {
		fmt.Println("Not CallbackEvent")
		return
	}

	ce.call()
}


type CallbackEvent struct {
	core.QEvent

	f func(...interface{})
	args []interface{}
}

func NewCallbackEvent(f func(...interface{}), args []interface{}) *CallbackEvent {
	e := core.NewQEvent(core.QEvent__Type(event_type))
	ret := &CallbackEvent{
		QEvent: *e,
		f: f,
		args: args,
	}

	qt.Register(e.Pointer(), ret)
	ret.ConnectDestroyQEvent(func () {
		qt.Unregister(e.Pointer())
	})

	return ret
}

func (this *CallbackEvent) call() {
	this.f(this.args...)
}

type funcTask struct {
	f func(...interface{})
	args []interface{}
}

func (this *funcTask) Do() {
	this.f(this.args...)
}

func CallFromThread(f func(...interface{}), args ...interface{}) {
	if atomic.LoadUint32(&initialized) == 0 {
		panic("reactor uninitialized")
	}
	e := NewCallbackEvent(f, args)
	core.QCoreApplication_PostEvent(__invoker, e, int(core.Qt__NormalEventPriority))
}

func CallInThread(f func(...interface{}), args ...interface{}) {
	if atomic.LoadUint32(&initialized) == 0 {
		panic("reactor uninitialized")
	}

	_pool.PostTask(&funcTask{f: f, args: args})
}

func Initialize() {
	_pool = NewPool(0)
	_pool.Start()
	atomic.StoreUint32(&initialized, 1)
}

func Destroy() {
	_pool.Stop()
	_pool = nil
	atomic.StoreUint32(&initialized, 0)
}
