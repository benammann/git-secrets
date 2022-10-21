package config_generic

import (
	"context"
	"github.com/benammann/git-secrets/pkg/gcp"
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

	// resolvedPayload holds the resolved secret payload and acts as in memory cache
	resolvedPayload string

}

func (s *GcpSecret) GetType() string {
	return "gcp"
}

func (s *GcpSecret) GetName() string {
	return s.Name
}

func (s *GcpSecret) GetOriginContext() *Context {
	return s.OriginContext
}

func (s *GcpSecret) GetPlainValue(ctx context.Context) (string, error) {

	if s.resolvedPayload != "" {
		return s.resolvedPayload, nil
	}

	file := s.OriginContext.GlobalConfig.GetGcpCredentialsFile(s.OriginContext.GcpCredentials)

	resolvedSecret, errResolve := gcp.ResolveSecret(ctx, s.ResourceId, file)

	if errResolve == nil {
		s.resolvedPayload = resolvedSecret
	}
	return resolvedSecret, errResolve

}