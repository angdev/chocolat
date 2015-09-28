package api

import (
	"github.com/jmoiron/jsonq"
	"github.com/mitchellh/mapstructure"
	"net/url"
	"strings"
)

type Keen struct {
	Addons []KeenAddon
}

type KeenAddon struct {
	Name   string
	Input  map[string]interface{}
	Output string
}

func (this *KeenAddon) ResolveInput(doc map[string]interface{}) error {
	var err error
	jq := jsonq.NewQuery(doc)

	for k, v := range this.Input {
		if this.Input[k], err = jq.Interface(strings.Split(v.(string), ".")...); err != nil {
			return err
		}
	}

	return nil
}

type KeenAddonApplier interface {
	Apply(doc map[string]interface{}) error
}

type KeenUrlParser KeenAddon

func (this *KeenUrlParser) Apply(doc map[string]interface{}) error {
	var input struct {
		Url string
	}
	if err := mapstructure.Decode(this.Input, &input); err != nil {
		return err
	}

	if parsed, err := url.Parse(input.Url); err != nil {
		return err
	} else {
		doc[this.Output] = &struct {
			Protocol    string
			Domain      string
			Path        string
			Anchor      string
			QueryString url.Values
		}{
			parsed.Scheme,
			parsed.Host,
			parsed.Path,
			parsed.Fragment,
			parsed.Query(),
		}

		return nil
	}
}

func resolveKeenAddon(addon *KeenAddon) KeenAddonApplier {
	switch addon.Name {
	case "keen:url_parser":
		return (*KeenUrlParser)(addon)
	default:
		return nil
	}
}

func resolveKeenObject(doc map[string]interface{}) error {
	keenObject, ok := doc["keen"]
	if !ok {
		return nil
	}

	var keen Keen
	if err := mapstructure.Decode(keenObject, &keen); err != nil {
		return err
	}

	for _, addon := range keen.Addons {
		if err := addon.ResolveInput(doc); err != nil {
			return err
		}

		if applier := resolveKeenAddon(&addon); applier != nil {
			if err := applier.Apply(doc); err != nil {
				return err
			}
		}
	}

	delete(doc, "keen")

	return nil
}
