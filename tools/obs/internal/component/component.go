package component

type Output struct {
	Name  string
	Value string
}

type Component struct {
	Name         string
	Description  string
	DependsOn    []string
	Capabilities []string
	Dir          string
	Outputs      []Output
	Config       map[string]string
}

var registry = make(map[string]*Component)

func Register(c *Component) {
	registry[c.Name] = c
}

func Get(name string) (*Component, bool) {
	c, ok := registry[name]
	return c, ok
}

func All() map[string]*Component {
	return registry
}
