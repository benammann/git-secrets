package render

import (
	"bytes"
	"fmt"
	config_generic "github.com/benammann/git-secrets/pkg/config/generic"
	"github.com/spf13/afero"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

type RenderingEngine struct {
	repository *config_generic.Repository
	fsIn afero.Fs
	fsOut afero.Fs
}

type RenderingContext struct {
	ContextName string
	File        *config_generic.FileToRender
	Secrets     config_generic.SecretsMap
	Configs     config_generic.ConfigMap
}

func NewRenderingEngine(repository *config_generic.Repository, fsIn afero.Fs, fsOut afero.Fs) *RenderingEngine {
	return &RenderingEngine{
		repository: repository,
		fsIn: fsIn,
		fsOut: fsOut,
	}
}

func (e *RenderingEngine) createTemplate(fileIn string) (*template.Template, error)  {
	return createTemplate(e.fsIn, fileIn)
}

// CreateRenderingContext creates the context which is used in the templates
func (e *RenderingEngine) CreateRenderingContext(fileToRender *config_generic.FileToRender) (*RenderingContext, error) {

	// decode the secrets
	secretsMap, errSecrets := e.repository.GetSecretsMapDecoded()
	if errSecrets != nil {
		return nil, fmt.Errorf("could not create context secrets: %s", errSecrets.Error())
	}

	// get the config values
	configMap := e.repository.GetConfigMap()

	return &RenderingContext{
		ContextName: e.repository.GetCurrent().Name,
		Secrets:     secretsMap,
		Configs:     configMap,
		File:        fileToRender,
	}, nil

}

// RenderFile renders the file and returns the rendered contents as string
func (e *RenderingEngine) RenderFile(fileToRender *config_generic.FileToRender) (usedContext *RenderingContext, fileContents string, err error) {

	var bytesOut bytes.Buffer
	usedContext, errExecute := e.ExecuteTemplate(fileToRender, &bytesOut)
	if errExecute != nil {
		return nil, "", fmt.Errorf("could not execute template: %s", errExecute.Error())
	}

	return usedContext, bytesOut.String(), nil

}

// WriteFile renders the file and writes it to its destination
func (e *RenderingEngine) WriteFile(fileToRender *config_generic.FileToRender) (usedContext *RenderingContext, err error) {

	// open the file
	fsFileOut, errFsFile := e.fsOut.OpenFile(fileToRender.FileOut, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if fsFileOut != nil {
		defer fsFileOut.Close()
	}
	if errFsFile != nil {
		return nil, fmt.Errorf("could not open %s: %s", fileToRender.FileOut, errFsFile.Error())
	}

	// execute the template and link the file stream
	usedContext, errExecute := e.ExecuteTemplate(fileToRender, fsFileOut)
	if errExecute != nil {
		return usedContext, fmt.Errorf("could not execute template: %s", errExecute.Error())
	}

	return usedContext, nil

}

// ExecuteTemplate executes the template and creates the rendering context
func (e *RenderingEngine) ExecuteTemplate(fileToRender *config_generic.FileToRender, writer io.Writer) (usedContext *RenderingContext, err error) {

	// create the rendering context
	usedContext, err = e.CreateRenderingContext(fileToRender)
	if err != nil {
		return nil, fmt.Errorf("could not create rendering context: %s", err.Error())
	}

	// create the template and execute
	tpl, errTpl := e.createTemplate(fileToRender.FileIn)
	if errTpl != nil {
		return nil, fmt.Errorf("error while reading template %s: %s", fileToRender.FileIn, errTpl.Error())
	}
	err = tpl.ExecuteTemplate(writer, filepath.Base(fileToRender.FileIn), usedContext)

	// return
	return usedContext, err

}
