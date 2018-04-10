package infinite_stream_go

import "testing"

func integerStartFrom(n int) Stream {
	return Cons(n, func() Stream {
		return integerStartFrom(n + 1)
	})
}

func TestStream_Car(t *testing.T) {
	cp := Cons(1, func() Stream {
		return Tail{}
	})
	if cp.Car() != 1 {
		t.Errorf("cp.Car() = %v, want 1", cp.Car())
	}
}

func TestStream_Cdr(t *testing.T) {
	cp := Cons(1, func() Stream {
		return Tail{}
	})
	if _, ok := cp.Cdr().(Tail); !ok {
		t.Errorf("cp.Cdr() = %v, want Tail", cp.Cdr())
	}
}

func TestFrom(t *testing.T) {
	s := From(1, 2, 3)
	if s.Car() != 1 {
		t.Errorf("s.Car() = %v, want 1", s.Car())
	}
	if s.Cdr().Car() != 2 {
		t.Errorf("s.Cdr().Car() = %v, want 2", s.Cdr().Car())
	}
	if s.Cdr().Cdr().Car() != 3 {
		t.Errorf("s.Cdr().Cdr().Car() = %v, want 3", s.Cdr().Cdr().Car())
	}
	if s.Cdr().Cdr().Cdr().IsTail() != true {
		t.Error("s.Cdr().Cdr().Cdr().IsTail() is false, want true")
	}
}

func TestRef(t *testing.T) {
	ns := integerStartFrom(1)
	if Ref(0, ns) != 1 {
		t.Errorf("Ref(ns, 0) = %v, want 1", Ref(0, ns))
	}
	if Ref(1, ns) != 2 {
		t.Errorf("ns.Ref(1) = %v, want 2", Ref(1, ns))
	}
	if Ref(99, ns) != 100 {
		t.Errorf("Ref(ns,99) = %v, want 100", Ref(99, ns))
	}

}

func TestTake(t *testing.T) {
	ns := integerStartFrom(1)
	taken := Take(2, ns)
	if taken.Car() != 1 {
		t.Errorf("taken.Car() = %v, want 1", taken.Car())
	}
	if taken.Cdr().Car() != 2 {
		t.Errorf("taken.Cdr().Car() = %v, want 2", taken.Cdr().Car())
	}
	if taken.Cdr().Cdr().IsTail() != true {
		t.Error("taken.Cdr().Cdr().IsTail() = false, want true")
	}
}

func TestConj(t *testing.T) {
	s1 := Cons(1, TailCdr())
	s2 := Conj(s1, 2)
	if s2.Car() != 1 {
		t.Errorf("s2.Car() = %v, want 1", s2.Car())
	}
	if s2.Cdr().Car() != 2 {
		t.Errorf("s2.Cdr().Car() = %v, want 2", s2.Cdr().Car())
	}
	if s2.Cdr().Cdr().IsTail() != true {
		t.Error("s2.Cdr().Cdr().IsTail() is false, want true")
	}
}

func TestConcat(t *testing.T) {
	s1 := From(1, 2)
	s2 := From(3, 4)
	concated := Concat(s1, func() Stream { return s2 })
	if Ref(0, concated) != 1 {
		t.Errorf("Ref(0, concated) = %v, want 1", Ref(0, concated))
	}
	if Ref(1, concated) != 2 {
		t.Errorf("Ref(1, concated) = %v, want 2", Ref(1, concated))
	}
	if Ref(2, concated) != 3 {
		t.Errorf("Ref(2, concated) = %v, want 3", Ref(2, concated))
	}
	if Ref(3, concated) != 4 {
		t.Errorf("Ref(3, concated) = %v, want 4", Ref(3, concated))
	}
}

func TestReduce(t *testing.T) {
	ns := integerStartFrom(1)
	ten := Take(10, ns)
	reduced := Reduce(func(result interface{}, input interface{}) interface{} {
		sum, ok1 := result.(int)
		n, ok2 := input.(int)
		if ok1 && ok2 {
			return sum + n
		}
		return result
	}, 0, ten)
	v, ok := reduced.(int)
	if !ok {
		t.Errorf("reduce value = %v, want int", v)
	}
	if v != 55 {
		t.Errorf("reduce value = %v, want 55", v)
	}
}

