package pact

import (
	"errors"
)

type Builder interface {
	ServiceConsumer(consumer string) Builder
	HasPactWith(serviceProvider string) Builder
	GetMockProviderService() (p ProviderService, serverUrl string)
	Build() error
}

type Config struct {
	PactPath string
	LogPath  string
}

type ConsumerPactBuilder struct {
	consumer        string
	serviceProvider string
	providerService *mockProviderService
}

var (
	errInvalidConsumer = errors.New("ConsumerName has not been set, please supply a consumer name using the ServiceConsumer method.")
	errInvalidProvider = errors.New("ProviderName has not been set, please supply a provider name using the HasPactWith method.")
)

func NewConsumerPactBuilder(pactConfig *Config) Builder {
	return &ConsumerPactBuilder{providerService: newMockProviderService(pactConfig)}
}

func (c *ConsumerPactBuilder) ServiceConsumer(consumer string) Builder {
	c.consumer = consumer
	return c
}

func (c *ConsumerPactBuilder) HasPactWith(serviceProvider string) Builder {
	c.serviceProvider = serviceProvider
	return c
}

func (c *ConsumerPactBuilder) GetMockProviderService() (ProviderService, string) {
	url := c.providerService.start()
	return c.providerService, url
}

func (c *ConsumerPactBuilder) Build() error {
	if c.consumer == "" {
		return errInvalidConsumer
	} else if c.serviceProvider == "" {
		return errInvalidProvider
	}
	if err := c.providerService.persistPact(c.consumer, c.serviceProvider); err != nil {
		return err
	} else {
		c.providerService.stop()
		return nil
	}
}
