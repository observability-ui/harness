package recipes

import (
	"obs/internal/recipe"
	"obs/recipes/lp"
	"obs/recipes/mp"
)

func init() {
	recipe.DefaultRegistry.Register("start", &mp.StartMonitoringPlugin{})
	recipe.DefaultRegistry.Register("deploy", &mp.DeployMonitoringPlugin{})
	recipe.DefaultRegistry.Register("start", &lp.StartLoggingPlugin{})
}
