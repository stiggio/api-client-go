package stigg

import (
	"net/http"
)

const (
	defaultApiUrl     = "https://api.stigg.io/graphql"
	defaultEdgeApiUrl = "https://edge.api.stigg.io"
)

// ClientConfig carries the optional settings of NewStiggClient.
type ClientConfig struct {
	// ApiUrl overrides the GraphQL API endpoint. Empty means https://api.stigg.io/graphql.
	//
	// ApiUrl and EdgeApiUrl default independently, so pointing ApiUrl at a non-production
	// deployment without also setting EdgeApiUrl leaves Edge-eligible reads on production.
	ApiUrl string
	// EdgeApiUrl overrides the Edge API endpoint. Empty means https://edge.api.stigg.io.
	EdgeApiUrl string
	// DisableEdge routes every operation to the GraphQL API.
	DisableEdge bool
	// HttpClient is shared by the GraphQL and Edge transports. Nil means a retryable default with
	// a 30 second timeout and RetryCount retries.
	HttpClient *http.Client
}
