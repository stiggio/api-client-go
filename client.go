package stigg

import (
	"context"
	"fmt"
	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"time"
)

const RetryCount = 3

func NewStiggClient(apiKey string, httpClient *http.Client, baseUrl *string) StiggClient {
	if baseUrl == nil {
		defaultStiggBaseUrl := "https://api.stigg.io/graphql"
		baseUrl = &defaultStiggBaseUrl
	}
	if httpClient == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = RetryCount
		retryClient.Logger = nil

		httpClient = retryClient.StandardClient()
		httpClient.Timeout = time.Second * 30
	}

	return NewClient(httpClient, *baseUrl, func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("x-api-key", fmt.Sprintf(apiKey))
		req.Header.Set("Cache-Control", "no-cache")
		return next(ctx, req, gqlInfo, res)
	})
}
