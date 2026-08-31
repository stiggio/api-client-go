![img_2.png](img_2.png)

# Stigg Go client.

## Getting started

0. Install

```go
go get github.com/stiggio/api-client-go/v7
```

1. Retrieve your server API key from [Stigg console](https://app.stigg.io/account/settings)
2. Init your Go SDK

```go
apiKey := "XXXXXXXXX"
client := stigg.NewStiggClient(apiKey)
```

For full usage specification, see our [Go SDK documentation](https://docs.stigg.io/docs/go-sdk)

## Upgrading

`NewStiggClient` takes the API key followed by an optional `ClientConfig`, replacing the three positional arguments:

```go
// before
client := stigg.NewStiggClient(apiKey, nil, nil)
client := stigg.NewStiggClient(apiKey, httpClient, baseUrl)

// after
client := stigg.NewStiggClient(apiKey)
client := stigg.NewStiggClient(apiKey, stigg.ClientConfig{
	HttpClient: httpClient, // optional, was the 2nd argument
	ApiUrl:     *baseUrl,   // optional, was the 3rd argument
})
```

`ApiUrl` is a plain `string` rather than the old `*string`; leave it empty for the default endpoint.

## Edge API routing

The following read operations are served from the [Edge API](https://docs.stigg.io/api-and-sdks/integration/graphql/edge-api-rate-limits) (`https://edge.api.stigg.io`) instead of the GraphQL API, which is significantly faster for cacheable reads:

| Operation | Edge endpoint |
| --- | --- |
| `GetPaywall` (public only) | `/v1[/p/{productId}]/paywall.json` |
| `GetEntitlements` | `/v1/c/{customerId}/entitlements.json` |
| `GetEntitlement` | `/v1/c/{customerId}/entitlements.json?featureId=…` |
| `GetEntitlementsState` | `/v1/c/{customerId}/entitlements-state.json` |
| `GetActiveSubscriptionsList` | `/v1/c/{customerId}/subscriptions.json` |

Routing is transparent — the request and response types are unchanged. Matching `node-server-sdk`, the only input that forces the GraphQL API is a `GetPaywall` with `CustomerID` set, since the Edge API serves only the public paywall. The other four operations always use the Edge API.

The Edge API rate limit is 6,000 RPM per API key, shared across all Edge endpoints.

### Configuration

Edge routing is on by default and needs no configuration:

```go
client := stigg.NewStiggClient(apiKey)
```

To override the Edge endpoint, or to send every operation to the GraphQL API:

```go
client := stigg.NewStiggClient(apiKey, stigg.ClientConfig{
	EdgeApiUrl:  "https://edge.api.stigg.io", // optional
	DisableEdge: false,                       // true sends everything to the GraphQL API
})
```

`ApiUrl` and `EdgeApiUrl` default independently, so pointing `ApiUrl` at a non-production deployment without also setting `EdgeApiUrl` leaves Edge-eligible reads on production.
