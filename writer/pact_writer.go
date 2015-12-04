package writer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type PactWriter interface {
	Write() error
}

type PactFileWriter struct {
	pact     *PactFile
	pactPath string
}

func NewPactFileWriter(pact *PactFile, path string) PactWriter {
	return &PactFileWriter{
		pact:     pact,
		pactPath: path,
	}
}

func (p *PactFileWriter) Write() error {
	if _, err := os.Stat(p.pactPath); os.IsNotExist(err) {
		os.MkdirAll(p.pactPath, 0777)
	}

	filename := path.Join(p.pactPath, p.pact.FileName())

	data, err := p.pact.ToJson()

	if err != nil {
		return err
	}

	//indent
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "\t"); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, out.Bytes(), 0777); err != nil {
		return err
	}

	return nil
}
