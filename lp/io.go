package lp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func lineOfInts(line string) ([]int, error) {
	words := strings.Fields(line)
	nums := make([]int, len(words))
	for i, str := range words {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		nums[i] = int(num)
	}
	return nums, nil
}

func lineOfFloats(line string) ([]float64, error) {
	words := strings.Fields(line)
	nums := make([]float64, len(words))
	for i, str := range words {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		nums[i] = num
	}
	return nums, nil
}

func readLine(scanner *bufio.Scanner) error {
	if !scanner.Scan() {
		err := scanner.Err()
		if err == nil {
			return io.EOF
		}
		return err
	}
	return nil
}

// Reads a dictionary in the University of Colorado format.
func ReadDictColoradoFrom(r io.Reader) (*Dict, error) {
	scanner := bufio.NewScanner(r)

	// First line contains dimensions.
	if err := readLine(scanner); err != nil {
		return nil, err
	}
	dims, err := lineOfInts(scanner.Text())
	if err != nil {
		return nil, err
	}
	if len(dims) != 2 {
		return nil, errors.New("wrong number of dimensions")
	}
	m, n := dims[0], dims[1]

	// Second line contains basic indices.
	if err := readLine(scanner); err != nil {
		return nil, err
	}
	basic, err := lineOfInts(scanner.Text())
	if err != nil {
		return nil, err
	}
	if len(basic) != m {
		return nil, errors.New("wrong number of basic vars")
	}

	// Third line contains non-basic indices.
	if err := readLine(scanner); err != nil {
		return nil, err
	}
	nonbasic, err := lineOfInts(scanner.Text())
	if err != nil {
		return nil, err
	}
	if len(nonbasic) != n {
		return nil, errors.New("wrong number of non-basic vars")
	}

	// Fourth line contains b.
	if err := readLine(scanner); err != nil {
		return nil, err
	}
	b, err := lineOfFloats(scanner.Text())
	if err != nil {
		return nil, err
	}
	if len(b) != m {
		return nil, errors.New("constraint constants are wrong length")
	}

	// Following m lines contain coefficients.
	a := make([][]float64, m)
	for i := 0; i < m; i++ {
		if err := readLine(scanner); err != nil {
			return nil, err
		}
		coeff, err := lineOfFloats(scanner.Text())
		if err != nil {
			return nil, err
		}
		if len(coeff) != n {
			return nil, errors.New("constraint coefficients are wrong length")
		}
		a[i] = coeff
	}

	// Last line is objective constant and coefficients.
	if err := readLine(scanner); err != nil {
		return nil, err
	}
	obj, err := lineOfFloats(scanner.Text())
	if err != nil {
		return nil, err
	}
	if len(obj) != n+1 {
		msg := fmt.Sprint("objective coefficients are wrong length: ", len(obj))
		return nil, errors.New(msg)
	}
	d, c := obj[0], obj[1:]

	return &Dict{basic, nonbasic, a, b, c, d}, nil
}
