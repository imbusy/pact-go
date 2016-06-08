package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/imbusy/pact-go/consumer"
	version "github.com/hashicorp/go-version"
)

const pactSpecificationVersion = "1.1.0"

var (
	errEmptyProvider    = errors.New("Pactfile is invalid, provider name should not be empty.")
	errEmptyConsumer    = errors.New("Pactfile is invalid, consumer name should not be empty.")
	errIncompatiblePact = fmt.Errorf("Incompatible pact specification! We only support version %s.", pactSpecificationVersion)
)

type Participant struct {
	Name string `json:"name"`
}

type metadata struct {
	PactSpecificationVersion string `json:"pactSpecificationVersion"`
}

type PactFile struct {
	Consumer     *Participant            `json:"consumer"`
	Provider     *Participant            `json:"provider"`
	Interactions []*consumer.Interaction `json:"interactions"`
	Metadata     *metadata               `json:"metaData"`
}

func NewPactFile(consumer string, provider string, interactions []*consumer.Interaction) *PactFile {
	return &PactFile{
		Consumer:     &Participant{Name: consumer},
		Provider:     &Participant{Name: provider},
		Interactions: interactions,
		Metadata:     &metadata{PactSpecificationVersion: pactSpecificationVersion},
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

func (p *PactFile) Validate() error {
	if p.Provider == nil || p.Provider.Name == "" {
		return errEmptyProvider
	}

	if p.Consumer == nil || p.Consumer.Name == "" {
		return errEmptyConsumer
	}

	psv, err := version.NewVersion(pactSpecificationVersion)
	if err != nil {
		return err
	}

	fpsv, err := version.NewVersion(p.Metadata.PactSpecificationVersion)
	if err != nil {
		return err
	}

	//should be backwards compatible with version 1.0.0
	if fpsv.GreaterThan(psv) {
		return errIncompatiblePact
	}

	return nil
}
