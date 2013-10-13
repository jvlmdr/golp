package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fname := os.Args[1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dict, err := ReadDictFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	dict.Fprint(os.Stderr)

	in, final := dict.nextIn()
	if final {
		fmt.Println("FINAL")
		return
	}
	out, unbounded := dict.nextOut(in)
	if unbounded {
		fmt.Println("UNBOUNDED")
		return
	}

	fmt.Println(dict.NonBasic[in])
	fmt.Println(dict.Basic[out])
	dict = dict.Pivot(in, out)
	fmt.Println(dict.D)
	dict.Fprint(os.Stderr)

	//	d := new(Dictionary)
	//	d.Basic = []int{4, 5, 6}
	//	d.NonBasic = []int{1, 2, 3}
	//	d.A = [][]float64{
	//		[]float64{-2, -3, -1},
	//		[]float64{-4, -1, -2},
	//		[]float64{-3, -4, -2},
	//	}
	//	d.B = []float64{5, 11, 8}
	//	d.C = []float64{5, 4, 3}
	//	d.D = 0

	//	d = d.Pivot(0, 0)
	//	fmt.Println(d)
}

type LP struct {
	C    []float64
	Ineq []Inequality
}

type Inequality struct {
	A []float64
	B float64
}

type Affine struct {
	A []float64
	B float64
}
