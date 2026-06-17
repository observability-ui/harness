package recipe

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	byCmd   map[string]map[string]Recipe // command -> name -> recipe
	aliases map[string]map[string]string // command -> alias -> name
}

func NewRegistry() *Registry {
	return &Registry{
		byCmd:   make(map[string]map[string]Recipe),
		aliases: make(map[string]map[string]string),
	}
}

func (r *Registry) Register(rec Recipe) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cmd := rec.Command()
	if r.byCmd[cmd] == nil {
		r.byCmd[cmd] = make(map[string]Recipe)
		r.aliases[cmd] = make(map[string]string)
	}
	if _, exists := r.byCmd[cmd][rec.Name()]; exists {
		return fmt.Errorf("recipe %q already registered for command %q", rec.Name(), cmd)
	}
	r.byCmd[cmd][rec.Name()] = rec
	for _, alias := range rec.Aliases() {
		r.aliases[cmd][alias] = rec.Name()
	}
	return nil
}

func (r *Registry) Lookup(command, nameOrAlias string) (Recipe, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	recipes := r.byCmd[command]
	if recipes == nil {
		return nil, false
	}
	if rec, ok := recipes[nameOrAlias]; ok {
		return rec, true
	}
	if name, ok := r.aliases[command][nameOrAlias]; ok {
		return recipes[name], true
	}
	return nil, false
}

func (r *Registry) List(command string) []Recipe {
	r.mu.RLock()
	defer r.mu.RUnlock()

	recipes := r.byCmd[command]
	result := make([]Recipe, 0, len(recipes))
	for _, rec := range recipes {
		result = append(result, rec)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name() < result[j].Name() })
	return result
}

func (r *Registry) ListAll() []Recipe {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Recipe
	for _, recipes := range r.byCmd {
		for _, rec := range recipes {
			result = append(result, rec)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name() < result[j].Name() })
	return result
}

// DefaultRegistry is the global registry used by CLI commands.
var DefaultRegistry = NewRegistry()
