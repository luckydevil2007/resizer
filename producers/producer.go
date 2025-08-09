package producers

import (
	"time"

	"github.com/IBM/sarama"
)

type EventProducer struct {
	producer sarama.SyncProducer
}

type TransformEvent struct {
	Image_id       int       `json:"image_id"`
	Timestamp      time.Time `json:"timestamp"`
	Transform_type string    `json:"transform_type"`
	User_id        int       `json:user_id`
}

type LoginEvent struct {
	User_id   int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewEventProducer(brokers []string /*, topic string*/) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &EventProducer{producer: producer}, nil
}

/*
func (ep *EventProducer) SendTransformEvent(t entities.ImageTransform) error {
	var transformType int
	transformType = 0
	if t.Resize > 0 {
		transformType++
	}
	if t.Rotate > 0 {
		transformType += 2
	}

	event := TransformEvent{
		Image_id:       t.ID,
		Timestamp:      time.Now(),
		User_id:        0,
		Transform_type: strconv.Itoa(transformType),
	}

	// random key
	key := "key-" + strconv.Itoa(rand.Intn(10))

	payload, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: "transforms_topic",
		Value: sarama.ByteEncoder(payload),
	}
	partition, offset, err := ep.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send event: %v", err)
		return err
	}
	log.Printf("Sent event ID %d to partition %d at offset %d: %s", event.Image_id, partition, offset, string(payload))
	return nil
}

func (ep *EventProducer) SendLoginEvent(u *entities.User) error {
	event := LoginEvent{
		User_id:   u.ID,
		Timestamp: time.Now(),
	}

	// random key
	key := "key-" + strconv.Itoa(rand.Intn(10))

	payload, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: "transforms_topic",
		Value: sarama.ByteEncoder(payload),
	}
	partition, offset, err := ep.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send event: %v", err)
		return err
	}
	log.Printf("Sent event ID %d to partition %d at offset %d: %s", event.User_id, partition, offset, string(payload))
	return nil
}

func (ep *EventProducer) Close() error {
	return ep.producer.Close()
}
*/
