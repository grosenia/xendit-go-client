package main

import (
	"fmt"
	"time"

	xenditgo "github.com/grosenia/xendit-go-client"
	viper "github.com/spf13/viper"
)

var xenditclient xenditgo.Client
var disbursementGateway xenditgo.DisbursementGateway

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

	// argsWithoutProg := os.Args[1:]

	// // Example
	// args1 := argsWithoutProg[0]
	// args2 := argsWithoutProg[1]

	// externalID := args1
	// amount, err := strconv.ParseFloat(args2, 64)
	// if err != nil {
	// 	return
	// }

	// fmt.Println("External ID: " + externalID)
	// fmt.Println("Amount: ", amount)

	var disbursements []xenditgo.DisbursementItem
	var disbursement xenditgo.DisbursementItem

	headerId := fmt.Sprintf("test_%d", time.Now().Unix())
	for i := 0; i < 3; i++ {
		itemID := fmt.Sprintf("%s_%d", headerId, i+1)
		disbursement.ExternalID = itemID
		disbursement.Amount = 20000
		disbursement.BankCode = "BCA"
		disbursement.BankAccountName = "Stanley Nguyen"
		disbursement.BankAccountNumber = "12345678"
		disbursement.Description = "Reimbursement for pair of shoes (1)"
		disbursements = append(disbursements, disbursement)
	}

	Reference := headerId
	key := headerId

	var batchRequest = &xenditgo.XenditCreateBatchReq{
		HeaderID:      Reference,
		Disbursements: disbursements,
	}

	resp, err := disbursementGateway.CreateBatchDisbursement(key, batchRequest)
	if err != nil {
		panic(fmt.Errorf("fatal error CreateBatchDisbursement: %s", err))
	}
	fmt.Println(fmt.Sprintf("Status: %v", resp.Status))
}

func setupClient() {
	xenditclient = xenditgo.NewClient()

	// TODO: should put in config file
	xenditclient.SecretAPIKey = viper.GetString("KEY_WRITE_MONEY_OUT")
	xenditclient.APIEnvType = xenditgo.Sandbox

	disbursementGateway = xenditgo.DisbursementGateway{
		Client: xenditclient,
	}
}
