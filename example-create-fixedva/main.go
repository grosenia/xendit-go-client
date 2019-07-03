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
	externalID := generateOrderID()
	bankCode := "BCA" // BCA / BNI / MANDIRI / PERMATA / BRI
	name := "Ali Irawan"

	var createFixedVaRequest = &xenditgo.XenditCreateFixedVaReq{
		ExternalID: externalID,
		BankCode:   bankCode,
		Name:       name,
	}

	resp, err := invoiceGateway.CreateFixedVa(createFixedVaRequest)

	fmt.Println("Created fixed VA Response: ")
	fmt.Println(resp)
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
