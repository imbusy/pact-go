package pact

type PactBuilder interface {
	ServiceConsumer(consumer string) PactBuilder
	HasPactWith(serviceProvider string) PactBuilder
	GetMockProviderService() ProviderService
	Build() error
}

type PactConfig struct {
	PactPath string
	LogPath  string
}

type ConsumerPactBuilder struct {
	consumer        string
	serviceProvider string
	config          *PactConfig
	providerService ProviderService
}

func NewConsumerPactBuilder(pactConfig *PactConfig) *ConsumerPactBuilder {
	return &ConsumerPactBuilder{config: pactConfig, providerService: NewMockProviderService(pactConfig)}
}

func (c *ConsumerPactBuilder) ServiceConsumer(consumer string) PactBuilder {
	c.consumer = consumer
	return c
}

func (c *ConsumerPactBuilder) HasPactWith(serviceProvider string) PactBuilder {
	c.serviceProvider = serviceProvider
	return c
}

func (c *ConsumerPactBuilder) GetMockProviderService() ProviderService {
	return c.providerService
}

func (c *ConsumerPactBuilder) Build() error {
	return nil
}
