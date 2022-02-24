package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"log"
)

func main() {
	hs := client.Setup2TestClients()
	user1 := hs.Users[0]
	user1Bal, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user1.AccountId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	user2 := hs.Users[1]
	user2Bal, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(user1Bal.Hbars.AsTinybar(), user2Bal.Hbars.AsTinybar())
}
