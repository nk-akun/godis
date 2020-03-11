package main

import (
	"errors"
	"log"
)

func errorTrace(err error) error {
	if err != nil {
		log.Println("error tracing: ", err.Error())
	}
	return err
}

func errorNew(errorMsg string) error {
	return errors.New("error: " + errorMsg)
}

func f2() error {
	return errorTrace(errorNew("test error!"))
}

func f1() error {
	if err := f2(); err != nil {
		errorTrace(err)
		return err
	}
	return nil
}

func main() {
	if err := f1(); err != nil {
		log.Println(err.Error())
	}
}
