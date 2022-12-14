package stigg

import (
	"context"
	"fmt"
	"github.com/Yamashou/gqlgenc/clientv2"
	"net/http"
	"time"
)

func NewStiggClient(apiKey string, httpClient *http.Client, baseUrl *string) StiggGraphQLClient {
	if baseUrl == nil {
		defaultStiggBaseUrl := "https://api.stigg.io/graphql"
		baseUrl = &defaultStiggBaseUrl
	}
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	return NewClient(httpClient, *baseUrl, func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("x-api-key", fmt.Sprintf(apiKey))
		req.Header.Set("Cache-Control", "no-cache")
		return next(ctx, req, gqlInfo, res)
	})

}
