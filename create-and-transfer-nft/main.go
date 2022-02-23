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

	// Create the NFT
	nftCreate, err := hedera.NewTokenCreateTransaction().
		SetTokenType(hedera.TokenTypeNonFungibleUnique).
		SetTokenName("TOPKITE NA BULGARIQ").
		SetTokenSymbol("BALLS").
		SetDecimals(0).
		SetInitialSupply(0).
		SetTreasuryAccountID(user1.AccountId).
		SetAdminKey(user1.PublicKey).
		SetWipeKey(user1.PublicKey).
		SetSupplyKey(user1.PublicKey).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	nftCreateTxSign := nftCreate.Sign(user1.PrivateKey)
	nftCreateSubmit, err := nftCreateTxSign.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	nftCreateRx, err := nftCreateSubmit.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	tokenId := *nftCreateRx.TokenID
	fmt.Printf("Token created: %s\n", tokenId)

	// Mint an NFT
	// Best thing in the world
	CID := "QmVpK83eWSNxGaecqb7MGPats1BSMmpbWmp7zrY8pXg2Yk"
	mintTx, err := hedera.NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadata([]byte(CID)).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	mintTxSign := mintTx.Sign(user1.PrivateKey)
	mintTxSubmit, err := mintTxSign.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	mintRx, err := mintTxSubmit.GetReceipt(user1.C)
	fmt.Println("New NFT", tokenId, "with serial", mintRx.SerialNumbers)

	// Associate User2's account with the token
	associateUser2Tx, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(user2.AccountId).
		SetTokenIDs(tokenId).
		FreezeWith(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	signTx := associateUser2Tx.Sign(user2.PrivateKey)
	associateUser2TxSubmit, err := signTx.Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	associateUser2Rx, err := associateUser2TxSubmit.GetReceipt(user2.C)
	fmt.Println("NFT Association with User2's account:", associateUser2Rx.Status)

	// Transfer the NFT to User2's account
	nftId := hedera.NftID{
		TokenID:      tokenId,
		SerialNumber: 1,
	}
	tokenTransferTx, err := hedera.NewTransferTransaction().
		AddNftTransfer(nftId, user1.AccountId, user2.AccountId).
		FreezeWith(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	signTransferTx := tokenTransferTx.Sign(user1.PrivateKey)
	tokenTransferSubmit, err := signTransferTx.Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	tokenTransferRx, err := tokenTransferSubmit.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("NFT Transfer from User1 (Treasurer) to User2:", tokenTransferRx.Status)

	// Check the balance of the treasury account after the transfer
	balanceCheckUser1, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user1.AccountId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Treasury balance:", balanceCheckUser1.Tokens, "NFTs of ID", tokenId)

	// Check the balance of Alice's account after the transfer
	balanceCheckUser2, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("User2's balance:", balanceCheckUser2.Tokens, "NFTs of ID", tokenId)
}
