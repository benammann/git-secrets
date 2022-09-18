package config_generic

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/spf13/afero"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

func NewGcpSecret(name string, resourceId string, originContext *Context) *GcpSecret {
	return &GcpSecret{
		Name:          name,
		ResourceId:    resourceId,
		OriginContext: originContext,
	}
}

type GcpSecret struct {

	// Name describes the name of the GcpSecret
	Name string

	// ResourceId hold the encodedValue in base64 of the GcpSecret
	ResourceId string

	// OriginContext references the configured context to decode the GcpSecret
	OriginContext *Context
}

func (s *GcpSecret) GetName() string {
	return s.Name
}

func (s *GcpSecret) GetOriginContext() *Context {
	return s.OriginContext
}

func (s *GcpSecret) GetPlainValue() (string, error) {

	file := s.OriginContext.GlobalConfig.GetGcpCredentialsFile(s.OriginContext.GcpCredentials)
	credentials, errCredentials := initCredentialsFromFile(afero.NewOsFs(), file, "https://www.googleapis.com/auth/cloud-platform")
	if errCredentials != nil {
		return "", errCredentials
	}

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	res, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: s.ResourceId,
	})

	if err != nil {
		return "", err
	}

	return string(res.Payload.GetData()), nil
}

func initCredentialsFromFile(fs afero.Fs, fileName string, scopes ...string) (*google.Credentials, error) {

	ctx := context.Background()

	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := google.CredentialsFromJSON(ctx, data, scopes...)
	if err != nil {
		log.Fatal(err)
	}

	return creds, err

}