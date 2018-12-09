package provider

import (
	"fmt"
	"strings"
)

type Provider interface {
	Match(val string) bool
	Decode(val string) (string, error)
}

var registry []Provider

func init() {
	registry = make([]Provider, 0)
}

func Populate(env []string) []string {
	newEnv := make([]string, 0)
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		for _, p := range registry {
			if p.Match(pair[1]) {
				if decoded, err := p.Decode(pair[1]); err == nil {
					pair[1] = decoded
				} else {
					panic(fmt.Sprintf("Error decoding variable '%v': %v\n", pair[0], err))
				}
				break
			}
		}

		newEnv = append(newEnv, pair[0] + "=" + pair[1])
	}

	return newEnv
}

func Register(provider Provider) {
	registry = append(registry, provider)
}
