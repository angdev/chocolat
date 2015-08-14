package api

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
)

type Params struct {
	req      *rest.Request
	requires []string
}

func NewParams(req *rest.Request) *Params {
	return &Params{req: req}
}

func (this *Params) Require(keys ...string) *Params {
	this.requires = keys
	return this
}

func (this *Params) Parse(out interface{}) error {
	values := make(map[string]interface{})

	if err := this.req.DecodeJsonPayload(&values); err != nil {
		return err
	}

	// Query parameters
	for _, require := range this.requires {
		v := this.req.FormValue(require)
		if _, ok := values[require]; !ok && v != "" {
			values[require] = v
		}
	}

	// Path parameters
	for _, require := range this.requires {
		v := this.req.PathParam(require)
		if _, ok := values[require]; !ok && v != "" {
			values[require] = v
		}
	}

	for _, require := range this.requires {
		if _, ok := values[require]; !ok {
			return ParamsMissingError
		}
	}

	if result, err := json.Marshal(values); err != nil {
		return err
	} else {
		return json.Unmarshal(result, out)
	}
}
