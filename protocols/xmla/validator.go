// protocols/xmla/validator.go
package xmla

import "errors"

func (env *Envelope) Validate() error {
	if env.Body.Discover == nil {
		return errors.New("discover request is required")
	}
	if env.Body.Discover.RequestType == "" {
		return errors.New("request type is empty")
	}
	return nil
}
