package main

import (
	"bufio"
	"io"
	"reflect"
	"os"
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

type LoadFunc func(string, interface{}) error

func MakeLoadFunc(readFrom interface{}) LoadFunc {
	return func(fname string, dst interface{}) error {
		file, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer file.Close()

		in := []reflect.Value{reflect.ValueOf(file)}
		out := reflect.ValueOf(readFrom).Call(in)
		outval, errval := out[0], out[1]

		if !errval.IsNil() {
			return errval.Interface().(error)
		}

		dstval := reflect.ValueOf(dst).Elem()
		// If we can assign to the destination, then do so.
		if outval.Type().AssignableTo(dstval.Type()) {
			dstval.Set(outval)
			return nil
		}
		// Otherwise, try de-referencing output value if possible.
		if outval.Kind() == reflect.Ptr {
			yval := outval.Elem()
			if yval.Type().AssignableTo(dstval.Type()) {
				dstval.Set(yval)
				return nil
			}
		}
		panic("could not assign to dst")
	}
}

func Load(readFrom interface{}, fname string, dst interface{}) error {
	return MakeLoadFunc(readFrom)(fname, dst)
}

type SaveFunc func(string, interface{}) error

func MakeSaveFunc(readFrom interface{}) SaveFunc {
	return func(fname string, dst interface{}) error {
		file, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer file.Close()

		in := []reflect.Value{reflect.ValueOf(file), reflect.ValueOf(dst)}
		out := reflect.ValueOf(readFrom).Call(in)
		errval := out[0]

		if !errval.IsNil() {
			return errval.Interface().(error)
		}
		return nil
	}
}

func Save(readFrom interface{}, fname string, dst interface{}) error {
	return MakeSaveFunc(readFrom)(fname, dst)
}
