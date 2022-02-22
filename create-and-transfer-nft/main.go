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

	nftId := hedera.NftID{
		TokenID:      tokenId,
		SerialNumber: mintRx.SerialNumbers[0],
	}
	nftInfo, err := hedera.NewTokenNftInfoQuery().
		SetNftID(nftId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(nftInfo)
}
