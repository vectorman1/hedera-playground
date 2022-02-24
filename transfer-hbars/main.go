package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"log"
)

func main() {
	hs := client.Setup2TestClients()

	user1 := hs.Users[0]
	user2 := hs.Users[1]
	trans := hedera.NewTransferTransaction().
		AddHbarTransfer(user1.AccountId, hedera.HbarFrom(-1000000000, hedera.HbarUnits.Tinybar)).
		AddHbarTransfer(user2.AccountId, hedera.HbarFrom(1000000000, hedera.HbarUnits.Tinybar))

	txResponse, err := trans.Execute(user1.C)
	if err != nil {
		panic(err)
	}
	receipt, err := txResponse.GetReceipt(user1.C)
	if err != nil {
		panic(err)
	}
	log.Println(receipt.Status)
}
