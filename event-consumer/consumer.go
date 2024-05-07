package eventconsumer

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	handler "github.com/thegoodparticle/challan-data-svc/rest-handler"
)

var (
	messageCountStart = 0
)

func SetupConsumer(brokers []string, topicName string, h *handler.Handler) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()

	consumer, err := master.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				messageCountStart++
				log.Println("Received messages", string(msg.Key), string(msg.Value))
				h.PostEvent(msg.Value)
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Println("Processed", messageCountStart, "messages")
}
