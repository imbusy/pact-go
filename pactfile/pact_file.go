package pactfile

import (
	"encoding/json"
	"fmt"
	"github.com/SEEK-Jobs/pact-go/consumer"
)

type Participant struct {
	Name string `json:"name"`
}

type PactFile struct {
	Consumer     *Participant            `json:"consumer"`
	Provider     *Participant            `json:"provider"`
	Interactions []*consumer.Interaction `json:"interactions"`
	Metadata     interface{}             `json:"metaData"`
}

func NewPactFile(consumer string, provider string, interactions []*consumer.Interaction) *PactFile {
	return &PactFile{
		Consumer:     &Participant{Name: consumer},
		Provider:     &Participant{Name: provider},
		Interactions: interactions,
		Metadata: struct {
			PactSpecification string `json:"pactSpecification"`
		}{PactSpecification: "1.1.0"},
	}
}

func (p *PactFile) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PactFile) FileName() string {
	return fmt.Sprintf("%s-%s.json", p.Consumer.Name, p.Provider.Name)
}
