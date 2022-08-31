package config_generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRenderTarget(t *testing.T) {
	newRenderTarget := NewRenderTarget("test")
	assert.Equal(t, "test", newRenderTarget.Name)
	assert.Len(t, newRenderTarget.FilesToRender, 0)
}

func TestRenderTarget_AddFileToRender(t *testing.T) {
	newRenderTarget := NewRenderTarget("test")
	assert.NoError(t, newRenderTarget.AddFileToRender("fileIn", "fileOut"))
	assert.Len(t, newRenderTarget.FilesToRender, 1)
	assert.Equal(t, "fileIn", newRenderTarget.FilesToRender[0].FileIn)
	assert.Equal(t, "fileOut", newRenderTarget.FilesToRender[0].FileOut)
	assert.Error(t, newRenderTarget.AddFileToRender("fileIn", "fileOut"))
}

func TestRepository_AddRenderTarget(t *testing.T) {
	newRenderTarget := NewRenderTarget("test")
	repo := initRepository(t, TestFileBlankDefault, "default")
	assert.NoError(t, repo.AddRenderTarget(newRenderTarget))
	assert.Error(t, repo.AddRenderTarget(newRenderTarget))
}

func TestRepository_GetRenderTarget(t *testing.T) {

	newRenderTarget := NewRenderTarget("test")
	repo := initRepository(t, TestFileBlankDefault, "default")
	assert.NoError(t, repo.AddRenderTarget(newRenderTarget))

	assert.Equal(t, newRenderTarget, repo.GetRenderTarget("test"))
	assert.Nil(t, repo.GetRenderTarget("missing"))

}

func TestRepository_HasRenderTarget(t *testing.T) {

	newRenderTarget := NewRenderTarget("test")
	repo := initRepository(t, TestFileBlankDefault, "default")
	assert.NoError(t, repo.AddRenderTarget(newRenderTarget))

	assert.True(t, repo.HasRenderTarget("test"))
	assert.False(t, repo.HasRenderTarget("missing"))

}

func TestRepository_RenderTargetNames(t *testing.T) {
	newRenderTarget := NewRenderTarget("test")
	newRenderTarget2 := NewRenderTarget("test2")
	repo := initRepository(t, TestFileBlankDefault, "default")
	assert.NoError(t, repo.AddRenderTarget(newRenderTarget))
	assert.NoError(t, repo.AddRenderTarget(newRenderTarget2))
	assert.Equal(t, []string{"test", "test2"}, repo.RenderTargetNames())
}
