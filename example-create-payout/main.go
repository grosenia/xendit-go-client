package main

import (
	"fmt"
	"strconv"
	"os"
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

	argsWithoutProg := os.Args[1:]

	// Example
	args1 := argsWithoutProg[0]
	args2 := argsWithoutProg[1]

	externalID := args1
	amount, err := strconv.ParseFloat(args2, 64)
	if err != nil {
		return
	}
	
	fmt.Println("External ID: " + externalID)
	fmt.Println("Amount: ", amount)

	var payoutRequest = &xenditgo.XenditCreatePayoutReq{
		ExternalID: externalID,
		Amount:     amount,
	}

	resp, err := invoiceGateway.CreatePayout(payoutRequest)

	fmt.Println("Response: %v", resp)
	fmt.Println("Payout URL: " + resp.PayoutURL)
	fmt.Println("Pass Code: " + resp.PassCode)
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
