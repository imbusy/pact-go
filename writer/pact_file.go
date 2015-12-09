package writer

import (
	"encoding/json"
	"fmt"
	"strings"

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
	consumer := strings.Replace(strings.ToLower(p.Consumer.Name), " ", "_", -1)
	provider := strings.Replace(strings.ToLower(p.Provider.Name), " ", "_", -1)
	return fmt.Sprintf("%s-%s.json", consumer, provider)
}
