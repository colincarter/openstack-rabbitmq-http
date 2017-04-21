package app

import "github.com/streadway/amqp"

// Application represents environment variables
type Application struct {
	RabbitConn *amqp.Connection
}

// App Global application state
var App Application
