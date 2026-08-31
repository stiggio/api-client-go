package stigg

import (
	"context"
	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"time"
)

const RetryCount = 3

func NewStiggClient(apiKey string, config ...ClientConfig) StiggClient {
	var cfg ClientConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = defaultApiUrl
	}
	if cfg.EdgeApiUrl == "" {
		cfg.EdgeApiUrl = defaultEdgeApiUrl
	}
	if cfg.HttpClient == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = RetryCount
		retryClient.Logger = nil

		cfg.HttpClient = retryClient.StandardClient()
		cfg.HttpClient.Timeout = time.Second * 30
	}

	graphqlClient := NewClient(cfg.HttpClient, cfg.ApiUrl, func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("x-graphql-operation-name", gqlInfo.Request.OperationName)

		return next(ctx, req, gqlInfo, res)
	})

	if cfg.DisableEdge {
		return graphqlClient
	}

	return &edgeRoutingClient{
		StiggClient: graphqlClient,
		edge: &edgeApiClient{
			apiKey:     apiKey,
			baseUrl:    cfg.EdgeApiUrl,
			httpClient: cfg.HttpClient,
		},
	}
}
