package stigg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	stiggUrl := "http://localhost:4000/graphql"
	t.Run("Test client startup", func(t *testing.T) {
		c := NewStiggClient("$2b$10$GD4f3l9udQYfipImpm1eWrOP7fENFQiiybUyBvckgryVu1j.3Zgq8u:d62b64ac-0ccf-410c-bb2f-05474e2152fe", nil, &stiggUrl)
		assert.NotNil(t, c)
	})

	t.Run("Test getCustomerByID", func(t *testing.T) {

		client := NewStiggClient("$2b$10$GD43l9udQYfipImpm1eWrOP7fENFQiiybUyBvckgryVu1j.3Zgq8u:d62b64ac-0ccf-410c-bb2f-05474e2152fe", nil, &stiggUrl)

		customerID := "customer-demo-0111"

		input := GetCustomerByRefIDInput{CustomerID: customerID}
		customer, err := client.GetCustomerByID(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, customerID, customer.GetCustomerByRefID.ID, "Couldn't get customer")
	})
}
