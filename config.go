package pact

import (
	"github.com/SEEK-Jobs/pact-go/util"
	"log"
	"os"
)

var (
	DefaultLogger         = log.New(os.Stderr, "\t", 0)
	DefaultVerifierConfig = &VerfierConfig{Logger: DefaultLogger}
	DefaultBuilderConfig  = &BuilderConfig{Logger: DefaultLogger}
	DefaultPactUriConfig  = &PactUriConfig{}
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

type PactUriConfig struct {
	Username string
	Password string
}
