package recipes

import (
	"obsui/internal/recipe"
	"obsui/recipes/startmp"
)

func init() {
	recipe.DefaultRegistry.Register(&startmp.StartMonitoringPlugin{})
}
