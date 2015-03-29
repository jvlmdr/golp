package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jvlmdr/golp/lp"
)

// The result of trying to pivot.
type Solution struct {
	Unbounded bool
	Final     bool
	// Enter and Leave are labels not indices!
	Enter int
	Leave int
	// The dictionary after the pivot operation.
	Dict *lp.Dict
}

func WriteSolutionTo(w io.Writer, soln Solution) error {
	if soln.Final {
		_, err := fmt.Fprintln(w, "FINAL")
		return err
	}
	if soln.Unbounded {
		_, err := fmt.Fprintln(w, "UNBOUNDED")
		return err
	}

	if _, err := fmt.Fprintln(w, soln.Enter); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, soln.Leave); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, soln.Dict.Obj()); err != nil {
		return err
	}
	return nil
}

func ReadSolutionFrom(r io.Reader) (Summary, error) {
	scanner := bufio.NewScanner(r)
	var line string

	if err := readLine(scanner); err != nil {
		return Summary{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// First line contains UNBOUNDED or an integer.
	if line == "UNBOUNDED" {
		return Summary{Unbounded: true}, nil
	}
	enter, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return Summary{}, err
	}

	if err := readLine(scanner); err != nil {
		return Summary{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// Second line contains an integer.
	leave, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return Summary{}, err
	}

	if err := readLine(scanner); err != nil {
		return Summary{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	// Third line contains a number.
	objective, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return Summary{}, err
	}

	soln := Summary{Enter: int(enter), Leave: int(leave), Objective: objective}
	return soln, nil
}
