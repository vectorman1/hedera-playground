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

	helloHedera, err := ioutil.ReadFile("deploy-smartcontract/contracts/HelloHedera.json")
	if err != nil {
		log.Fatalln(err)
	}

	var contract contract
	err = json.Unmarshal(helloHedera, &contract)
	if err != nil {
		log.Fatalln(err)
	}

	bytecode := contract.Data.Bytecode.Object
	fileTx, err := hedera.NewFileCreateTransaction().
		SetContents([]byte(bytecode)).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	fileRx, err := fileTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	byteCodeFileId := *fileRx.FileID
	log.Println("Created file of the contract's bytecode with ID: ", byteCodeFileId)

	contractTx, err := hedera.NewContractCreateTransaction().
		SetBytecodeFileID(byteCodeFileId).
		SetGas(100000).
		SetConstructorParameters(hedera.NewContractFunctionParameters().
			AddString("Hello from Hedera!")).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	contractRx, err := contractTx.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	contractId := *contractRx.ContractID
	log.Println("Created contract on hedera with ID: ", contractId)

	callFnQuery, err := hedera.NewContractCallQuery().
		SetContractID(contractId).
		SetGas(10000000).
		SetFunction("getMessage", nil).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	getMessage := callFnQuery.GetString(0)
	log.Println("Contract response:", getMessage)
}
