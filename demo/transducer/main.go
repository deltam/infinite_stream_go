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

func inc(input interface{}) interface{} {
	if n, ok := input.(int); ok {
		return n + 1
	}
	return input
}

func square(input interface{}) interface{} {
	if n, ok := input.(int); ok {
		return n * n
	}
	return input
}

func isEven(input interface{}) bool {
	if n, ok := input.(int); ok {
		return n%2 == 0
	}
	return false
}

func twins(a interface{}) interface{} {
	return is.From(a, a)
}

func mapcatting(fn func(interface{}) interface{}) is.Transducer {
	return func(reducing is.Reducer) is.Reducer {
		return func(result interface{}, input interface{}) interface{} {
			v := fn(input)
			if s, ok := v.(is.Stream); ok {
				return is.Reduce(reducing, result, s)
			} else {
				return result
			}
		}
	}
}

func taking(n int) is.Transducer {
	return func(reducing is.Reducer) is.Reducer {
		return func(result interface{}, input interface{}) interface{} {
			if 0 < n {
				n--
				return reducing(result, input)
			} else {
				return result
			}
		}
	}
}

func main() {
	fmt.Println("natural numbers")
	ns := integerStartFrom(1)
	displayLine(ns, 20)

	fmt.Println("map inc")
	mapInc := is.Reduce(
		is.Map(inc)(is.ConjReducer),
		is.Empty(),
		is.Take(10, ns))
	displayLine(mapInc.(is.Stream), 20)

	fmt.Println("map square")
	mapSquare := is.Reduce(
		is.Map(square)(is.ConjReducer),
		is.Empty(),
		is.Take(10, ns))
	displayLine(mapSquare.(is.Stream), 20)

	fmt.Println("filter even")
	filterEven := is.Reduce(
		is.Filter(isEven)(is.ConjReducer),
		is.Empty(),
		is.Take(10, ns))
	displayLine(filterEven.(is.Stream), 20)

	fmt.Println("filter even -> map square")
	filterMap := is.Reduce(
		is.Map(square)(
			is.Filter(isEven)(is.ConjReducer)),
		is.Empty(),
		is.Take(20, ns))
	displayLine(filterMap.(is.Stream), 20)

	fmt.Println("comp(filter even -> map square)")
	xform1 := is.Comp(
		is.Filter(isEven),
		is.Map(square),
	)
	v1 := xform1(is.ConjReducer)(is.Empty(), 1)
	fmt.Printf("xform 1 -> %v\n", v1)
	v2 := xform1(is.ConjReducer)(is.Empty(), 2)
	fmt.Printf("xform 2 -> %v\n", v2)
	fmt.Println()

	fmt.Println("transduce filter even -> map inc")
	xform2 := is.Comp(
		is.Filter(isEven),
		is.Map(inc),
	)
	transduced := is.Transduce(
		xform2,
		is.ConjReducer,
		is.Empty(),
		is.Take(20, ns),
	)
	displayLine(transduced.(is.Stream), 20)

	fmt.Println("sequence filter even -> map square")
	xform3 := is.Comp(
		is.Filter(isEven),
		is.Filter(func(input interface{}) bool {
			if n, ok := input.(int); ok {
				return n < 10
			} else {
				return false
			}
		}),
		is.Map(square),
		is.Map(inc),
	)
	sequenced := is.Sequence(xform3, is.Take(30, ns))
	displayLine(sequenced, 20)

	fmt.Println("mapcatting")
	mapcat := is.Sequence(
		mapcatting(twins),
		is.From(10, 20, 30),
	)
	displayLine(mapcat, 20)

	fmt.Println("taking")
	taken := is.Sequence(
		taking(3),
		is.From(1, 2, 3, 4, 5),
	)
	displayLine(taken, 20)
}
