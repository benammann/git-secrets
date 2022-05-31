package main

import (
	"github.com/benammann/git-secrets/cmd"
	"time"
)

var version string = "v0.0.0-local"
var commit string = "local-rev"
var date string = time.Now().Format(time.RFC3339)

func main() {
	cmd.Execute(version, commit, date)
}
