package pact

import (
	"errors"
)

//Builder Pact Builder
type Builder interface {
	//ServiceConsumer consumer which creates the pact
	ServiceConsumer(consumer string) Builder
	//HasPactWith sets the provider with which cosumer has pact
	HasPactWith(serviceProvider string) Builder
	//GetMockProviderService Returns the mock provider service and it's url
	GetMockProviderService() (p ProviderService, serverURL string)
	//Build Builds the pact and persists it
	Build() error
}

type consumerPactBuilder struct {
	consumer        string
	serviceProvider string
	providerService *mockProviderService
}

var (
	errInvalidConsumer = errors.New("ConsumerName has not been set, please supply a consumer name using the ServiceConsumer method.")
	errInvalidProvider = errors.New("ProviderName has not been set, please supply a provider name using the HasPactWith method.")
)

//NewConsumerPactBuilder Factory method to create new consumer pact builder
func NewConsumerPactBuilder(pactConfig *BuilderConfig) Builder {
	if pactConfig == nil {
		pactConfig = DefaultBuilderConfig
	}
	return &consumerPactBuilder{providerService: newMockProviderService(pactConfig)}
}

func (c *consumerPactBuilder) ServiceConsumer(consumer string) Builder {
	c.consumer = consumer
	return c
}

func (c *consumerPactBuilder) HasPactWith(serviceProvider string) Builder {
	c.serviceProvider = serviceProvider
	return c
}

func (c *consumerPactBuilder) GetMockProviderService() (ProviderService, string) {
	url := c.providerService.start()
	return c.providerService, url
}

func (c *consumerPactBuilder) Build() error {
	if c.consumer == "" {
		return errInvalidConsumer
	} else if c.serviceProvider == "" {
		return errInvalidProvider
	}
	if err := c.providerService.persistPact(c.consumer, c.serviceProvider); err != nil {
		return err
	}

	c.providerService.stop()
	return nil

}
