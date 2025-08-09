package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/luckydevil2007/audionotes/database"
)

type TransformEvent struct {
	Image_id       int       `json:"image_id"`
	Timestamp      time.Time `json:"timestamp"`
	Transform_type string    `json:"transform_type"`
	User_id        int       `json:user_id`
}

func main() {

	consumer := Consumer{
		ready: make(chan bool),
	}
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split("localhost:9092", ","), "ThisConsumer", config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	defer cancel()

	//	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split("transforms_topic", ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up

	wg.Wait()

}

type Consumer struct {
	ready        chan bool
	clickhouseDb *sql.DB
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	consumer.clickhouseDb = database.OpenClickhouse()
	ctreateTableQuery := "CREATE TABLE  IF NOT EXISTS events (image_id UInt32, timestamp DateTime,transform_type UInt32, user_id UInt32) ENGINE = MergeTree() ORDER BY (timestamp, image_id)"
	consumer.clickhouseDb.Exec(ctreateTableQuery)

	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			session.MarkMessage(message, "")
			var msg TransformEvent
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				return fmt.Errorf("error unmarshaling message: %w", err)
			}
			consumer.clickhouseDb.Exec(`INSERT INTO events (image_id, timestamp, transform_type, user_id) VALUES ($1, $2, $3, $4 )`, msg.Image_id, msg.Timestamp, msg.Transform_type, msg.User_id)

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}

}
