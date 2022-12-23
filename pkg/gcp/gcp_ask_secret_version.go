package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

type SecretsVersionList struct {
	parentProject *cloudresourcemanager.Project
	parentSecret *secretmanagerpb.Secret
	versions []*secretmanagerpb.SecretVersion
}


func ListSecretVersions(ctx context.Context, parentProject *cloudresourcemanager.Project, parentSecret *secretmanagerpb.Secret) (*SecretsVersionList, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	projectsList := client.ListSecretVersions(ctx, &secretmanagerpb.ListSecretVersionsRequest{
		Parent: parentSecret.Name,
	})

	var versions []*secretmanagerpb.SecretVersion

	for {
		resp, err := projectsList.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not fetch secret versions: %s", err.Error())
		}
		versions = append(versions, resp)
	}

	if len(versions) < 1 {
		return nil, fmt.Errorf("there are no versions configured in %s", parentProject.Name)
	}

	return NewSecretVersionList(parentProject, parentSecret, versions), nil

}

func NewSecretVersionList(parentProject *cloudresourcemanager.Project, parentSecret *secretmanagerpb.Secret, versions []*secretmanagerpb.SecretVersion) *SecretsVersionList {
	return &SecretsVersionList{
		parentProject: parentProject,
		parentSecret: parentSecret,
		versions: versions,
	}
}

func (l *SecretsVersionList) GetVersionByName(name string) *secretmanagerpb.SecretVersion {
	for _, version := range l.versions {
		if version.Name == name {
			return version
		}
	}
	return nil
}

func (l *SecretsVersionList) VersionNames() (names []string) {
	for _, version := range l.versions {
		names = append(names, version.Name)
	}
	return names
}

func (l *SecretsVersionList) AskSecretVersion() (*secretmanagerpb.SecretVersion, error) {
	secretName := ""
	prompt := &survey.Select{
		Message: "Secret Version:",
		Options: l.VersionNames(),
	}
	errAsk := survey.AskOne(prompt, &secretName)
	return l.GetVersionByName(secretName), errAsk
}