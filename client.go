package stigg

import (
	"context"
	"fmt"
	"github.com/Yamashou/gqlgenc/clientv2"
	stigg "github.com/stiggio/api-client-go/codegen/generated"
	"net/http"
)

func NewClient(sdkKey string, baseUrl *string) (*stigg.Client, error) {
	if baseUrl == nil {
		defaultStiggBaseUrl := "https://api.stigg.io/graphql"
		baseUrl = &defaultStiggBaseUrl
	}

	c := &stigg.Client{
		Client: clientv2.NewClient(http.DefaultClient, *baseUrl, func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
			req.Header.Set("x-api-key", fmt.Sprintf(sdkKey))
			req.Header.Set("Cache-Control", "no-cache")
			return next(ctx, req, gqlInfo, res)
		}),
	}
	return c, nil
}
