package main

import (
	"errors"
	"fmt"
	"strings"
)

type Logger interface {
	Log(msg string, attrs ...interface{}) error
}

type StdLogger struct {}

func (s StdLogger) Log(msg string, attrs ...interface{}) error {
	if len(attrs) == 0 {
		fmt.Println(msg)
		return nil
	}
	if len(attrs) % 2 != 0 {
		return errors.New("there must be an even number of attrs")
	}

	var nameValuePairs []string
	for i := 0; i < len(attrs); i += 2 {
		nameValuePairs = append(nameValuePairs, fmt.Sprintf("%v=%v", attrs[i], attrs[i+1]))
	}

	fmt.Printf("%s : [%s]\n", msg, strings.Join(nameValuePairs, ", "))

	return nil
}
