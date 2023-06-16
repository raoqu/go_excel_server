package main

import (
	"github.com/TwiN/go-color"
)

func printError(err error) {
	if err != nil {
		println(color.InRed(err.Error()))
	}
}

func green(msg string) {
	println(color.InGreen(msg))
}

func red(msg string) {
	println(color.InRed(msg))
}

func yellow(msg string) {
	println(color.InYellow(msg))
}

func purple(msg string) {
	println(color.InPurple(msg))
}

func cyan(msg string) {
	println(color.InCyan(msg))
}

func gray(msg string) {
	println(color.InGray(msg))
}
