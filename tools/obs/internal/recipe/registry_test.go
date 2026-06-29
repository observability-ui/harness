package recipe_test

import (
	"testing"

	"github.com/spf13/pflag"
	"obs/internal/recipe"
)

type stubRecipe struct {
	name    string
	aliases []string
	flags   *pflag.FlagSet
}

func (s *stubRecipe) Name() string              { return s.name }
func (s *stubRecipe) Aliases() []string          { return s.aliases }
func (s *stubRecipe) Description() string        { return s.name + " recipe" }
func (s *stubRecipe) Flags() *pflag.FlagSet {
	if s.flags != nil {
		return s.flags
	}
	return pflag.NewFlagSet(s.name, pflag.ContinueOnError)
}
func (s *stubRecipe) Requirements(_ *pflag.FlagSet) []recipe.Requirement { return nil }
func (s *stubRecipe) Steps(_ *recipe.Config) ([]*recipe.Step, error) { return nil, nil }

func newStubWithFlags(name string, aliases []string, flags func(fs *pflag.FlagSet)) *stubRecipe {
	s := &stubRecipe{name: name, aliases: aliases}
	if flags != nil {
		s.flags = pflag.NewFlagSet(name, pflag.ContinueOnError)
		flags(s.flags)
	}
	return s
}

func TestRegistryRegisterAndLookup(t *testing.T) {
	reg := recipe.NewRegistry()
	r := &stubRecipe{name: "monitoring-plugin", aliases: []string{"mp"}}

	if err := reg.Register("start", r); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	got, ok := reg.Lookup("start", "monitoring-plugin")
	if !ok || got.Name() != "monitoring-plugin" {
		t.Fatalf("Lookup by name failed: ok=%v got=%v", ok, got)
	}

	got, ok = reg.Lookup("start", "mp")
	if !ok || got.Name() != "monitoring-plugin" {
		t.Fatalf("Lookup by alias failed: ok=%v got=%v", ok, got)
	}

	_, ok = reg.Lookup("deploy", "mp")
	if ok {
		t.Fatal("Lookup wrong command should return false")
	}
}

func TestRegistryDuplicateRegister(t *testing.T) {
	reg := recipe.NewRegistry()
	r := &stubRecipe{name: "mp", aliases: nil}
	if err := reg.Register("start", r); err != nil {
		t.Fatal(err)
	}
	if err := reg.Register("start", r); err == nil {
		t.Fatal("duplicate Register should return error")
	}
}

func TestRegistryList(t *testing.T) {
	reg := recipe.NewRegistry()
	reg.Register("start", &stubRecipe{name: "a"})
	reg.Register("deploy", &stubRecipe{name: "b"})
	reg.Register("start", &stubRecipe{name: "c"})

	startRecipes := reg.List("start")
	if len(startRecipes) != 2 {
		t.Fatalf("expected 2 start recipes, got %d", len(startRecipes))
	}

	all := reg.ListAll()
	if len(all) != 3 {
		t.Fatalf("expected 3 total recipes, got %d", len(all))
	}
}
