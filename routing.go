package stigg

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
)

// edgeRoutingClient serves the Edge-eligible read operations from the Edge API and delegates every
// other operation to the embedded GraphQL client.
//
// The routing rules mirror node-server-sdk: a customer-scoped paywall cannot be served from the Edge
// API (src/api/managementApi.ts:689), and every other Edge-eligible read goes to the Edge API
// unconditionally (src/api/managementApi.ts:407, src/api/entitlements/entitlementsApi.ts:59).
//
// Per-call interceptors are GraphQL request interceptors, so they do not apply on the Edge path.
// Use ClientConfig.DisableEdge for an operation that has to run through them.
type edgeRoutingClient struct {
	StiggClient
	edge *edgeApiClient
}

func (c *edgeRoutingClient) GetPaywall(ctx context.Context, input GetPaywallInput, interceptors ...clientv2.RequestInterceptor) (*GetPaywall, error) {
	if input.CustomerID != nil {
		return c.StiggClient.GetPaywall(ctx, input, interceptors...)
	}

	return c.edge.GetPaywall(ctx, input)
}

func (c *edgeRoutingClient) GetEntitlements(ctx context.Context, query FetchEntitlementsQuery, interceptors ...clientv2.RequestInterceptor) (*GetEntitlements, error) {
	return c.edge.GetEntitlements(ctx, query)
}

func (c *edgeRoutingClient) GetEntitlementsState(ctx context.Context, query FetchEntitlementsQuery, interceptors ...clientv2.RequestInterceptor) (*GetEntitlementsState, error) {
	return c.edge.GetEntitlementsState(ctx, query)
}

func (c *edgeRoutingClient) GetEntitlement(ctx context.Context, query FetchEntitlementQuery, interceptors ...clientv2.RequestInterceptor) (*GetEntitlement, error) {
	return c.edge.GetEntitlement(ctx, query)
}

func (c *edgeRoutingClient) GetActiveSubscriptionsList(ctx context.Context, input GetActiveSubscriptionsInput, interceptors ...clientv2.RequestInterceptor) (*GetActiveSubscriptionsList, error) {
	return c.edge.GetActiveSubscriptionsList(ctx, input)
}
