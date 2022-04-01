package render

import (
	"encoding/base64"
	"html/template"
)

// getTemplateFunctions are added to the template and can be executed
func getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"Base64Encode": templateFunctionBase64Encode,
	}
}

// templateFunctionBase64Encode takes the current value and returns it as a base64 value
// should be used for kubernetes secrets
func templateFunctionBase64Encode(args ...interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(args[0].(string)))
}

// createNewTemplate creates a new template engine with all the extensions based on the file name
func createNewTemplate(pathToFile string) (*template.Template, error) {

	// create the new engine with file base name
	tpl := template.New("")

	// add the template functions
	tpl.Funcs(getTemplateFunctions())

	tpl, err := tpl.ParseFiles(pathToFile)

	return tpl, err
}
