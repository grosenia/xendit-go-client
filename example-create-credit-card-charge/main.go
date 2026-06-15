package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	fmt.Println("Load Config success")
	fmt.Println("Setup client")

	setupClient()

	tokenID := viper.GetString("TOKEN_ID")
	authenticationID := viper.GetString("AUTHENTICATION_ID")
	amount := viper.GetFloat64("CHARGE_AMOUNT")
	externalID := viper.GetString("EXTERNAL_ID")

	// Optional CLI override: token_id [authentication_id] [amount] [external_id]
	args := os.Args[1:]
	if len(args) >= 1 && args[0] != "" {
		tokenID = args[0]
	}
	if len(args) >= 2 && args[1] != "" {
		authenticationID = args[1]
	}
	if len(args) >= 3 && args[2] != "" {
		parsedAmount, parseErr := strconv.ParseFloat(args[2], 64)
		if parseErr != nil {
			fmt.Println("Invalid amount:", parseErr)
			return
		}
		amount = parsedAmount
	}
	if len(args) >= 4 && args[3] != "" {
		externalID = args[3]
	}

	if tokenID == "" {
		fmt.Println("TOKEN_ID is required.")
		fmt.Println("Generate one with tokenize.html (sandbox), then set TOKEN_ID in config.props")
		fmt.Println("Usage: go run main.go [token_id] [authentication_id] [amount] [external_id]")
		return
	}
	if amount <= 0 {
		fmt.Println("CHARGE_AMOUNT must be greater than 0")
		return
	}
	if externalID == "" {
		externalID = generateOrderID()
	}

	fmt.Println("External ID:", externalID)
	fmt.Println("Amount:", amount)
	fmt.Println("Token ID:", tokenID)
	if authenticationID != "" {
		fmt.Println("Authentication ID:", authenticationID)
	}

	chargeRequest := &xenditgo.XenditCreateCreditCardChargeReq{
		TokenID:          tokenID,
		ExternalID:       externalID,
		Amount:           amount,
		AuthenticationID: authenticationID,
		Capture:          true,
	}

	resp, err := invoiceGateway.CreateCreditCardCharge(chargeRequest)
	if err != nil {
		fmt.Println("Error creating credit card charge:", err)
		return
	}

	if resp.ErrorStatus {
		fmt.Println("Error:", resp.Error())
		return
	}

	pretty, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("Credit card charge response:")
	fmt.Println(string(pretty))

	if resp.PayerAuthenticationURL != "" {
		fmt.Println()
		fmt.Println("3DS required — open this URL in a browser to complete authentication:")
		fmt.Println(resp.PayerAuthenticationURL)
	}
}

func setupClient() {
	xenditclient = xenditgo.NewClient()
	xenditclient.SecretAPIKey = viper.GetString("KEY_WRITE_MONEY_IN")
	xenditclient.APIEnvType = xenditgo.Sandbox
	xenditclient.LogLevel = 3

	invoiceGateway = xenditgo.InvoiceGateway{
		Client: xenditclient,
	}
}

func generateOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
