package cmd

import (
	"bufio"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/benammann/git-secrets/pkg/utility"
	"github.com/fatih/color"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type DecodedSecret struct {
	secret       config_generic.Secret
	decodedValue string
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use: "scan",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(projectCfgError)
	},
	Short: "Searches project files for leaked secrets",
	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()

		scanAll, _ := cmd.Flags().GetBool(FlagAll)
		verbose, _ := cmd.Flags().GetBool(FlagVerbose)

		stagedFiles, errStagedFiles := utility.GetStagedFiles(scanAll)
		if errStagedFiles != nil {
			cobra.CheckErr(errStagedFiles)
		}

		if len(stagedFiles) == 0 {
			fmt.Println("no staged files found. use --all to scan all files")
			return
		}

		var decodedSecrets []*DecodedSecret

		for _, context := range projectCfg.GetContexts() {
			contextSecrets := projectCfg.GetSecretsByContext(context.Name)
			for _, secret := range contextSecrets {
				decodedValue, errDecode := secret.GetPlainValue(cmd.Context())
				if errDecode != nil {
					color.Yellow("Warning: could not decode secret %s from context %s, skipping this secret\n", secret.GetName(), secret.GetOriginContext().Name)
					continue
				}
				decodedSecrets = append(decodedSecrets, &DecodedSecret{secret: secret, decodedValue: decodedValue})
			}
		}

		var leakedSecrets []struct {
			fileName    string
			line        int
			lineContent string
			secret      *DecodedSecret
		}

		registerLeakedSecret := func(fileName string, line int, lineContent string, secret *DecodedSecret) {

			lineContent = strings.ReplaceAll(lineContent, secret.decodedValue, "************")

			leakedSecrets = append(leakedSecrets, struct {
				fileName    string
				line        int
				lineContent string
				secret      *DecodedSecret
			}{fileName: fileName, line: line, lineContent: lineContent, secret: secret})
		}

		var failures []struct {
			fileName string
			err      error
		}

		registerFailure := func(fileName string, err error) {
			failures = append(failures, struct {
				fileName string
				err      error
			}{fileName: fileName, err: err})
		}

		stagedFilesChunks := utility.ChunkStringSlice(stagedFiles, 10000)

		for _, stagedFilesChunk := range stagedFilesChunks {
			var filesWaitGroup sync.WaitGroup
			for _, stagedFile := range stagedFilesChunk {
				filesWaitGroup.Add(1)
				go func(stagedFile string) {

					if verbose {
						fmt.Println(stagedFile)
					}

					f, err := os.OpenFile(stagedFile, os.O_RDONLY, os.ModePerm)
					if err != nil {
						registerFailure(stagedFile, fmt.Errorf("open file error: %s", err.Error()))
						return
					}
					defer f.Close()

					sc := bufio.NewScanner(f)

					currentLine := 0
					for sc.Scan() {

						lineContent := sc.Text()
						currentLine++

						for _, secret := range decodedSecrets {
							if strings.Contains(lineContent, secret.decodedValue) {
								registerLeakedSecret(stagedFile, currentLine, lineContent, secret)
							}
						}

					}
					if err := sc.Err(); err != nil {
						registerFailure(stagedFile, fmt.Errorf("scan file error: %s", err.Error()))
						return
					}

					filesWaitGroup.Done()
				}(stagedFile)
			}
			filesWaitGroup.Wait()
		}

		elapsed := time.Since(start)

		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		if len(leakedSecrets) > 0 || len(failures) > 0 {

			if verbose {
				fmt.Printf("\n")
			}

			for _, failure := range failures {
				fmt.Printf("%s - error: %s\n", red(failure.fileName), yellow(failure.err.Error()))
			}

			for _, leakedSecret := range leakedSecrets {
				fmt.Printf("%s:%s - secret %s from context %s is present\n", red(leakedSecret.fileName), yellow(leakedSecret.line), yellow(leakedSecret.secret.secret.GetName()), yellow(leakedSecret.secret.secret.GetOriginContext().Name))
				fmt.Printf("%s%d | %s\n\n", yellow("> "), leakedSecret.line, leakedSecret.lineContent)
			}

			color.Red("Searched in %d files for %d secrets in %s \n", len(stagedFiles), len(decodedSecrets), elapsed)

			os.Exit(1)
		} else {

			if verbose {
				fmt.Printf("\n")
			}
			color.Green("All files are clean of any leaked secret contained in .git-secrets.json\n")
			color.Green("Searched in %d files for %d secrets in %s \n", len(stagedFiles), len(decodedSecrets), elapsed)
		}

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolP(FlagAll, "a", false, "Scan all files that are contained in the git repo")
	scanCmd.Flags().BoolP(FlagVerbose, "v", false, "List the scanned files")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
