package resolver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	catalog["ccp"] = &ccpResolver{
		endpoint: "https://did.baidu.com/v1/did/resolve",
	}
}

// "ccp" DID Method
// https://did.baidu.com/did-spec/
type ccpResolver struct {
	endpoint string
}

// Resolve a specific DID instance as defined in the "ccp" method specification
func (sr *ccpResolver) Resolve(value string) ([]byte, error) {
	id, err := verify(value, "ccp")
	if err != nil {
		return nil, err
	}

	// Submit request
	res, err := http.Get(fmt.Sprintf("%s/%s", sr.endpoint, id))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	// Return response
	return ioutil.ReadAll(res.Body)
}
