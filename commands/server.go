package commands

import (
	"fmt"
	"log"
	"sync"
	"time"

	"bitbucket.org/colincarter/openstack-rabbitmq-http/handlers"
	"bitbucket.org/colincarter/openstack-rabbitmq-http/producer"

	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
	graceful "gopkg.in/tylerb/graceful.v1"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	rabbitConn      *amqp.Connection
	rabbitChanWg    sync.WaitGroup
	rabbitEventChan chan []byte
	rabbitMeterChan chan []byte
	quitChan        chan bool
	producers       []*producer.Producer
)

func init() {
	rabbitEventChan = make(chan []byte)
	rabbitMeterChan = make(chan []byte)
	quitChan = make(chan bool)
}

// Server runs the app
func Server(c *cli.Context) error {
	defer cleanup()

	var err error

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.String("rabbit-user"),
		c.String("rabbit-password"),
		c.String("rabbit-host"),
		c.Int("rabbit-port"))

	rabbitConn, err = amqp.Dial(amqpURI)
	if err != nil {
		return err
	}

	for i := 0; i < c.Int("concurrency"); i++ {
		p, err := producer.NewProducer(
			rabbitConn,
			rabbitEventChan,
			rabbitMeterChan,
			quitChan,
			&rabbitChanWg,
			c.String("rabbit-exchange"),
			c.String("failure-dir"))
		if err != nil {
			return err
		}
		producers = append(producers, p)
		go p.Run()
	}

	router := httprouter.New()
	router.POST("/events", handlers.HandleEvents(rabbitEventChan))
	router.POST("/meters", handlers.HandleMeters(rabbitMeterChan))
	router.GET("/ping", handlers.HandlePing)

	log.Printf("Listening on %s:%s", c.String("listen"), c.String("port"))

	graceful.Run(fmt.Sprintf("%s:%s", c.String("listen"), c.String("port")), 5*time.Second, router)

	return nil
}

func cleanup() {
	log.Println("Cleaning up")

	for _, p := range producers {
		p.Shutdown()
	}

	rabbitChanWg.Wait()

	if quitChan != nil {
		close(quitChan)
	}

	if rabbitEventChan != nil {
		close(rabbitEventChan)
	}

	if rabbitMeterChan != nil {
		close(rabbitMeterChan)
	}

	if rabbitConn != nil {
		rabbitConn.Close()
	}
}
