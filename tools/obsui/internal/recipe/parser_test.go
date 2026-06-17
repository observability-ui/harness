package recipe_test

import (
	"testing"

	"github.com/spf13/pflag"
	"obsui/internal/recipe"
)

func TestParseRecipeArgs_Single(t *testing.T) {
	reg := recipe.NewRegistry()
	reg.Register("start", newStubWithFlags("monitoring-plugin", []string{"mp"}, func(fs *pflag.FlagSet) {
		fs.String("version", "", "version to use")
	}))

	segments, err := recipe.ParseRecipeArgs(reg, "start", []string{"mp", "--version=4.18"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}
	if segments[0].Recipe.Name() != "monitoring-plugin" {
		t.Fatalf("expected monitoring-plugin, got %s", segments[0].Recipe.Name())
	}
	v, _ := segments[0].Flags.GetString("version")
	if v != "4.18" {
		t.Fatalf("expected version=4.18, got %q", v)
	}
}

func TestParseRecipeArgs_Multiple(t *testing.T) {
	reg := recipe.NewRegistry()
	reg.Register("start", newStubWithFlags("monitoring-plugin", []string{"mp"}, func(fs *pflag.FlagSet) {
		fs.String("version", "", "")
	}))
	reg.Register("start", newStubWithFlags("console", []string{"con"}, func(fs *pflag.FlagSet) {
		fs.String("version", "", "")
	}))

	segments, err := recipe.ParseRecipeArgs(reg, "start", []string{"mp", "--version=4.18", "con", "--version=4.19"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}
	v0, _ := segments[0].Flags.GetString("version")
	v1, _ := segments[1].Flags.GetString("version")
	if v0 != "4.18" || v1 != "4.19" {
		t.Fatalf("versions wrong: %q %q", v0, v1)
	}
}

func TestParseRecipeArgs_UnknownRecipe(t *testing.T) {
	reg := recipe.NewRegistry()
	_, err := recipe.ParseRecipeArgs(reg, "start", []string{"nonexistent"})
	if err == nil {
		t.Fatal("unknown recipe should return error")
	}
}

func TestParseRecipeArgs_NoArgs(t *testing.T) {
	reg := recipe.NewRegistry()
	_, err := recipe.ParseRecipeArgs(reg, "start", []string{})
	if err == nil {
		t.Fatal("no args should return error")
	}
}
