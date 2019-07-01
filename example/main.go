package main

import (
	"fmt"
	"strconv"
	"time"

	xenditgo "github.com/grosenia/xendit-go-client"
	viper "github.com/spf13/viper"
)

var xenditclient xenditgo.Client
var invoiceGateway xenditgo.InvoiceGateway

func main() {

	fmt.Println("Load Config...")

	viper.SetConfigType("props")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	fmt.Println("Load Config success")
	fmt.Println("Setup client")

	setupClient()

	// Example
	orderID := generateOrderID()
	payerEmail := "boylevantz@gmail.com"
	amount := 62000.00
	description := "Handuk 2 pcs"
	shouldSendEmail := true
	shouldExcludeCreditCard := true

	// This VA ID is already created before
	// TODO: should take out the Account ID to ENV

	callbackVirtualAccountID := viper.GetString("SAMPLE_VA_ACCOUNT_ID")

	invoiceDuration := 86400

	fmt.Println("Generated Order ID: " + orderID)
	fmt.Println("Call back VA ID: " + callbackVirtualAccountID)

	var invoiceRequest = &xenditgo.XenditCreateInvoiceReq{
		ExternalID:               orderID,
		PayerEmail:               payerEmail,
		Amount:                   amount,
		Description:              description,
		ShouldSendEmail:          shouldSendEmail,
		ShouldExcludeCreditCard:  shouldExcludeCreditCard,
		CallbackVirtualAccountID: callbackVirtualAccountID,
		InvoiceDuration:          invoiceDuration,
	}

	resp, err := invoiceGateway.CreateInvoice(invoiceRequest)
	fmt.Println("Invoice URL: " + resp.InvoiceURL)

}

func setupClient() {
	xenditclient = xenditgo.NewClient()

	// TODO: should put in config file
	xenditclient.SecretAPIKey = viper.GetString("KEY_WRITE_MONEY_IN")
	xenditclient.APIEnvType = xenditgo.Sandbox

	invoiceGateway = xenditgo.InvoiceGateway{
		Client: xenditclient,
	}

}

func generateOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
