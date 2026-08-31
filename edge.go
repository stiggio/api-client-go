package stigg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Yamashou/gqlgenc/clientv2"
)

const (
	edgeApiVersionHeader = "X-API-VERSION"
	edgeApiVersionValue  = "1"
)

// edgeApiClient reads the Edge API, a REST-style JSON layer that returns GraphQL-shaped envelopes
// for the subset of read operations it can serve without customer-specific resolution.
type edgeApiClient struct {
	apiKey     string
	baseUrl    string
	httpClient *http.Client
}

func (c *edgeApiClient) GetPaywall(ctx context.Context, input GetPaywallInput) (*GetPaywall, error) {
	params := url.Values{}
	if input.BillingCountryCode != nil {
		params.Set("billingCountryCode", *input.BillingCountryCode)
	}
	if input.FetchAllCountriesPrices != nil {
		params.Set("fetchAllCountriesPrices", strconv.FormatBool(*input.FetchAllCountriesPrices))
	}
	if input.IncludeHiddenPlans != nil {
		params.Set("includeHiddenPlans", strconv.FormatBool(*input.IncludeHiddenPlans))
	}

	path := "/v1/paywall.json"
	if input.ProductID != nil {
		path = fmt.Sprintf("/v1/p/%s/paywall.json", url.PathEscape(*input.ProductID))
	}

	var res GetPaywall
	if err := c.get(ctx, path, params, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *edgeApiClient) GetEntitlements(ctx context.Context, query FetchEntitlementsQuery) (*GetEntitlements, error) {
	params := url.Values{}
	if query.ResourceID != nil {
		params.Set("resourceId", *query.ResourceID)
	}

	var res GetEntitlements
	if err := c.get(ctx, c.customerPath(query.CustomerID, "entitlements.json"), params, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *edgeApiClient) GetEntitlementsState(ctx context.Context, query FetchEntitlementsQuery) (*GetEntitlementsState, error) {
	params := url.Values{}
	if query.ResourceID != nil {
		params.Set("resourceId", *query.ResourceID)
	}

	var res GetEntitlementsState
	if err := c.get(ctx, c.customerPath(query.CustomerID, "entitlements-state.json"), params, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEntitlement shares the entitlements.json endpoint with GetEntitlements; the featureId param is
// what narrows the response to a single entitlement.
func (c *edgeApiClient) GetEntitlement(ctx context.Context, query FetchEntitlementQuery) (*GetEntitlement, error) {
	params := url.Values{}
	params.Set("featureId", query.FeatureID)
	if query.ResourceID != nil {
		params.Set("resourceId", *query.ResourceID)
	}
	if query.Options != nil {
		if query.Options.RequestedUsage != nil {
			usage := *query.Options.RequestedUsage
			// The GraphQL path rejects these in encoding/json; without this they reach the URL as "NaN".
			if math.IsNaN(usage) || math.IsInf(usage, 0) {
				return nil, fmt.Errorf("requestedUsage must be finite, got %v", usage)
			}
			params.Set("requestedUsage", strconv.FormatFloat(usage, 'f', -1, 64))
		}
		for _, requestedValue := range query.Options.RequestedValues {
			params.Add("requestedValues", requestedValue)
		}
	}

	var res GetEntitlement
	if err := c.get(ctx, c.customerPath(query.CustomerID, "entitlements.json"), params, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *edgeApiClient) GetActiveSubscriptionsList(ctx context.Context, input GetActiveSubscriptionsInput) (*GetActiveSubscriptionsList, error) {
	params := url.Values{}
	if input.ResourceID != nil {
		params.Set("resourceId", *input.ResourceID)
	} else {
		for _, resourceID := range input.ResourceIds {
			params.Add("resourceIds", resourceID)
		}
	}

	var res GetActiveSubscriptionsList
	if err := c.get(ctx, c.customerPath(input.CustomerID, "subscriptions.json"), params, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *edgeApiClient) customerPath(customerID, resource string) string {
	return fmt.Sprintf("/v1/c/%s/%s", url.PathEscape(customerID), resource)
}

func (c *edgeApiClient) get(ctx context.Context, path string, params url.Values, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+path, nil)
	if err != nil {
		return fmt.Errorf("create request struct failed: %w", err)
	}
	req.URL.RawQuery = params.Encode()

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set(edgeApiVersionHeader, edgeApiVersionValue)
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return parseEdgeResponse(body, resp.StatusCode, res)
}

// Mirrors clientv2's unexported response parsing so an Edge failure surfaces to callers as the same
// *clientv2.ErrorResponse a GraphQL failure would.
func parseEdgeResponse(body []byte, httpCode int, result interface{}) error {
	errResponse := &clientv2.ErrorResponse{}
	isKOCode := httpCode < 200 || 299 < httpCode
	if isKOCode {
		errResponse.NetworkError = &clientv2.HTTPError{
			Code:    httpCode,
			Message: fmt.Sprintf("Response body %s", string(body)),
		}
	}

	if err := unmarshalEdgeResponse(body, result); err != nil {
		var gqlErrors *clientv2.GqlErrorList
		if errors.As(err, &gqlErrors) {
			errResponse.GqlErrors = &gqlErrors.Errors
		} else if !isKOCode {
			return err
		}
	}

	if errResponse.HasErrors() {
		return errResponse
	}

	return nil
}

func unmarshalEdgeResponse(data []byte, res interface{}) error {
	var resp struct {
		Data   json.RawMessage `json:"data"`
		Errors json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(data), err)
	}

	// An explicit "errors": null is an empty error list, not a failure.
	if len(resp.Errors) > 0 && string(resp.Errors) != "null" {
		gqlErrors := &clientv2.GqlErrorList{}
		if err := json.Unmarshal(data, gqlErrors); err != nil {
			return fmt.Errorf("failed to parse graphql errors. Response content %s - %w", string(data), err)
		}

		return gqlErrors
	}

	if len(resp.Data) == 0 {
		return fmt.Errorf("response contained neither data nor errors: %s", string(data))
	}

	// Not encoding/json: union fragments carry only `graphql:"... on X"` tags and would silently
	// decode as zero values.
	if err := unmarshalEdgeData(resp.Data, res); err != nil {
		return fmt.Errorf("failed to decode data into response: %w", err)
	}

	return nil
}
