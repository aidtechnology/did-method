package resolver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	catalog["stack"] = &stackResolver{
		endpoint: "https://core.blockstack.org/v1/dids",
	}
}

// "stack" DID Method
// https://github.com/blockstack/blockstack-core/blob/master/docs/blockstack-did-spec.md
type stackResolver struct {
	endpoint string
}

// Resolve a specific DID instance as defined in the "stack" method specification
func (sr *stackResolver) Resolve(value string) ([]byte, error) {
	id, err := verify(value, "stack")
	if err != nil {
		return nil, err
	}

	// Submit request
	res, err := http.Get(fmt.Sprintf("%s/%s", sr.endpoint, id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Return response
	return ioutil.ReadAll(res.Body)
}
