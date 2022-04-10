package writer

type ConfigWriter interface {
	AddSecret(contextName string, secretName string, secretEncodedValue string) error
	AddContext(contextName string) error
	AddFileToRender(contextName string, fileIn string, fileOut string) error
	WriteConfig() error
}
