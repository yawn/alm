package main

import (
	"github.com/yawn/alm/command"
)

var (
	build   string
	time    string
	version string
)

func main() {
	command.Execute(version, build, time)
}
