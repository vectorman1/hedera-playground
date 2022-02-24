package main

import (
	"encoding/json"
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"io/ioutil"
	"log"
)

func main() {
	hs := client.Setup2TestClients()
	user1 := hs.Users[0]
	user2 := hs.Users[1]

	hts, err := ioutil.ReadFile("deploy-smartcontract-hts/contracts/HTS.json")
	if err != nil {
		log.Fatalln(err)
	}

	var contract contract
	err = json.Unmarshal(hts, &contract)
	if err != nil {
		log.Fatalln(err)
	}

	htsByteCode := contract.Data.Bytecode.Object
	fileCreateTx, err := hedera.NewFileCreateTransaction().
		SetContents([]byte(htsByteCode)).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	fileCreateRx, err := fileCreateTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	byteCodeFileId := *fileCreateRx.FileID
	log.Println("Created file with ID:", byteCodeFileId)

	contractCreateTx, err := hedera.NewContractCreateTransaction().
		SetGas(100000).
		SetBytecodeFileID(byteCodeFileId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	contractCreateRx, err := contractCreateTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	contractId := *contractCreateRx.ContractID
	log.Println("Created contract with Id:", contractId)

	tokenCreateTx, err := hedera.NewTokenCreateTransaction().
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetTokenName("Stoyan Kolev Coin").
		SetTokenSymbol("STK").
		SetDecimals(2).
		SetInitialSupply(10000).
		SetTreasuryAccountID(user1.AccountId).
		SetSupplyType(hedera.TokenSupplyTypeInfinite).
		SetSupplyKey(user1.PrivateKey).
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
	log.Printf("Created fungible token %s with ID %s", "Stoyan Kolev Coin", tokenId)

	// Associate the token with the receiver's account
	tokenIdSol := tokenId.ToSolidityAddress()
	accountIdSol := user2.AccountId.ToSolidityAddress()
	htsAssocParams, err := hedera.NewContractFunctionParameters().
		AddAddress(accountIdSol)
	if err != nil {
		log.Fatalln(err)
	}
	htsAssocParams, err = htsAssocParams.
		AddAddress(tokenIdSol)
	if err != nil {
		log.Fatalln(err)
	}

	htsAssociateTx, err := hedera.NewContractExecuteTransaction().
		SetContractID(contractId).
		SetFunction("tokenAssociate", htsAssocParams).
		SetGas(100000).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	htsAssociateRx, err := htsAssociateTx.GetReceipt(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	associateStatus := htsAssociateRx.Status
	log.Println("Associate token status", associateStatus)

	// Check records of association
	//Get the child transaction record
	childRecord, err := hedera.NewTransactionRecordQuery().
		//Set boolean equal to true
		SetIncludeChildren(true).
		//Parent transaction ID
		SetTransactionID(htsAssociateTx.TransactionID).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Log the child record
	log.Printf("The associate child transaction record %v\n", childRecord.Children)

	//The balance of the account
	accountBalance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("The account token balance %v\n", accountBalance.Tokens)

	// Set up the parameters we will use to call the smartcontract transfer function
	htsTokenTransferParams, err := hedera.NewContractFunctionParameters().
		AddAddress(tokenIdSol)
	if err != nil {
		log.Fatalln(err)
	}
	htsTokenTransferParams, err = htsTokenTransferParams.
		AddAddress(user1.AccountId.ToSolidityAddress())
	if err != nil {
		log.Fatalln(err)
	}
	htsTokenTransferParams, err = htsTokenTransferParams.
		AddAddress(user2.AccountId.ToSolidityAddress())
	if err != nil {
		log.Fatalln(err)
	}
	htsTokenTransferParams.AddInt64(250)

	//Transfer the token
	transferTx := hedera.NewContractExecuteTransaction().
		//The contract ID
		SetContractID(contractId).
		//The max gas
		SetGas(2000000).
		//The contract function to call and parameters
		SetFunction("tokenTransfer", htsTokenTransferParams)

	//Sign with treasury key to authorize the transfer from the treasury account
	signTx, err := transferTx.Sign(user1.PrivateKey).Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Get the receipt
	transferTxReceipt, err := signTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Get transaction status
	transferTxStatus := transferTxReceipt.Status

	log.Printf("The transfer transaction status %v\n", transferTxStatus)

	//Verify the transfer by checking the balance
	transferAccountBalance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(user2.AccountId).
		Execute(user2.C)
	if err != nil {
		log.Fatalln(err)
	}

	//Log the account token balance
	log.Printf("The account token balance %v\n", transferAccountBalance.Tokens)
}
