package command

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const app = "alm"

var this = Version{}

// Execute is the main entry point into the app
func Execute(version, build, time string) {

	whitespace := regexp.MustCompile(`(\s{2,})`)

	this.Build = build
	this.Time = time
	this.Version = version

	if _, err := rootCmd.ExecuteC(); err != nil {

		fmt.Fprintf(os.Stderr, "error: %s\n",
			whitespace.ReplaceAllString(err.Error(), ""))

		os.Exit(-1)

	}

}

