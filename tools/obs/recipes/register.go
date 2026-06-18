package recipes

import (
	"obs/internal/recipe"
	"obs/recipes/mp"
)

func init() {
	recipe.DefaultRegistry.Register("start", &mp.StartMonitoringPlugin{})
}
