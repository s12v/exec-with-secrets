package provider

import (
	"errors"
	"strings"
	"testing"
)

type DummyProvider struct {
	result string
	broken bool
}

func (*DummyProvider) Match(val string) bool {
	return strings.HasPrefix(val, "{dummy}")
}

func (p *DummyProvider) Decode(val string) (string, error) {
	if p.broken {
		return val, errors.New("test error")
	}
	return p.result, nil
}

func TestRegister(t *testing.T) {
	registry = make([]Provider, 0)
	dummy := DummyProvider{"bar", false}
	Register(&dummy)
	if registry[0] != &dummy {
		t.Fatal("expected the provider to be added")
	}
}

func TestPopulate(t *testing.T) {
	env := []string{
		"PLAINTEXT=this is a test",
		"SECRET={dummy}ddd",
	}

	registry = []Provider{&DummyProvider{"bar", false}}
	env = Populate(env)

	if env[0] != "PLAINTEXT=this is a test" {
		t.Fatalf("unexpected env[0]: %v", env[0])
	}

	if env[1] != "SECRET=bar" {
		t.Fatalf("unexpected env[1]: %v", env[1])
	}
}

func TestPopulate_Error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic")
		}
	}()

	registry = make([]Provider, 0)
	env := []string{
		"PLAINTEXT=this is a test",
		"SECRET={dummy}ddd",
	}

	registry = []Provider{&DummyProvider{"bar", true}}
	env = Populate(env)
}
