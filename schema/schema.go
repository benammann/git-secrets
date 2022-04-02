package schema

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
)

//go:embed def
var schemaDefinitions embed.FS

type FileName string

var V1 FileName = "v1.json"

// GetSchemaContents reads the schema contents from the embed fs
func GetSchemaContents(schemaFileName FileName) []byte {
	fileContents, fileErr := schemaDefinitions.ReadFile(fmt.Sprintf("def/%s", schemaFileName))
	if fileErr != nil {
		cobra.CheckErr(fmt.Errorf("could not load schema %s: %s", fileErr.Error()))
	}
	return fileContents
}
