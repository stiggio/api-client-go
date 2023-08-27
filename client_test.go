package stigg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	apiUrl := "https://api.stigg.io/graphql"
	apiKey := os.Getenv("SERVER_API_KEY")

	t.Run("Test NewStiggClient", func(t *testing.T) {
		c := NewStiggClient("", nil, &apiUrl)
		assert.NotNil(t, c)
	})

	t.Run("Test GetCustomerByID", func(t *testing.T) {
		client := NewStiggClient(apiKey, nil, &apiUrl)

		customerID := "customer-demo-01"
		input := GetCustomerByRefIDInput{CustomerID: customerID}
		customer, err := client.GetCustomerByID(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, customerID, customer.GetCustomerByRefID.RefID, "Couldn't get customer")
	})

	t.Run("Test ProvisionSubscription", func(t *testing.T) {
		client := NewStiggClient(apiKey, nil, &apiUrl)

		customerID := "customer-demo-01"
		billingPeriod := BillingPeriodMonthly
		unitQuantity := 5.0
		planID := "plan-revvenu-essentials"
		input := ProvisionSubscriptionInput{
			CustomerID:    customerID,
			PlanID:        planID,
			BillingPeriod: &billingPeriod,
			UnitQuantity:  &unitQuantity,
			CheckoutOptions: &CheckoutOptions{
				SuccessURL: "https://www.google.com/search?q=success",
				CancelURL:  "https://www.google.com/search?q=cancel",
			},
		}
		resp, err := client.ProvisionSubscription(ctx, input)

		assert.Nil(t, err)
		assert.Contains(t, *resp.ProvisionSubscription.CheckoutURL, "checkout.stripe.com", "Invalid checkout URL is missing")
	})
}
