package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"log"
)

func main() {
	hs := client.Setup2TestClients()
	user1 := hs.Users[0]
	user2 := hs.Users[1]

	// Creating the token
	tokenName := "Stoyan Kolev Coin"
	supplyKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		log.Fatalln(err)
	}
	tokenCreateTx, err := hedera.NewTokenCreateTransaction().
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetTokenName(tokenName).
		SetTokenSymbol("STK").
		SetDecimals(2).
		SetInitialSupply(10000).
		SetTreasuryAccountID(user1.AccountId).
		SetSupplyType(hedera.TokenSupplyTypeInfinite).
		SetSupplyKey(supplyKey).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	tokenCreateSign := tokenCreateTx.Sign(user1.PrivateKey)
	tokenCreateExec, err := tokenCreateSign.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	tokenCreateRx, err := tokenCreateExec.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	tokenId := *tokenCreateRx.TokenID
	fmt.Println(fmt.Sprintf("Created fungible token %s with ID %s", tokenName, tokenId))

	// Associate User2 with the Token
	associateUser2Tx, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(user2.AccountId).
		SetTokenIDs(tokenId).
		FreezeWith(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	associateUser2Exec, err := associateUser2Tx.Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	associateUser2Rx, err := associateUser2Exec.GetReceipt(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Associated", user2.AccountId, "with token", tokenId, associateUser2Rx.Status)

	//Check the balance before the transfer for the treasury account
	balanceCheckTreasury, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user1.AccountId).
		Execute(user1.C)
	fmt.Println("Treasury balance:", balanceCheckTreasury.Tokens, "units of token ID", tokenId)
	if err != nil {
		log.Fatalln(err)
	}

	//Check the balance before the transfer for Alice's account
	balanceCheckAlice, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	fmt.Println("User2's balance:", balanceCheckAlice.Tokens, "units of token ID", tokenId)
	if err != nil {
		log.Fatalln(err)
	}

	tokenTransferTx, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenId, user1.AccountId, -2500).
		AddTokenTransfer(tokenId, user2.AccountId, 2500).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	signTransferTx := tokenTransferTx.Sign(user1.PrivateKey)
	tokenTransferSubmit, err := signTransferTx.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	tokenTransferRx, err := tokenTransferSubmit.GetReceipt(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Token transfer from Treasury to User2:", tokenTransferRx.Status)

	//Check the balance before the transfer for the treasury account
	balanceCheckTreasury, err = hedera.NewAccountBalanceQuery().
		SetAccountID(user1.AccountId).
		Execute(user1.C)
	fmt.Println("Treasury balance:", balanceCheckTreasury.Tokens, "units of token ID", tokenId)
	if err != nil {
		log.Fatalln(err)
	}

	//Check the balance before the transfer for Alice's account
	balanceCheckAlice, err = hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	fmt.Println("User2's balance:", balanceCheckAlice.Tokens, "units of token ID", tokenId)
	if err != nil {
		log.Fatalln(err)
	}
}
