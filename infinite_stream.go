package infinite_stream_go

type Stream interface {
	Car() interface{}
	Cdr() Stream
	IsTail() bool
}

func Delay(fn func() Stream) func() Stream {
	already_run := false
	var ret Stream
	return func() Stream {
		if !already_run {
			ret = Force(fn)
			already_run = true
		}
		return ret
	}
}

func Force(f func() Stream) Stream {
	return f()
}

type ConsPair struct {
	car interface{}
	cdr func() Stream
}

func Cons(a interface{}, b func() Stream) (cp ConsPair) {
	cp.car = a
	cp.cdr = Delay(b)
	return cp
}

func (cp ConsPair) Car() interface{} {
	return cp.car
}

func (cp ConsPair) Cdr() Stream {
	return Force(cp.cdr)
}

func (cp ConsPair) IsTail() bool {
	return false
}

type Tail struct {
	// dummy
}

func (t Tail) Car() interface{} {
	return nil
}

func (t Tail) Cdr() Stream {
	return t
}

func (t Tail) IsTail() bool {
	return true
}

func TailCdr() func() Stream {
	return func() Stream {
		return Tail{}
	}
}

func Empty() Stream {
	return Tail{}
}

///////////////////////////////////////////////////////////////
// operators

func From(items ...interface{}) Stream {
	if len(items) == 0 {
		return Tail{}
	} else {
		return Cons(items[0], func() Stream {
			return From(items[1:]...)
		})
	}
}

func Ref(n int, s Stream) interface{} {
	if n == 0 {
		return s.Car()
	} else {
		return Ref(n-1, s.Cdr())
	}
}

func Take(n int, s Stream) Stream {
	if 0 < n {
		return Cons(s.Car(), func() Stream {
			return Take(n-1, s.Cdr())
		})
	} else {
		return Tail{}
	}
}

///////////////////////////////////////////////////////////////
// transducer

type Reducer func(interface{}, interface{}) interface{}
type Transducer func(Reducer) Reducer

func Reduce(fn Reducer, init interface{}, s Stream) interface{} {
	if s.IsTail() {
		return init
	}
	result := fn(init, s.Car())
	return Reduce(fn, result, s.Cdr())
}

func Conj(s Stream, item interface{}) Stream {
	if s.IsTail() {
		return Cons(item, TailCdr())
	}
	return Cons(s.Car(), func() Stream {
		return Conj(s.Cdr(), item)
	})
}

func ConjReducer(result interface{}, input interface{}) interface{} {
	if s, ok := result.(Stream); ok {
		return Conj(s, input)
	}
	return nil
}

func Map(fn func(interface{}) interface{}) Transducer {
	return func(reducing Reducer) Reducer {
		return func(result interface{}, input interface{}) interface{} {
			return reducing(result, fn(input))
		}
	}
}

func Filter(pred func(interface{}) bool) Transducer {
	return func(reducing Reducer) Reducer {
		return func(result interface{}, input interface{}) interface{} {
			if pred(input) {
				return reducing(result, input)
			} else {
				return result
			}
		}
	}
}

func Comp(transducer ...Transducer) Transducer {
	return func(reducing Reducer) Reducer {
		comp := reducing
		for i := len(transducer) - 1; i >= 0; i-- {
			comp = transducer[i](comp)
		}
		return comp
	}
}

func Transduce(tr Transducer, reducing Reducer, init interface{}, s Stream) interface{} {
	return Reduce(tr(reducing), init, s)
}

func Sequence(tr Transducer, s Stream) Stream {
	if s.IsTail() {
		return Tail{}
	}

	isChanged := false
	var val interface{}
	tr(func(_ interface{}, input interface{}) interface{} {
		val = input
		isChanged = true
		return nil
	})(nil, s.Car())

	if isChanged {
		return Cons(
			val,
			func() Stream {
				return Sequence(tr, s.Cdr())
			})
	} else {
		return Sequence(tr, s.Cdr())
	}
}
