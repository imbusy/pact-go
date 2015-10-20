package consumer

type PactBuilder interface {
	ServiceConsumer(consumer string) PactBuilder
	HasPactWith(serviceProvider string) PactBuilder
	GetMockProviderService() *MockProviderService
	Build() error
}

type PactConfig struct {
	PactPath string
	LogPath  string
}

type ConsumerPactBuilder struct {
	consumer            string
	serviceProvider     string
	config              *PactConfig
	mockProviderService *MockProviderService
}

func NewConsumerPactBuilder(pactConfig *PactConfig) *ConsumerPactBuilder {
	return &ConsumerPactBuilder{config: pactConfig}
}

func (c *ConsumerPactBuilder) ServiceConsumer(consumer string) PactBuilder {
	c.consumer = consumer
	return c
}

func (c *ConsumerPactBuilder) HasPactWith(serviceProvider string) PactBuilder {
	c.serviceProvider = serviceProvider
	return c
}

func (c *ConsumerPactBuilder) GetMockProviderService() *MockProviderService {
	c.mockProviderService = NewMockProviderService(c.config)
	return c.mockProviderService
}

func (c *ConsumerPactBuilder) Build() error {
	return nil
}
