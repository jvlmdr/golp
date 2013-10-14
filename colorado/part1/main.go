package main

import (
	"github.com/jackvalmadre/golp/lp"

	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	const (
		dir = "files"
		n   = 5
	)

	for i := 1; i <= n; i++ {
		fmt.Printf("case %d:\n", i)

		dictFile := path.Join(dir, fmt.Sprintf("part%d.dict", i))
		solnFile := path.Join(dir, fmt.Sprintf("part%d.output", i))

		dict, err := loadDict(dictFile)
		if err != nil {
			log.Fatal(err)
		}

		if err := pivotToFile(dict, solnFile); err != nil {
			log.Print(err)
		}
	}
}

func pivotToFile(dict *lp.Dict, fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return pivot(dict, file)
}

func pivot(dict *lp.Dict, w io.Writer) error {
	// Get entering variable, check final.
	enter, final := dict.Enter()
	if final {
		if _, err := fmt.Fprintln(w, "FINAL"); err != nil {
			return err
		}
		return nil
	}

	// Get leaving variable, check unbounded.
	leave, unbounded := dict.Leave(enter)
	if unbounded {
		if _, err := fmt.Fprintln(w, "UNBOUNDED"); err != nil {
			return err
		}
		return nil
	}

	if _, err := fmt.Fprintln(w, dict.NonBasic[enter]); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, dict.Basic[leave]); err != nil {
		return err
	}

	dict = dict.Pivot(enter, leave)

	if _, err := fmt.Fprintln(w, dict.D); err != nil {
		return err
	}
	return nil
}

func loadDict(fname string) (*lp.Dict, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return lp.ReadDictFrom(file)
}
