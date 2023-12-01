package main

import (
	"fmt"
	"strconv"
	"time"

	xenditgo "github.com/grosenia/xendit-go-client"
	viper "github.com/spf13/viper"
)

var xenditclient xenditgo.Client
var qrCodeGateway xenditgo.QrcodeGateway

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
	referenceId := "ORDER_15263"
	typeQr := "DYNAMIC"
	currency := "IDR"
	amount := 256421.00
	expiresAt := "2023-12-02T09:56:43.60445Z"
	// shouldExcludeCreditCard := true

	fmt.Println("Generated Order ID: " + referenceId)

	var qrRequest = &xenditgo.XenditCreteQrcodeReq{
		ReferenceId: referenceId,
		Type:        typeQr,
		Amount:      amount,
		Currency:    currency,
		ExpiredAt:   expiresAt,
	}

	resp, err := qrCodeGateway.CreateQrCode(qrRequest)
	fmt.Println("Error Code: " + resp.ErrorCode)
	fmt.Println("Error Message: " + resp.ErrorMessage)

	//simulate payment
	fmt.Println("respon id: " + resp.Id)
	qrId := resp.Id
	fmt.Println("respon id: " + qrId)
	resp2, err := qrCodeGateway.PayQrCode(qrId)
	fmt.Println("Error Code: " + resp2.ErrorCode)
	fmt.Println("Error Message: " + resp2.ErrorMessage)

}

func setupClient() {
	xenditclient = xenditgo.NewClient()

	// TODO: should put in config file
	xenditclient.SecretAPIKey = viper.GetString("KEY_WRITE_MONEY_IN")
	xenditclient.ApiVersion = viper.GetString("API_VERSION")
	xenditclient.APIEnvType = xenditgo.Sandbox

	qrCodeGateway = xenditgo.QrcodeGateway{
		Client: xenditclient,
	}

}

func generateOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
