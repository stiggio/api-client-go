![img_2.png](img_2.png)

# Stigg Go client.

## Getting started

0. Install

```go
go get github.com/stiggio/api-client-go/v5
```

1. Retrieve your server API key from [Stigg console](https://app.stigg.io/account/settings)
2. Init your Go SDK

```go
apiKey := "XXXXXXXXX"
client := stigg.NewStiggClient(apiKey, nil, nil)
```

For full usage specification, see our [Go Sdk documentation](https://docs.stigg.io/docs/go-sdk)
