package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Solution struct {
	Unbounded bool
	Enter     int
	Leave     int
	Objective float64
}

func ReadSolutionFrom(r io.Reader) (Solution, error) {
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
