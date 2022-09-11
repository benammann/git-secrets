package render

import (
	"encoding/base64"
	"github.com/spf13/afero"
	"github.com/tcnksm/go-gitconfig"
	"html/template"
	"io/fs"
)

type AferoConvFs struct {
	aferoFs afero.Fs
}

func (ac AferoConvFs) Open(name string) (fs.File, error) {
	return ac.aferoFs.Open(name)
}

// getTemplateFunctions are added to the template and can be executed
func getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"Base64Encode": templateFunctionBase64Encode,
		"GitConfig":    templateFunctionGitConfig,
	}
}

// templateFunctionBase64Encode takes the current value and returns it as a base64 value
// should be used for kubernetes secrets
func templateFunctionBase64Encode(args ...interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(args[0].(string)))
}

// templateFunctionGitConfig allows reading values from git
func templateFunctionGitConfig(args ...interface{}) interface{} {
	val, err := gitconfig.Local(args[0].(string))
	if err != nil {
		globalVal, errGlobal := gitconfig.Global(args[0].(string))
		if errGlobal != nil {
			return ""
		}
		return globalVal
	}
	return val
}

func createTemplate(fs afero.Fs, pathToFile string) (*template.Template, error) {

	// create the new engine with file base name
	tpl := template.New("")

	// add the template functions
	tpl.Funcs(getTemplateFunctions())

	tpl, err := tpl.ParseFS(AferoConvFs{aferoFs: fs}, pathToFile)

	if err != nil {
		return nil, err
	}

	return tpl, err
}