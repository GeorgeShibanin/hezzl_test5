package nats

import (
	"encoding/json"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsQueue struct {
	conn *nats.Conn
}

func initConnection(conn *nats.Conn) *NatsQueue {
	return &NatsQueue{conn: conn}
}

func Init(url string) (*NatsQueue, error) {
	sc, err := nats.Connect(url)
	log.Println("Connected to " + url)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}

	log.Println("Connected Nats")
	return initConnection(sc), nil
}

func (n *NatsQueue) PushMessage(message storage.Item) (string, error) {
	m, err := json.Marshal(message)
	if err != nil {
		//Обработать ошибку
		log.Fatal(err)
		return "FAIL to Marahall", err
	}
	err = n.conn.Publish("logs", m)
	if err != nil {
		return "FAIL to Publish", err
	}
	return "WellDoneYouSend", nil
}

func (n *NatsQueue) GetMessage() (storage.Item, error) {
	subj, i := "logs", 0
	var data storage.Item
	mcb := func(msg *nats.Msg) {
		i++
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Println("Trouble with message in channel(cant unmarshal)")
			return
		}
	}
	_, err := n.conn.Subscribe(subj, mcb)
	if err != nil {
		log.Fatal(err)
		return storage.Item{}, err
	}
	return data, nil
}
