package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"log"
	"time"
)

func main() {
	hs := client.Setup2TestClients()

	// Create topic with first user
	user1 := hs.Users[0]
	txResponse, err := hedera.NewTopicCreateTransaction().
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	topicCreateReceipt, err := txResponse.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	topicId := *topicCreateReceipt.TopicID
	fmt.Println(topicId)

	// Subscribe to the topic with the second user
	user2 := hs.Users[1]
	_, err = hedera.NewTopicMessageQuery().
		SetTopicID(topicId).
		Subscribe(user2.C, func(msg hedera.TopicMessage) {
			fmt.Println(msg.ConsensusTimestamp.String(), "received topic message ", string(msg.Contents), "\r")
		})

	// Send message on the topic with the first user
	submitMessage, err := hedera.NewTopicMessageSubmitTransaction().
		SetMessage([]byte("Hello world!")).
		SetTopicID(topicId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	topicSubmitReceipt, err := submitMessage.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	sendStatus := topicSubmitReceipt.Status
	fmt.Println("Sending message", sendStatus)
	time.Sleep(30000)
}
