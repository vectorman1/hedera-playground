package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	client "hedera-playground/_client"
	"log"
	"time"
)

func main() {
	hs := client.Setup1TestClient()

	// Create topic with first user
	user1 := hs.Users[0]
	txResponse, err := hedera.NewTopicCreateTransaction().
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	topicCreateRx, err := txResponse.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}

	topicId := *topicCreateRx.TopicID
	log.Println("Created topic with ID:", topicId, topicCreateRx.Status)

	// Subscribe to the topic
	received := make(chan struct{})
	_, err = hedera.NewTopicMessageQuery().
		SetTopicID(topicId).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(user1.C, func(msg hedera.TopicMessage) {
			log.Println(msg.ConsensusTimestamp.String(), "received topic message", string(msg.Contents), "\r")
			received <- struct{}{}
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Subscribed to topic")

	// Send message on the topic with the first user
	submitMessage, err := hedera.NewTopicMessageSubmitTransaction().
		SetMessage([]byte("Hello world!")).
		SetTopicID(topicId).
		Execute(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Sent a message to topic")

	submitMessageRx, err := submitMessage.GetReceipt(user1.C)
	if err != nil {
		log.Fatalln(err)
	}
	
	sendStatus := submitMessageRx.Status
	log.Println("Sending message status:", sendStatus.String())
	log.Println("Waiting for receive")
	<-received
}
