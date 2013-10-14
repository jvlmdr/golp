package lp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Dict struct {
	Basic    []int
	NonBasic []int
	// basic = A nonbasic + b
	A [][]float64
	B []float64
	// objective = c' nonbasic + d
	C []float64
	D float64
}

func NewDict(m, n int) *Dict {
	d := new(Dict)
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

func (dict *Dict) Enter() (enter int, final bool) {
	// Find variable with lowest index.
	var (
		found bool
		arg   int
		min   int
	)
	for i := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[i] > 0 {
			index := dict.NonBasic[i]
			if !found || index < min {
				found = true
				arg, min = i, index
			}
		}
	}
	return arg, !found
}

func (dict *Dict) Leave(enter int) (leave int, unbounded bool) {
	// Find leaving variable.
	var (
		found    bool
		arg      int
		minVal   float64
		minIndex int
	)
	for i := range dict.Basic {
		// Must have negative constraint coefficient.
		if dict.A[i][enter] >= 0 {
			continue
		}
		val := -dict.B[i] / dict.A[i][enter]
		index := dict.Basic[i]

		if found {
			if val > minVal {
				continue
			} else if val == minVal {
				if index > minIndex {
					continue
				}
			}
		}
		found = true
		arg = i
		minVal = val
		minIndex = index
	}
	return arg, !found
}

func (src *Dict) Pivot(enter, leave int) *Dict {
	m := len(src.Basic)
	n := len(src.NonBasic)
	dst := NewDict(m, n)

	// Copy the variable indices.
	copy(dst.Basic, src.Basic)
	copy(dst.NonBasic, src.NonBasic)
	// Swap enter and leave variables.
	dst.Basic[leave] = src.NonBasic[enter]
	dst.NonBasic[enter] = src.Basic[leave]

	dst.B[leave] = -src.B[leave] / src.A[leave][enter]
	for j := range dst.A[leave] {
		if j == enter {
			dst.A[leave][j] = 1 / src.A[leave][enter]
		} else {
			dst.A[leave][j] = -src.A[leave][j] / src.A[leave][enter]
		}
	}

	for i := range dst.A {
		if i == leave {
			continue
		}
		dst.B[i] = src.B[i] + src.A[i][enter]*dst.B[leave]
		for j := range dst.A[i] {
			if j == enter {
				dst.A[i][j] = src.A[i][enter] / src.A[leave][j]
			} else {
				dst.A[i][j] = src.A[i][j] + src.A[i][enter]*dst.A[leave][j]
			}
		}
	}

	dst.D = src.D + src.C[enter]*dst.B[leave]
	for j := range dst.C {
		if j == enter {
			dst.C[j] = src.C[enter] / src.A[leave][j]
		} else {
			dst.C[j] = src.C[j] + src.C[enter]*dst.A[leave][j]
		}
	}
	return dst
}

func (dict *Dict) longestCoeff(format string) int {
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

func (dict *Dict) longestIndex(format string) int {
	var n int
	for _, ni := range dict.Basic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	for _, ni := range dict.NonBasic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	return n
}

func (dict *Dict) Fprint(w io.Writer) error {
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

func ReadDictFrom(r io.Reader) (*Dict, error) {
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
