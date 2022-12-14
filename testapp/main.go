package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stiggio/api-client-go"
	"net/http"
	"strings"
)

func extractApiKeyReq(req *http.Request) (apiKey string, err error) {
	apiKey = req.Header.Get("X-API-Key")
	//trimming so the bash getCustomerByRefId demo app script is happy
	apiKey = strings.Trim(apiKey, "\"")
	if apiKey == "" {
		return "", errors.New("no apiKey provided")
	}

	return apiKey, nil
}

func getClient(req *http.Request) (stigg.StiggGraphQLClient, error) {
	apiKey, err := extractApiKeyReq(req)
	if err != nil {
		return nil, errors.New("api key header was not provided")
	}
	stiggGraphUrl := req.URL.Query().Get("stigg-graph-url")
	client := stigg.NewStiggClient(apiKey, nil, &stiggGraphUrl)

	return client, nil
}
func getCustomerByRefId(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var input stigg.GetCustomerByRefIDInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	customer, err := client.GetCustomerByID(ctx, input)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(customer)
	w.Write(bytesResp)

}
func getCoupons(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	resp, err := client.GetCoupons(ctx)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp)
	w.Write(bytesResp)
}
func getPaywall(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.GetPaywallInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	coupons, err := client.GetPaywall(ctx, input)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(coupons)
	w.Write(bytesResp)

	return
}
func getEntitlements(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.FetchEntitlementsQuery
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.GetEntitlements(ctx, input)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp)
	w.Write(bytesResp)

	return
}
func getEntitlement(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.FetchEntitlementQuery
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.GetEntitlement(ctx, input)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp)
	w.Write(bytesResp)

	return
}
func createCustomer(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.CustomerInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.CreateCustomer(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.CreateCustomer)
	w.Write(bytesResp)

	return
}
func importCustomer(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.ImportCustomerInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.ImportCustomer(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.ImportCustomer)
	w.Write(bytesResp)

	return
}
func createSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.SubscriptionInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	subscription, err := client.CreateSubscription(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(subscription.CreateSubscription)
	w.Write(bytesResp)

	return
}
func provisionCustomer(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.ProvisionCustomerInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.ProvisionCustomer(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.ProvisionCustomer)
	w.Write(bytesResp)
}
func provisionSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.ProvisionSubscriptionInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	subscription, err := client.ProvisionSubscription(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(subscription.ProvisionSubscriptionV2)
	w.Write(bytesResp)

	return
}
func updateSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.UpdateSubscriptionInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	subscription, err := client.UpdateSubscription(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(subscription.UpdateSubscription)
	w.Write(bytesResp)

	return
}
func updateCustomer(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.UpdateCustomerInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.UpdateCustomer(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.UpdateCustomer)
	w.Write(bytesResp)

	return
}
func cancelSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.SubscriptionCancellationInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.CancelSubscription(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.CancelSubscription)
	w.Write(bytesResp)

	return
}
func initiateCheckout(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.InitiateCheckoutInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.InitiateCheckout(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.InitiateCheckout)
	w.Write(bytesResp)

	return
}
func estimateSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.EstimateSubscriptionInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.EstimateSubscription(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.EstimateSubscription)
	w.Write(bytesResp)

	return
}
func estimateSubscriptionUpdate(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.EstimateSubscriptionUpdateInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.EstimateSubscriptionUpdate(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.EstimateSubscriptionUpdate)
	w.Write(bytesResp)

	return
}
func reportUsage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := getClient(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var input stigg.UsageMeasurementCreateInput
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&input); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't marshal request\n%+v", err)))
		return
	}

	resp, err := client.ReportUsage(ctx, input)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Couldn't get customer\n%+v", err)))
		return
	}
	bytesResp, _ := json.Marshal(resp.CreateUsageMeasurement)
	w.Write(bytesResp)

	return
}
func main() {
	//queries
	http.HandleFunc("/getCustomerByRefId", getCustomerByRefId)
	http.HandleFunc("/getCoupons", getCoupons)
	http.HandleFunc("/getPaywall", getPaywall)
	http.HandleFunc("/getEntitlements", getEntitlements)
	http.HandleFunc("/getEntitlement", getEntitlement)
	// mutations
	http.HandleFunc("/createCustomer", createCustomer)
	http.HandleFunc("/provisionCustomer", provisionCustomer)
	http.HandleFunc("/importCustomer", importCustomer)
	http.HandleFunc("/updateCustomer", updateCustomer)
	http.HandleFunc("/createSubscription", createSubscription)
	http.HandleFunc("/provisionSubscription", provisionSubscription)
	http.HandleFunc("/updateSubscription", updateSubscription)
	http.HandleFunc("/cancelSubscription", cancelSubscription)
	http.HandleFunc("/initiateCheckout", initiateCheckout)
	http.HandleFunc("/estimateSubscription", estimateSubscription)
	http.HandleFunc("/estimateSubscriptionUpdate", estimateSubscriptionUpdate)
	http.HandleFunc("/reportUsage", reportUsage)

	port := "6666"
	fmt.Printf("Starting server at port ", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Errorf("cant start server", err)
	}
}
