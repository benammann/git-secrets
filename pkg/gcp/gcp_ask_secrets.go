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

type SecretsList struct {
	parentProject *cloudresourcemanager.Project
	secrets []*secretmanagerpb.Secret
}


func ListSecrets(ctx context.Context, parentProject *cloudresourcemanager.Project) (*SecretsList, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	projectsList := client.ListSecrets(ctx, &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%d", parentProject.ProjectNumber),
	})

	var secrets []*secretmanagerpb.Secret

	for {
		resp, err := projectsList.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not fetch secret version: %s", err.Error())
		}
		secrets = append(secrets, resp)
	}

	if len(secrets) < 1 {
		return nil, fmt.Errorf("there are no secrets configured in %s", parentProject.Name)
	}

	return NewSecretsList(parentProject, secrets), nil

}

func NewSecretsList(parentProject *cloudresourcemanager.Project, secrets []*secretmanagerpb.Secret) *SecretsList {
	return &SecretsList{
		parentProject: parentProject,
		secrets: secrets,
	}
}

func (l *SecretsList) GetSecretByName(name string) *secretmanagerpb.Secret {
	for _, secret := range l.secrets {
		if secret.Name == name {
			return secret
		}
	}
	return nil
}

func (l *SecretsList) SecretNames() (names []string) {
	for _, secret := range l.secrets {
		names = append(names, secret.Name)
	}
	return names
}

func (l *SecretsList) AskSecret() (*secretmanagerpb.Secret, error) {
	secretName := ""
	prompt := &survey.Select{
		Message: "Secret Name:",
		Options: l.SecretNames(),
	}
	errAsk := survey.AskOne(prompt, &secretName)
	return l.GetSecretByName(secretName), errAsk
}