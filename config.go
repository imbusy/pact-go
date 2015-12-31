package pact

import (
	"github.com/SEEK-Jobs/pact-go/util"
	"log"
	"os"
)

var (
	DefaultVerifierConfig = &VerfierConfig{Logger: log.New(os.Stderr, "\t", 0)}
	DefaultBuilderConfig  = &BuilderConfig{Logger: log.New(os.Stderr, "\t", 0)}
)

//BuilderConfig configuration needed to build pacts
type BuilderConfig struct {
	PactPath string
	Logger   util.Logger
}

//VerifierConfig configuration needed to verify pacts
type VerfierConfig struct {
	Logger util.Logger
}
