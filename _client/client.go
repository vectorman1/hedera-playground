package client

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"log"
)

type hederaUser struct {
	C          *hedera.Client
	AccountId  hedera.AccountID
	PrivateKey hedera.PrivateKey
	PublicKey  hedera.PublicKey
}

type hederaSuite struct {
	Users []*hederaUser
}

type accountCredentials struct {
	AccountId  string
	PrivateKey string
}

func Setup1TestClient() *hederaSuite {
	hs, err := newHederaSuite([]*accountCredentials{
		{
			"0.0.30779785",
			"302e020100300506032b6570042204207baac8952b698f13f3874dc9dbff19e5b9914f51f8c78616e934701923a3a278",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	return hs
}

func Setup2TestClients() *hederaSuite {
	hs, err := newHederaSuite([]*accountCredentials{
		{
			"0.0.30779785",
			"302e020100300506032b6570042204207baac8952b698f13f3874dc9dbff19e5b9914f51f8c78616e934701923a3a278",
		},
		{
			"0.0.30783406",
			"302e020100300506032b6570042204207832e442a414fcd03306c063f5f5ea6b915b67b2092ac376c014224cccd21ef2",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	return hs
}

func newHederaSuite(accounts []*accountCredentials) (*hederaSuite, error) {
	hs := &hederaSuite{}

	for _, acc := range accounts {
		err := hs.setupClient(acc.AccountId, acc.PrivateKey)
		if err != nil {
			return nil, err
		}
	}

	return hs, nil
}

func (h *hederaSuite) setupClient(accountId string, privateKey string) error {
	parsedAccountId, err := hedera.AccountIDFromString(accountId)
	if err != nil {
		return err
	}
	parsedPrivateKey, err := hedera.PrivateKeyFromString(privateKey)
	if err != nil {
		return err
	}

	client := hedera.ClientForTestnet()
	client.SetOperator(parsedAccountId, parsedPrivateKey)

	publicKey := client.GetOperatorPublicKey()

	h.Users = append(h.Users, &hederaUser{
		C:          client,
		AccountId:  parsedAccountId,
		PrivateKey: parsedPrivateKey,
		PublicKey:  publicKey,
	})

	return nil
}
