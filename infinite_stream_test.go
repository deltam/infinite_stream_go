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
