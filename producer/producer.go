package producer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// Producer represents channels for sending to Rabbit
type Producer struct {
	rabbitChannel          *amqp.Channel
	rabbitExchange         string
	wg                     *sync.WaitGroup
	eventsChan, metersChan chan []byte
	quitChan               chan bool
	failureDir             string
}

// NewProducer creates a new producer
func NewProducer(conn *amqp.Connection,
	eventsChan chan []byte,
	metersChan chan []byte,
	quitChan chan bool,
	wg *sync.WaitGroup,
	rabbitExchange string,
	failureDir string) (*Producer, error) {

	if failureDir != "" {
		if err := os.MkdirAll(failureDir, 0777); err != nil {
			log.Printf(
				"Unable to create failure directory %s - %s", failureDir, err)
			failureDir = ""
		}
	}

	c, err := connectAMQP(conn, rabbitExchange)
	if err != nil {
		return nil, err
	}

	return &Producer{
		rabbitChannel:  c,
		rabbitExchange: rabbitExchange,
		wg:             wg,
		eventsChan:     eventsChan,
		metersChan:     metersChan,
		quitChan:       quitChan,
		failureDir:     failureDir}, nil
}

// Run runs  a producer
func (p *Producer) Run() {
	p.wg.Add(1)
	defer p.wg.Done()

	sendMsg := func(routingKey string, data []byte) {
		timeStamp := time.Now()
		err := p.rabbitChannel.Publish(
			p.rabbitExchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp:    timeStamp,
				ContentType:  "text/plain",
				Body:         data,
			})
		if err != nil {
			log.Printf("Error publishing %s", string(data))
			p.writeFailure(
				fmt.Sprintf("%s/%s-%d.txt",
					p.failureDir,
					routingKey,
					timeStamp.UnixNano()),
				data)
		}
	}

	for {
		select {
		case event := <-p.eventsChan:
			sendMsg("raw_events", event)
		case meter := <-p.metersChan:
			sendMsg("raw_meters", meter)
		case <-p.quitChan:
			p.rabbitChannel.Close()
			return
		}
	}
}

// Shutdown producer
func (p *Producer) Shutdown() {
	p.quitChan <- true
}

func (p *Producer) writeFailure(filename string, data []byte) {
	if p.failureDir == "" {
		return
	}
	ioutil.WriteFile(filename, data, 0777)
}

func connectAMQP(conn *amqp.Connection, exchange string) (*amqp.Channel, error) {
	c, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
