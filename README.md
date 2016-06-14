# pact-go
A Go Lang implementation of the Ruby consumer driven contract library, Pact.
Pact is based off the specification found at https://github.com/bethesque/pact_specification.

Currently pact-go is compatible with v1.1 of [pact specification](https://github.com/pact-foundation/pact-specification/tree/version-1.1). The package is fairly stable but in beta phase, so potentially there could be breaking changes in future.

Read more about Pact and the problems it solves at [https://github.com/realestate-com-au/pact](https://github.com/realestate-com-au/pact)

Please feel free to contribute, we do accept pull requests.
### Installing
```shell
go get github.com/imbusy/pact-go
```

### Usage

#### Service Consumer
##### 1. Build your client
Which may look something like this
```go
package client

import (
  "fmt"
  "net/http"
  "encoding/json"
)
type ProviderAPIClient struct {
  baseURL string
}

type Resource struct {
  Name string
}

func (c *ProviderAPIClient) GetResource(id int)  (*Resource, error) {
  url := fmt.Sprintf("%s/%d", c.baseURL, id)
  req, _ := http.NewRequest("GET", url, nil)

  client := &http.Client{}
  resp, err:= client.Do(req)

  if err != nil { return nil, err}
  defer resp.Body.Close()

  var res Resource
  decoder := json.NewDecoder(resp.Body)
  if err := decoder.Decode(&res); err != nil {
    return nil, err
  }

  return &res, nil
}
```
##### 2. Configure pact
In your test file describe and configure your pact
```go
import (
  pact "github.com/imbusy/pact-go"
)

func buildPact() pact.Builder {
  return pact.
    NewConsumerPactBuilder(&pact.Config{PactPath: "./pacts"}).
		ServiceConsumer("my consumer").
		HasPactWith("my provider")
}
```
##### 3. Add your tests
Add your test method to register and verify your interactions with provider
```go
package client

import (
	pact "github.com/imbusy/pact-go"
	"github.com/imbusy/pact-go/provider"
	"net/http"
	"testing"
)

func buildPact() pact.Builder {
	return pact.
		NewConsumerPactBuilder(&pact.BuilderConfig{PactPath: "./pacts"}).
		ServiceConsumer("consumer client").
		HasPactWith("provider api")
}

func TestPactWithProvider(t *testing.T) {
	builder := buildPact()
	ms, msUrl := builder.GetMockProviderService()

	request := provider.NewJSONRequest("GET", "/23", "", nil)
	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJSONResponse(200, header)
	response.SetBody(`{"Name": "John"}`)

	//Register interaction for this test scope
	if err := ms.Given("there is a resource with id 23").
		UponReceiving("get request for resource with id 23").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
		t.FailNow()
	}

	//test
	client := &ProviderAPIClient{baseURL: msUrl}
	if _, err := client.GetResource(23); err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Verify registered interaction
	if err := ms.VerifyInteractions(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	//Clear interaction for this test scope, if you need to register and verify another interaction for another test scope
	ms.ClearInteractions()

	//Finally, build to produce the pact json file
	if err := builder.Build(); err != nil {
		t.Error(err)
	}
}
```

##### 4. Run your test
```shell
go test -v ./...
```

#### Service Provider
##### 1. Build the API
Build your api using any web framework like Goji, Gin, Gorilla or just net/http. Ensure your pipeline is composable to simplify testing and introduce mock behaviors.

##### 2. Honour Pacts
Create a new test file in your service provider api to verify pact with the consumer. Each test should contain one provider state. When you are done verifying all the
individual states, verify that all states were tested.

```go
package provider

import (
	"encoding/json"
	"fmt"
	pact "github.com/imbusy/pact-go"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type Resource struct {
	Name string
}

func TestProviderHonoursPactWithConsumer(t *testing.T) {
	//Stand up a test server with the mux which has all your
	//middlewares and handlers registered, for e.g.
	mux := http.NewServeMux()
	mux.HandleFunc("/23", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		fmt.Fprint(w, encoder.Encode(Resource{"John"}))
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	u, _ := url.Parse(server.URL)

	verifier := pact.NewPactFileVerifier(nil, nil, pact.DefaultVerifierConfig).
		HonoursPactWith("consumer client").
		ServiceProvider("provider api").
		//pact uri could be a local file
		PactUri("../pacts/consumer_client-provider_api.json", nil).
		//or could be a web url e.g. pact broker. You can also provide authorisation value in the config parameter
		//PactUri("http://pact-broker/pacts/provider/provider%20api/consumer/consumer%20client/version/latest", nil)

	if err := verifier.VerifyState("there is a resource with id 23", &http.Client{}, u); err != nil {
		t.Error(err)
	}

	if err := verifier.VerifyAllStatesTested(); err != nil {
		t.Error(err)
	}
}

```

#### 3. Run your test
Run your test to verify all the interactions.

```shell
go test -v ./...
```
