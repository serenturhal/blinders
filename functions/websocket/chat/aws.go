package main

import (
	"context"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	agm "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type CustomEndpointResolve struct {
	Domain, Stage string
}

func (c CustomEndpointResolve) ResolveEndpoint(_ context.Context, _ agm.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	var endpoint url.URL
	endpoint.Scheme = "https"
	endpoint.Host = c.Domain
	endpoint.Path = "v1"

	return smithyendpoints.Endpoint{
		URI:     endpoint,
		Headers: http.Header{},
	}, nil
}

// NewAPIGatewayManagementClient creates a new API Gateway Management Client instance
// from the provided parameters. The new client will have a custom endpoint
// that resolves to the application's deployed API.
func NewAPIGatewayManagementClient(cfg *aws.Config, domain, stage string) *agm.Client {
	cer := CustomEndpointResolve{
		Domain: domain,
		Stage:  stage,
	}
	return agm.NewFromConfig(*cfg,
		agm.WithEndpointResolverV2(cer))
}

func Publish(ctx context.Context, id string, data []byte) error {
	_, err := apiClient.PostToConnection(ctx, &agm.PostToConnectionInput{
		ConnectionId: &id,
		Data:         data,
	})

	return err
}
