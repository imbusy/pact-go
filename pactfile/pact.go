package pactfile

import (
	"encoding/json"
	"fmt"
	"github.com/bennycao/pact-go/consumer"
)

type Participant struct {
	Name string `json:"name"`
}

type Pact struct {
	Consumer     *Participant            `json:"consumer"`
	Provider     *Participant            `json:"provider"`
	Interactions []*consumer.Interaction `json:"interactions"`
	Metadata     interface{}             `json:"metaData"`
}

func NewPact(consumer string, provider string, interactions []*consumer.Interaction) *Pact {
	return &Pact{
		Consumer:     &Participant{Name: consumer},
		Provider:     &Participant{Name: provider},
		Interactions: interactions,
		Metadata: struct {
			PactSpecification string `json:"pactSpecification"`
		}{PactSpecification: "1.1.0"},
	}
}

func (p *Pact) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Pact) FileName() string {
	return fmt.Sprintf("%s-%s.json", p.Consumer.Name, p.Provider.Name)
}
