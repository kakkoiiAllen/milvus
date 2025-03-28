package kafka

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/milvus-io/milvus/internal/common"
	"github.com/milvus-io/milvus/internal/log"
	"github.com/milvus-io/milvus/internal/mq/msgstream/mqwrapper"
	"github.com/milvus-io/milvus/internal/util/paramtable"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var Params paramtable.BaseTable

func TestMain(m *testing.M) {
	Params.Init()
	mockCluster, err := kafka.NewMockCluster(1)
	defer mockCluster.Close()
	if err != nil {
		fmt.Printf("Failed to create MockCluster: %s\n", err)
		os.Exit(1)
	}

	broker := mockCluster.BootstrapServers()
	Params.Save("kafka.brokerList", broker)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func getKafkaBrokerList() string {
	brokerList := Params.Get("kafka.brokerList")
	log.Info("get kafka broker list.", zap.String("address", brokerList))
	return brokerList
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, common.Endian, tmp)
	return bytesBuffer.Bytes()
}
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, common.Endian, &tmp)
	return int(tmp)
}

// Consume1 will consume random messages and record the last MessageID it received
func Consume1(ctx context.Context, t *testing.T, kc *kafkaClient, topic string, subName string, c chan mqwrapper.MessageID, total *int) {
	consumer, err := kc.Subscribe(mqwrapper.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            subName,
		BufSize:                     1024,
		SubscriptionInitialPosition: mqwrapper.SubscriptionPositionEarliest,
	})
	assert.Nil(t, err)
	assert.NotNil(t, consumer)
	defer consumer.Close()

	// get random number between 1 ~ 5
	rand.Seed(time.Now().UnixNano())
	cnt := 1 + rand.Int()%5

	log.Info("Consume1 start")
	var msg mqwrapper.Message
	for i := 0; i < cnt; i++ {
		select {
		case <-ctx.Done():
			log.Info("Consume1 channel closed")
			return
		case msg = <-consumer.Chan():
			if msg == nil {
				return
			}

			log.Info("Consume1 RECV", zap.Any("v", BytesToInt(msg.Payload())))
			consumer.Ack(msg)
			(*total)++
		}
	}

	c <- msg.ID()
	log.Info("Consume1 randomly RECV", zap.Any("number", cnt))
	log.Info("Consume1 done")
}

// Consume2 will consume messages from specified MessageID
func Consume2(ctx context.Context, t *testing.T, kc *kafkaClient, topic string, subName string, msgID mqwrapper.MessageID, total *int) {
	consumer, err := kc.Subscribe(mqwrapper.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            subName,
		BufSize:                     1024,
		SubscriptionInitialPosition: mqwrapper.SubscriptionPositionEarliest,
	})
	assert.Nil(t, err)
	assert.NotNil(t, consumer)
	defer consumer.Close()

	err = consumer.Seek(msgID, true)
	assert.Nil(t, err)

	mm := <-consumer.Chan()
	consumer.Ack(mm)
	log.Info("skip the last received message", zap.Any("skip msg", mm.ID()))

	log.Info("Consume2 start")
	for {
		select {
		case <-ctx.Done():
			log.Info("Consume2 channel closed")
			return
		case msg, ok := <-consumer.Chan():
			if msg == nil || !ok {
				return
			}

			log.Info("Consume2 RECV", zap.Any("v", BytesToInt(msg.Payload())))
			consumer.Ack(msg)
			(*total)++
		}
	}
}

func Consume3(ctx context.Context, t *testing.T, kc *kafkaClient, topic string, subName string, total *int) {
	consumer, err := kc.Subscribe(mqwrapper.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            subName,
		BufSize:                     1024,
		SubscriptionInitialPosition: mqwrapper.SubscriptionPositionEarliest,
	})
	assert.Nil(t, err)
	assert.NotNil(t, consumer)
	defer consumer.Close()

	log.Info("Consume3 start")
	for {
		select {
		case <-ctx.Done():
			log.Info("Consume3 channel closed")
			return
		case msg, ok := <-consumer.Chan():
			if msg == nil || !ok {
				return
			}

			consumer.Ack(msg)
			(*total)++
			log.Info("Consume3 RECV", zap.Any("v", BytesToInt(msg.Payload())), zap.Int("total", *total))
		}
	}
}

func TestKafkaClient_ConsumeWithAck(t *testing.T) {
	kc := createKafkaClient(t)
	defer kc.Close()
	assert.NotNil(t, kc)

	rand.Seed(time.Now().UnixNano())
	topic := fmt.Sprintf("test-topic-%d", rand.Int())
	subName := fmt.Sprintf("test-subname-%d", rand.Int())
	arr := []int{111, 222, 333, 444, 555, 666, 777}
	c := make(chan mqwrapper.MessageID, 1)

	ctx, cancel := context.WithCancel(context.Background())

	var total1 int
	var total2 int
	var total3 int

	producer := createProducer(t, kc, topic)
	defer producer.Close()
	produceData(ctx, t, producer, arr)
	time.Sleep(100 * time.Millisecond)

	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	Consume1(ctx1, t, kc, topic, subName, c, &total1)

	lastMsgID := <-c
	log.Info("lastMsgID", zap.Any("lastMsgID", lastMsgID.(*kafkaID).messageID))

	ctx2, cancel2 := context.WithTimeout(ctx, 3*time.Second)
	Consume2(ctx2, t, kc, topic, subName, lastMsgID, &total2)
	cancel2()

	time.Sleep(5 * time.Second)
	ctx3, cancel3 := context.WithTimeout(ctx, 3*time.Second)
	Consume3(ctx3, t, kc, topic, subName, &total3)
	cancel3()

	cancel()
	assert.Equal(t, len(arr), total1+total2)

	assert.Equal(t, len(arr), total3)
}

