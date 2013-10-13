package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Dictionary struct {
	Basic    []int
	NonBasic []int
	// basic = A nonbasic + b
	A [][]float64
	B []float64
	// objective = c' nonbasic + d
	C []float64
	D float64
}

func NewDict(m, n int) *Dictionary {
	d := new(Dictionary)
	d.Basic = make([]int, m)
	d.NonBasic = make([]int, n)
	d.A = make([][]float64, m)
	for i := range d.A {
		d.A[i] = make([]float64, n)
	}
	d.B = make([]float64, m)
	d.C = make([]float64, n)
	return d
}

func (dict *Dictionary) nextIn() (in int, final bool) {
	for i := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[i] > 0 {
			return i, false
		}
	}
	return 0, true
}

func (dict *Dictionary) nextOut(in int) (out int, unbounded bool) {
	// Find leaving variable.
	var (
		arg   int
		min   float64
		found bool
	)
	for i := range dict.Basic {
		// Must have negative constraint coefficient.
		if dict.A[i][in] >= 0 {
			continue
		}
		val := -dict.B[i] / dict.A[i][in]
		if !found || val < min {
			arg, min = i, val
			found = true
		}
	}
	return arg, !found
}

func (src *Dictionary) Pivot(in, out int) *Dictionary {
	m := len(src.Basic)
	n := len(src.NonBasic)
	dst := NewDict(m, n)

	// Copy the variable indices.
	copy(dst.Basic, src.Basic)
	copy(dst.NonBasic, src.NonBasic)
	// Swap in and out variables.
	dst.Basic[out] = src.NonBasic[in]
	dst.NonBasic[in] = src.Basic[out]

	dst.B[out] = -src.B[out] / src.A[out][in]
	for j := range dst.A[out] {
		if j == in {
			dst.A[out][j] = 1 / src.A[out][in]
		} else {
			dst.A[out][j] = -src.A[out][j] / src.A[out][in]
		}
	}

	for i := range dst.A {
		if i == out {
			continue
		}
		dst.B[i] = src.B[i] + src.A[i][in]*dst.B[out]
		for j := range dst.A[i] {
			if j == in {
				dst.A[i][j] = src.A[i][in] / src.A[out][j]
			} else {
				dst.A[i][j] = src.A[i][j] + src.A[i][in]*dst.A[out][j]
			}
		}
	}

	dst.D = src.D + src.C[in]*dst.B[out]
	for j := range dst.C {
		if j == in {
			dst.C[j] = src.C[in] / src.A[out][j]
		} else {
			dst.C[j] = src.C[j] + src.C[in]*dst.A[out][j]
		}
	}
	return dst
}

func (dict *Dictionary) longestCoeff(format string) int {
	var n int
	for _, ai := range dict.A {
		for _, aij := range ai {
			n = max(n, len(fmt.Sprintf(format, aij)))
		}
	}
	for _, bi := range dict.B {
		n = max(n, len(fmt.Sprintf(format, bi)))
	}
	for _, ci := range dict.C {
		n = max(n, len(fmt.Sprintf(format, ci)))
	}
	n = max(n, len(fmt.Sprintf(format, dict.D)))
	return n
}

func (dict *Dictionary) longestIndex(format string) int {
	var n int
	for _, ni := range dict.Basic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	for _, ni := range dict.NonBasic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	return n
}

func (dict *Dictionary) Fprint(w io.Writer) error {
	coeffLen := dict.longestCoeff("%+-.2g")
	coeff := "%+-" + fmt.Sprintf("%d", coeffLen) + ".2g"
	indexLen := dict.longestIndex("%d")
	index := "%" + fmt.Sprintf("%d", indexLen) + "d"

	var b bytes.Buffer
	for i := range dict.Basic {
		fmt.Fprintf(&b, "x"+index+" =", dict.Basic[i])
		fmt.Fprintf(&b, " "+coeff, dict.B[i])
		for j := range dict.NonBasic {
			fmt.Fprintf(&b, "  "+coeff+" x"+index, dict.A[i][j], dict.NonBasic[j])
		}
		b.WriteString("\n")
		if _, err := io.Copy(w, &b); err != nil {
			return err
		}
	}

	spacer := strings.Repeat(" ", indexLen)
	fmt.Fprint(&b, "z"+spacer+" =")
	fmt.Fprintf(&b, " "+coeff, dict.D)
	for j := range dict.NonBasic {
		fmt.Fprintf(&b, "  "+coeff+" x"+index, dict.C[j], dict.NonBasic[j])
	}
	b.WriteString("\n")
	if _, err := io.Copy(w, &b); err != nil {
		return err
	}
	return nil
}

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

func ReadDictFrom(r io.Reader) (*Dictionary, error) {
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

	return &Dictionary{basic, nonbasic, a, b, c, d}, nil
}
