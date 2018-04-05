package infinite_stream_go

type Stream struct {
	car interface{}
	cdr func() Stream
}

func delay(fn func() Stream) func() Stream {
	already_run := false
	var ret Stream
	return func() Stream {
		if !already_run {
			ret = force(fn)
			already_run = true
		}
		return ret
	}
}

func force(fn func() Stream) Stream {
	return fn()
}

func Cons(a interface{}, b func() Stream) (cp Stream) {
	cp.car = a
	cp.cdr = delay(b)
	return cp
}

func (s Stream) Car() interface{} {
	return s.car
}

func (s Stream) Cdr() Stream {
	return force(s.cdr)
}

func (s Stream) Ref(n int) interface{} {
	if n == 0 {
		return s.Car()
	} else {
		return s.Cdr().Ref(n - 1)
	}
}

func (s Stream) Filter(pred func(interface{}) bool) Stream {
	if pred(s.Car()) {
		return Cons(s.Car(),
			func() Stream {
				return s.Cdr().Filter(pred)
			})
	} else {
		return s.Cdr().Filter(pred)
	}
}

func Iterate(fn func(interface{}) interface{}, init interface{}) Stream {
	return Cons(init, func() Stream {
		return Iterate(fn, fn(init))
	})
}