func TestKafkaClient_SeekPosition(t *testing.T) {
	kc := createKafkaClient(t)
	defer kc.Close()

	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	topic := fmt.Sprintf("test-topic-%d", rand.Int())
	subName := fmt.Sprintf("test-subname-%d", rand.Int())

	producer := createProducer(t, kc, topic)
	defer producer.Close()

	data := []int{1, 2, 3}
	ids := produceData(ctx, t, producer, data)

	consumer := createConsumer(t, kc, topic, subName, mqwrapper.SubscriptionPositionLatest)
	defer consumer.Close()

	err := consumer.Seek(ids[2], true)
	assert.Nil(t, err)

	select {
	case msg := <-consumer.Chan():
		consumer.Ack(msg)
		assert.Equal(t, 3, BytesToInt(msg.Payload()))
	case <-time.After(10 * time.Second):
		assert.FailNow(t, "should not wait")
	}
}

func TestKafkaClient_ConsumeFromLatest(t *testing.T) {
	kc := createKafkaClient(t)
	defer kc.Close()

	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	topic := fmt.Sprintf("test-topic-%d", rand.Int())
	subName := fmt.Sprintf("test-subname-%d", rand.Int())

	producer := createProducer(t, kc, topic)
	defer producer.Close()

	data := []int{1, 2}
	produceData(ctx, t, producer, data)

	consumer := createConsumer(t, kc, topic, subName, mqwrapper.SubscriptionPositionLatest)
	defer consumer.Close()

	go func() {
		time.Sleep(time.Second * 2)
		data := []int{3}
		produceData(ctx, t, producer, data)
	}()

	select {
	case msg := <-consumer.Chan():
		consumer.Ack(msg)
		assert.Equal(t, 3, BytesToInt(msg.Payload()))
	case <-time.After(5 * time.Second):
		assert.FailNow(t, "should not wait")
	}
}

func TestKafkaClient_EarliestMessageID(t *testing.T) {
	kafkaAddress := getKafkaBrokerList()
	kc := NewKafkaClientInstance(kafkaAddress)
	defer kc.Close()

	mid := kc.EarliestMessageID()
	assert.NotNil(t, mid)
}

func TestKafkaClient_MsgSerializAndDeserialize(t *testing.T) {
	kafkaAddress := getKafkaBrokerList()
	kc := NewKafkaClientInstance(kafkaAddress)
	defer kc.Close()

	mid := kc.EarliestMessageID()
	msgID, err := kc.BytesToMsgID(mid.Serialize())
	assert.NoError(t, err)
	assert.True(t, msgID.AtEarliestPosition())

	msgID, err = kc.StringToMsgID("1")
	assert.NoError(t, err)
	assert.NotNil(t, msgID)

	msgID, err = kc.StringToMsgID("1.0")
	assert.Error(t, err)
	assert.Nil(t, msgID)
}

func TestKafkaClient_NewKafkaClientInstanceWithConfig(t *testing.T) {
	config1 := &paramtable.KafkaConfig{Address: "addr", SaslPassword: "password"}
	assert.Panics(t, func() { NewKafkaClientInstanceWithConfig(config1) })

	config2 := &paramtable.KafkaConfig{Address: "addr", SaslUsername: "username"}
	assert.Panics(t, func() { NewKafkaClientInstanceWithConfig(config2) })

	config3 := &paramtable.KafkaConfig{Address: "addr", SaslUsername: "username", SaslPassword: "password"}
	client := NewKafkaClientInstanceWithConfig(config3)
	assert.NotNil(t, client)
	assert.NotNil(t, client.basicConfig)
}

func createKafkaClient(t *testing.T) *kafkaClient {
	kafkaAddress := getKafkaBrokerList()
	kc := NewKafkaClientInstance(kafkaAddress)
	assert.NotNil(t, kc)
	return kc
}

func createConsumer(t *testing.T,
	kc *kafkaClient,
	topic string,
	groupID string,
	initPosition mqwrapper.SubscriptionInitialPosition) mqwrapper.Consumer {
	consumer, err := kc.Subscribe(mqwrapper.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            groupID,
		BufSize:                     1024,
		SubscriptionInitialPosition: initPosition,
	})
	assert.Nil(t, err)
	return consumer
}

func createProducer(t *testing.T, kc *kafkaClient, topic string) mqwrapper.Producer {
	producer, err := kc.CreateProducer(mqwrapper.ProducerOptions{Topic: topic})
	assert.Nil(t, err)
	assert.NotNil(t, producer)
	return producer
}

func produceData(ctx context.Context, t *testing.T, producer mqwrapper.Producer, arr []int) []mqwrapper.MessageID {
	var msgIDs []mqwrapper.MessageID
	for _, v := range arr {
		msg := &mqwrapper.ProducerMessage{
			Payload:    IntToBytes(v),
			Properties: map[string]string{},
		}
		msgID, err := producer.Send(ctx, msg)
		msgIDs = append(msgIDs, msgID)
		assert.Nil(t, err)
	}

	producer.(*kafkaProducer).p.Flush(500)
	return msgIDs
}