func TestStep(t *testing.T) {
	// map
	mt := Map(func(input interface{}) interface{} {
		if n, ok := input.(int); ok {
			return n * n
		}
		return input
	})
	mv1, mChanged1 := Step(mt, 5)
	if mChanged1 != true {
		t.Error("mChanged1 is false, want true")
	}
	if mv1.Car() != 5*5 {
		t.Errorf("mv1.Car() = %v, want 25", mv1)
	}
	if mv1.Cdr().IsTail() != true {
		t.Error("mv1.Cdr().IsTail() is false, want true")
	}
	mv2, mChanged2 := Step(mt, 9)
	if mChanged2 != true {
		t.Error("mChanged2 is false, want true")
	}
	if mv2.Car() != 9*9 {
		t.Errorf("mv2.Car() = %v, want 81", mv2)
	}
	if mv2.Cdr().IsTail() != true {
		t.Error("mv2.Cdr().IsTail() is false, want true")
	}

	// filter
	ft := Filter(func(input interface{}) bool {
		if n, ok := input.(int); ok {
			return n%2 == 0
		}
		return false
	})
	_, fChanged1 := Step(ft, 1)
	if fChanged1 != false {
		t.Error("fChanged1 is true, want false")
	}
	fv1, fChanged4 := Step(ft, 2)
	if fChanged4 != true {
		t.Error("fChanged4 is false, want true")
	}
	if fv1.Car() != 2 {
		t.Errorf("fv1.Car() = %v, want 2", fv1)
	}
	if fv1.Cdr().IsTail() != true {
		t.Error("fv1.Cdr().IsTail() is false, want true")
	}

	// mapcat
	mct := func(reducing Reducer) Reducer {
		return func(result interface{}, input interface{}) interface{} {
			v := From(input, input) // twins
			return Reduce(reducing, result, v)
		}
	}
	mctv1, mctOk1 := Step(mct, 3)
	if !mctOk1 {
		t.Error("mctOk1 is false, want true")
	}
	if mctv1.Car() != 3 {
		t.Errorf("mctv1.Car() = %v, want 3", mctv1.Car())
	}
	if mctv1.Cdr().Car() != 3 {
		t.Errorf("mctv1.Cdr().Car() = %v, want 3", mctv1.Cdr().Car())
	}
	if mctv1.Cdr().Cdr().IsTail() != true {
		t.Error("mctv1.Cdr().Cdr().IsTail() is false, want true")
	}
}

func TestSequence(t *testing.T) {
	ns := integerStartFrom(1)

	// map
	s1 := Sequence(Map(func(input interface{}) interface{} {
		if n, ok := input.(int); ok {
			return n * n
		}
		return input
	}), ns)
	if Ref(0, s1) != 1 {
		t.Errorf("Ref(0, s1) = %v, want 1", Ref(0, s1))
	}
	if Ref(1, s1) != 2*2 {
		t.Errorf("Ref(1, s1) = %v, want 4", Ref(1, s1))
	}
	if Ref(2, s1) != 3*3 {
		t.Errorf("Ref(2, s1) = %v, want 4", Ref(2, s1))
	}

	// filter
	s2 := Sequence(Filter(func(input interface{}) bool {
		if n, ok := input.(int); ok {
			return n%2 == 0
		}
		return false
	}), ns)
	if Ref(0, s2) != 2 {
		t.Errorf("Ref(0, s2) = %v, want 2", Ref(0, s2))
	}
	if Ref(1, s2) != 4 {
		t.Errorf("Ref(2, s2) = %v, want 4", Ref(1, s2))
	}
	if Ref(2, s2) != 6 {
		t.Errorf("Ref(2, s2) = %v, want 6", Ref(2, s2))
	}
}
