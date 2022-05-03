package writer

type ConfigWriter interface {
	SetSecret(contextName string, secretName string, secretEncodedValue string, force bool) error
	SetConfig(contextName string, configName string, configValue string, force bool) error
	AddContext(contextName string) error
	AddFileToRender(contextName string, fileIn string, fileOut string) error
	WriteConfig() error
}
