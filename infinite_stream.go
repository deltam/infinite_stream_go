package infinite_stream_go

type Stream struct {
	car int
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

func Cons(a int, b func() Stream) (cp Stream) {
	cp.car = a
	cp.cdr = delay(b)
	return cp
}

func (cp Stream) Car() int {
	return cp.car
}

func (cp Stream) Cdr() Stream {
	return force(cp.cdr)
}

func (cp Stream) Ref(n int) int {
	if n == 0 {
		return cp.Car()
	} else {
		return cp.Cdr().Ref(n - 1)
	}
}

func (cp Stream) Filter(pred func(int) bool) Stream {
	if pred(cp.Car()) {
		return Cons(cp.Car(),
			func() Stream {
				return cp.Cdr().Filter(pred)
			})
	} else {
		return cp.Cdr().Filter(pred)
	}
}

func Iterate(fn func(int) int, init int) Stream {
	return Cons(init, func() Stream {
		return Iterate(fn, fn(init))
	})
}
