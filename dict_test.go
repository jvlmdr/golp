package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"strconv"
	"testing"
)

func TestPivot(t *testing.T) {
	const eps = 1e-3

	for i := 0; i < 10; i++ {
		t.Logf("case %d", i+1)
		dict, err := loadDict(fmt.Sprintf("dict%d", i+1))
		if err != nil {
			t.Fatal(err)
		}
		soln, err := loadSoln(fmt.Sprintf("dict%d.output", i+1))
		if err != nil {
			t.Fatal(err)
		}

		enter, final := dict.Enter()
		if final {
			t.Fatalf("wrong: final: got %v, want %v", final, false)
		}

		leave, unbounded := dict.Leave(enter)
		if unbounded != soln.Unbounded {
			t.Fatalf("wrong: unbounded: got %v, want %v", unbounded, soln.Unbounded)
		}
		if unbounded || soln.Unbounded {
			continue
		}

		if dict.NonBasic[enter] != soln.Enter {
			t.Fatalf("wrong enter index: got %d, want %d", dict.NonBasic[enter], soln.Enter)
		}
		if dict.Basic[leave] != soln.Leave {
			t.Fatalf("wrong leave index: got %d, want %d", dict.Basic[leave], soln.Leave)
		}

		dict = dict.Pivot(enter, leave)
		if !approx(dict.D, soln.Objective, eps) {
			t.Fatalf("wrong objective: got %g, want %g", dict.D, soln.Objective)
		}
	}
}

func approx(got, want, eps float64) bool {
	return math.Abs(got-want) <= eps*math.Abs(want)
}

type Solution struct {
	Unbounded bool
	Enter     int
	Leave     int
	Objective float64
}

func loadSoln(fname string) (Solution, error) {
	file, err := os.Open(fname)
	if err != nil {
		return Solution{}, err
	}
	defer file.Close()
	return readSolnFrom(file)
}

func readSolnFrom(r io.Reader) (Solution, error) {
	scanner := bufio.NewScanner(r)
	var line string

	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// First line contains UNBOUNDED or an integer.
	if line == "UNBOUNDED" {
		return Solution{Unbounded: true}, nil
	}
	enter, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return Solution{}, err
	}

	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// Second line contains an integer.
	leave, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return Solution{}, err
	}

	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// Third line contains a number.
	objective, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return Solution{}, err
	}

	soln := Solution{Enter: int(enter), Leave: int(leave), Objective: objective}
	return soln, nil
}
