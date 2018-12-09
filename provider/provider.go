package provider

import (
	"log"
	"strings"
)

var providers = make([]Provider, 0)

type Provider interface {
	Match(val string) bool
	Decode(va string) (string, error)
}

func Register(provider Provider) {
	providers = append(providers, provider)
}

func Populate(env []string) []string {
	newEnv := make([]string, 0)
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		for _, p := range providers {
			if p.Match(pair[1]) {
				if decoded, err := p.Decode(pair[1]); err == nil {
					pair[1] = decoded
				} else {
					log.Printf("Error decoding variable '%v': %v", pair[0], err)
				}
				break
			}
		}

		newEnv = append(newEnv, pair[0] + "=" + pair[1])
	}

	return newEnv
}
