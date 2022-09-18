package writer

type ConfigWriter interface {
	SetEncryptedSecret(contextName string, secretName string, secretEncodedValue string, force bool) error
	SetGcpSecret(contextName string, secretName string, resourceId string, force bool) error
	SetConfig(contextName string, configName string, configValue string, force bool) error
	AddContext(contextName string) error
	AddFileToRender(targetName string, fileIn string, fileOut string) error
	WriteConfig() error
}
