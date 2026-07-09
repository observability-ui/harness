package component

import "maps"

type RequiredFlag struct {
	Name  string
	Usage string
}

type Component struct {
	Name          string
	DependsOn     []string
	Dir           string
	Ports         []int
	Config        map[string]string
	RequiredFlags []RequiredFlag
}

var registry = make(map[string]*Component)

func Register(c *Component) {
	if _, exists := registry[c.Name]; exists {
		panic("duplicate component registration: " + c.Name)
	}
	registry[c.Name] = c
}

func Get(name string) (*Component, bool) {
	c, ok := registry[name]
	return c, ok
}

func All() map[string]*Component {
	cp := make(map[string]*Component, len(registry))
	maps.Copy(cp, registry)
	return cp
}

func ResetRegistry() {
	registry = make(map[string]*Component)
}
