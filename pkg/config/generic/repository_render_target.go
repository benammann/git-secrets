package config_generic

import "fmt"

type RenderTarget struct {
	Name          string
	FilesToRender []*FileToRender
}

type FileToRender struct {
	FileIn  string
	FileOut string
}

func NewRenderTarget(name string) *RenderTarget {
	return &RenderTarget{
		Name: name,
	}
}

func (c *Repository) AddRenderTarget(target *RenderTarget) error {
	if c.HasRenderTarget(target.Name) {
		return fmt.Errorf("the render target %s already exists", target.Name)
	}
	c.renderTargets = append(c.renderTargets, target)
	return nil
}

func (c *Repository) HasRenderTarget(targetName string) bool {
	return c.GetRenderTarget(targetName) != nil
}

func (c *Repository) GetRenderTarget(targetName string) *RenderTarget {
	for _, renderTarget := range c.renderTargets {
		if renderTarget.Name == targetName {
			return renderTarget
		}
	}
	return nil
}

func (c *Repository) RenderTargetNames() (names []string) {
	for _, renderTarget := range c.renderTargets {
		names = append(names, renderTarget.Name)
	}
	return names
}

// AddFileToRender adds a file to render which is later used by the rendering engine
func (c *RenderTarget) AddFileToRender(fileIn string, fileOut string) error {

	// check if output file is double defined
	for _, fileToRender := range c.FilesToRender {
		if fileToRender.FileOut == fileOut {
			return fmt.Errorf("output file %s is already defined on target %s", fileOut, c.Name)
		}
	}

	c.FilesToRender = append(c.FilesToRender, &FileToRender{
		FileIn:  fileIn,
		FileOut: fileOut,
	})

	return nil
}
