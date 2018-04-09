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
	fmt.Println()
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

/*
// エラトステネスのふるい
func sieve(s is.Stream) is.Stream {
	return is.Cons(
		s.Car(),
		func() is.Stream {
			if d, ok := s.Car().(int); ok {
				return sieve(
					s.Cdr().Filter(notDivisible(d)),
				)
			}
			return s
		},
	)
}
*/
func main() {
	fmt.Println("natural numbers")
	ns := integerStartFrom(1)
	fmt.Println(ns.Car())
	fmt.Println(ns.Cdr().Car())
	fmt.Println(is.Ref(100, ns))
	displayLine(ns, 20)

	/*
		fmt.Println("no sevens")
		noSevens := ns.Filter(notDivisible(7))
		displayLine(noSevens, 20)

		fmt.Println("primes")
		nsFrom2 := ns.Cdr()
		primes := sieve(nsFrom2)
		displayLine(primes, 20)
		fmt.Println(primes.Ref(1000))

		fmt.Println("natural numbers by Iterate")
		ns2 := is.IterateInt(func(n int) int { return n + 1 }, 1)
		displayLine(ns2, 20)
	*/
}
