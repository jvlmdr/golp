package main

import (
	"bufio"
	"io"
)

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
