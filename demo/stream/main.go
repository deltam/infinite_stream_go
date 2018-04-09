package main

import (
	"fmt"

	is "github.com/deltam/infinite_stream_go"
)

func displayLine(s is.Stream, limit int) {
	cur := s
	for i := 0; i < limit; i++ {
		fmt.Printf("%d ", cur.Car())
		cur = cur.Cdr()
		if cur.IsTail() {
			break
		}
	}
	fmt.Print("\n\n")
}

func integerStartFrom(n int) is.Stream {
	return is.Cons(n, func() is.Stream {
		return integerStartFrom(n + 1)
	})
}

// dで割り切れないことを判定する関数を返す
func notDivisible(d int) func(interface{}) bool {
	return func(param interface{}) bool {
		if n, ok := param.(int); ok {
			return 0 != n%d
		}
		return false
	}
}

// エラトステネスのふるい
func sieve(s is.Stream) is.Stream {
	if s.IsTail() {
		return is.Tail{}
	}
	return is.Cons(
		s.Car(),
		func() is.Stream {
			if d, ok := s.Car().(int); ok {
				return sieve(
					is.Sequence(
						is.Filter(notDivisible(d)),
						s.Cdr()))
			}
			return is.Tail{}
		},
	)
}

func main() {
	fmt.Println("natural numbers")
	ns := integerStartFrom(1)
	fmt.Println(ns.Car())
	fmt.Println(ns.Cdr().Car())
	fmt.Println(is.Ref(100, ns))
	displayLine(ns, 20)

	fmt.Println("no sevens")
	noSevens := is.Sequence(is.Filter(notDivisible(7)), is.Take(30, ns))
	displayLine(noSevens, 20)

	fmt.Println("primes")
	nsFrom2 := integerStartFrom(2)
	primes := sieve(nsFrom2)
	displayLine(primes, 20)
	fmt.Printf("prime 1000th: %d\n", is.Ref(1000, primes))
}
