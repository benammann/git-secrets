package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

type ResourceId string

var resolvedSecrets = make(map[ResourceId]string)

func ResolveSecret(ctx context.Context, resourceId string) (string, error) {

	cacheKey := ResourceId(resourceId)

	if resolvedSecrets[cacheKey] != "" {
		return resolvedSecrets[cacheKey], nil
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	res, errResolve := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceId,
	})
	if errResolve != nil {
		return "", fmt.Errorf("could not resolve secret %s: %s", resourceId, errResolve.Error())
	}

	resolvedSecrets[cacheKey] = string(res.Payload.GetData())

	return resolvedSecrets[cacheKey], nil

}