package main

import (
	"github.com/benammann/git-secrets/cmd"
)

var version = "0.0.0-local"
var commit = "n/a"
var date = "n/a"

func main() {
	cmd.Execute(version, commit, date)
}
