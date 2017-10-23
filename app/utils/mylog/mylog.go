package mylog

import (
	"fmt"

	"github.com/fatih/color"
)

// EnableTrace disable printing trace
var EnableTrace = false

func Trace(format string, values ...interface{}) {
	if !EnableTrace {
		return
	}
	for i, v := range values {
		values[i] = color.CyanString("%v", v)
	}
	fmt.Printf(format+"\n", values...)
}

func Info(format string, values ...interface{}) {
	for i, v := range values {
		values[i] = color.GreenString("%v", v)
	}
	fmt.Printf(format+"\n", values...)
}

func Warn(format string, values ...interface{}) {
	for i, v := range values {
		values[i] = color.YellowString("%v", v)
	}
	fmt.Printf(format+"\n", values...)
}

func Error(format string, values ...interface{}) {
	for i, v := range values {
		values[i] = color.RedString("%v", v)
	}
	fmt.Printf(format+"\n", values...)
}
