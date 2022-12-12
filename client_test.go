package stigg

import (
	"context"
	stigg "github.com/stiggio/api-client-go/codegen/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	stiggUrl := "http://localhost:4000/graphql"
	t.Run("Test client startup", func(t *testing.T) {
		_, err := NewClient("$2b$10$GD4f3l9udQYfipImpm1eWrOP7fENFQiiybUyBvckgryVu1j.3Zgq8u:d62b64ac-0ccf-410c-bb2f-05474e2152fe", &stiggUrl)

		if err != nil {
			t.Error("couldn't start client")
		}
	})

	t.Run("Test getCustomerByID", func(t *testing.T) {

		client, err := NewClient("$2b$10$GD43l9udQYfipImpm1eWrOP7fENFQiiybUyBvckgryVu1j.3Zgq8u:d62b64ac-0ccf-410c-bb2f-05474e2152fe", &stiggUrl)

		if err != nil {
			t.Error("couldn't start client")
		}
		customerID := "customer-demo-0111"

		input := stigg.GetCustomerByRefIDInput{CustomerID: customerID}
		customer, err := client.GetCustomerByID(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, customerID, customer.GetCustomerByRefID.ID, "Couldn't get customer")

	})
}
