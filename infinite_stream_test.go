package infinite_stream_go

import "testing"

var ns Stream = Iterate(func(n int) int { return n + 1 }, 1)

func TestCons(t *testing.T) {
	s1 := Cons(1, nil)
	if s1.car != 1 {
		t.Errorf("Cons(1,nil).Car() = %v, want 1", s1.car)
	}
	if s1.cdr == nil {
		t.Errorf("Cons(1,nil).Cdr() = func() Stream, want nil")
	}
}

func TestStream_Car(t *testing.T) {
	if ns.Car() != 1 {
		t.Errorf("ns.Car() = %v, want 1", ns.Car())
	}
}

func TestStream_Cdr(t *testing.T) {
	if ns.Cdr().Car() != 2 {
		t.Errorf("ns.Cdr().Car() = %v, want 2", ns.Cdr().Car())
	}
}
