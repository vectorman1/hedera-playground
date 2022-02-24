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
	transferTx := hedera.NewTransferTransaction().
		AddHbarTransfer(user1.AccountId, hedera.HbarFromTinybar(-10)).
		AddHbarTransfer(user2.AccountId, hedera.HbarFromTinybar(10))

	scheduleTx, err := hedera.NewScheduleCreateTransaction().
		SetScheduledTransaction(transferTx)
	if err != nil {
		log.Fatalln(err)
	}

	submitScheduleTx, err := scheduleTx.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Get the receipt of the transaction
	receipt, err := submitScheduleTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Get the schedule ID
	scheduleId := *receipt.ScheduleID
	log.Printf("The schedule ID %v\n", scheduleId)

	defer func() {
		deleteScheduleTx, err := hedera.NewScheduleDeleteTransaction().
			SetScheduleID(scheduleId).
			FreezeWith(user1.C)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = deleteScheduleTx.Execute(user1.C)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	//Get the scheduled transaction ID
	scheduleTxId := receipt.ScheduledTransactionID
	log.Printf("The scheduled transaction ID is %v\n", scheduleTxId)

	// Submit the signatures for the scheduled transaction
	signature1, err := hedera.NewScheduleSignTransaction().
		SetScheduleID(scheduleId).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Verify the transaction was successful and submit a schedule info request
	submitTx, err := signature1.Sign(user1.PrivateKey).Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	receipt1, err := submitTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	status := receipt1.Status
	log.Printf("The transaction status is %v\n", status)

	query1, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduleId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(query1)

	//Submit the second signature
	signature2, err := hedera.NewScheduleSignTransaction().
		SetScheduleID(scheduleId).
		FreezeWith(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Verify the transaction was successful and submit a schedule info request
	submitTx2, err := signature2.Sign(user2.PrivateKey).Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	receipt2, err := submitTx2.GetReceipt(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	status2 := receipt2.Status
	log.Printf("The transaction status is %v\n", status2)

	//Get the schedule info
	query2, err := hedera.NewScheduleInfoQuery().
		SetScheduleID(scheduleId).
		Execute(user2.C)

	log.Println(query2)

	scheduledTxRecord, err := scheduleTxId.GetRecord(user1.C)
	log.Printf("The scheduled transaction record is %v\n", scheduledTxRecord)
}
