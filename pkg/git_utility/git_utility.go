package gitutil

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedFiles returns the staged files in git
func GetStagedFiles(all bool) ([]string, error) {

	var gitCommand *exec.Cmd
	if all {
		gitCommand = exec.Command("git", "ls-files")
	} else {
		gitCommand = exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	}

	output, errExec := gitCommand.Output()
	if errExec != nil {
		return nil, fmt.Errorf("could not resolve staged files: %s / %s", errExec.Error(), string(output))
	}
	stagedFiles := strings.Split(strings.ReplaceAll(string(output), "\r\n", "\n"), "\n")
	var filteredStagedFiles []string
	for _, stagedFile := range stagedFiles {
		if stagedFile == "" || stagedFile == " " {
			continue
		}
		filteredStagedFiles = append(filteredStagedFiles, stagedFile)
	}
	return filteredStagedFiles, nil
}
