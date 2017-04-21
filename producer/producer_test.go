package producer_test

import (
	"testing"

	"bitbucket.org/colincarter/openstack-rabbitmq-http/producer"
)

func TestNewProducer(t *testing.T) {
	producer, err := producer.NewProducer()
}
