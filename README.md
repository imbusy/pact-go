# pact-go
A Go Lang implementation of the Ruby consumer driven contract library, Pact.
Pact is based off the specification found at https://github.com/bethesque/pact_specification.

Currently pact-go is compatible with v1.1 of [pact specification](https://github.com/pact-foundation/pact-specification/tree/version-1.1). At this stage only the consumer dsl is available to generate and verify pacts. We will be adding the provider side of functionality soon.

Read more about Pact and the problems it solves at [https://github.com/realestate-com-au/pact](https://github.com/realestate-com-au/pact)

Please feel free to contribute, we do accept pull requests.
### Installing
```shell
go get github.com/SEEK-Jobs/pact-go
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
  pact "github.com/SEEK-Jobs/pact-go"
)

func buildPact() {
  return pact.
    NewConsumerPactBuilder(&pact.Config{PactPath: "./pacts"}).
		ServiceConsumer("my consumer").
		HasPactWith("my provider")
}
```
##### 3. Add your tests
Add your test method to register and verify your interactions with provider
```go

import (
  pact "github.com/SEEK-Jobs/pact-go"
  "github.com/SEEK-Jobs/pact-go/provider"
  "net/http"
)

func buildPact() {
  return pact.
    NewConsumerPactBuilder(&pact.Config{PactPath: "./pacts"}).
		ServiceConsumer("my consumer").
		HasPactWith("my provider")
}

func TestPactWithProvider(t *testing.T) {
  builder := buildPact()
  ms, msUrl := builder.GetMockProviderService()

  request := provider.NewJSONRequest("GET", "/23", nil, nil)
	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJSONResponse(200, header)
	response.SetBody(`{"name": "John"}`)

  //Register interaction for this test scope
  if err := ps.Given("there is a resource with id 23").
		UponReceiving("get request for resource with id 23").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
		t.FailNow()
	}

  //test
  client := &ProviderClient{baseUrl: msUrl}
  if res, err := client.GetResource(23); err != nil {
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
