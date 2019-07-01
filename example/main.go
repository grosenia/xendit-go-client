package main

import (
	"fmt"
	"strconv"
	"time"

	xenditgo "github.com/grosenia/xendit-go-client"
)

var xenditclient xenditgo.Client

func main() {
	fmt.Println("Setup client")
	setupClient()

}

func setupClient() {
	xenditclient = xenditgo.NewClient()
	xenditclient.SecretAPIKey = "xnd_development_0xnxnk0EjMKa3YB78pQCeD6DrtGSyVJqDGxFonptcTCjmlM1UUsLlaYnZ4iuISA"
	xenditclient.APIEnvType = xenditgo.Sandbox

	/*
		InvoiceGateway = midtrans.InvoiceGateway{
			Client: midclient,
		}

		snapGateway = midtrans.SnapGateway{
			Client: midclient,
		} */
}

func generateOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
