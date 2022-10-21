package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/benammann/git-secrets/pkg/utility"
	"github.com/spf13/afero"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

const scopeGoogleCloudPlatform = "https://www.googleapis.com/auth/cloud-platform"

type ResourceId string
type CredentialsFile string

var fs = afero.NewOsFs()
var smClients = make(map[CredentialsFile]*secretmanager.Client)
var resolvedSecrets = make(map[ResourceId]string)

func getClient(ctx context.Context, gcpCredentialsFile CredentialsFile) (*secretmanager.Client, error) {

	credentials, errCredentials := initCredentialsFromFile(fs, gcpCredentialsFile, scopeGoogleCloudPlatform)
	if errCredentials != nil {
		return nil, errCredentials
	}

	if smClients[gcpCredentialsFile] != nil{
		return smClients[gcpCredentialsFile], nil
	}

	newClient, err := secretmanager.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	smClients[gcpCredentialsFile] = newClient

	addChannel, doneChannel := utility.GetContextChannels(ctx)

	addChannel<-1

	go func() {
		select {
		case <-ctx.Done():
			errClose := newClient.Close()
			delete(smClients, gcpCredentialsFile)
			doneChannel <-errClose==nil
		}
	}()

	return newClient, nil

}

func ListSecrets(ctx context.Context, gcpCredentialsFile string) ([]*secretmanagerpb.SecretVersion, error) {

	smClient, errClient := getClient(ctx, CredentialsFile(gcpCredentialsFile))
	if errClient != nil {
		return nil, fmt.Errorf("could not initialize secret manager client for %s: %s", gcpCredentialsFile, errClient.Error())
	}

	it := smClient.ListSecretVersions(ctx, &secretmanagerpb.ListSecretVersionsRequest{
		Parent: "projects/806001934377/secrets/awesomeSecret",
	})

	var secretVersions []*secretmanagerpb.SecretVersion

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not fetch secret version: %s", err.Error())
		}
		secretVersions = append(secretVersions, resp)
	}

	if len(secretVersions) < 1 {
		return nil, fmt.Errorf("could not fetch any secret versions from %s", gcpCredentialsFile)
	}

	return secretVersions, nil

}

func ResolveSecret(ctx context.Context, resourceId string, gcpCredentialsFile string) (string, error) {

	cacheKey := ResourceId(resourceId)

	if resolvedSecrets[cacheKey] != "" {
		return resolvedSecrets[cacheKey], nil
	}

	smClient, errClient := getClient(ctx, CredentialsFile(gcpCredentialsFile))
	if errClient != nil {
		return "", fmt.Errorf("could not initialize secret manager client for %s: %s", gcpCredentialsFile, errClient.Error())
	}

	res, errResolve := smClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceId,
	})
	if errResolve != nil {
		return "", fmt.Errorf("could not resolve secret %s: %s", resourceId, errResolve.Error())
	}

	resolvedSecrets[cacheKey] = string(res.Payload.GetData())

	return resolvedSecrets[cacheKey], nil

}

func initCredentialsFromFile(fs afero.Fs, fileName CredentialsFile, scopes ...string) (*google.Credentials, error) {

	ctx := context.Background()

	data, err := afero.ReadFile(fs, string(fileName))
	if err != nil {
		log.Fatal(err)
	}

	creds, err := google.CredentialsFromJSON(ctx, data, scopes...)
	if err != nil {
		log.Fatal(err)
	}

	return creds, err

}