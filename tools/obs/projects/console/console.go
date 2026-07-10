package console

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"obs/internal/runcontext"
	"obs/internal/task"
)

var Console = &task.Task{
	Name:  "console",
	Ports: []int{9000},
	Labels: map[string]string{
		"image": "quay.io/openshift/origin-console:latest",
	},
	Strategy: &ContainerStrategy{
		Image:    "quay.io/openshift/origin-console:latest",
		Platform: "linux/amd64",
	},
}

func init() {
	task.Register(Console)
}

type ContainerStrategy struct {
	Image    string
	Platform string
}

func (s *ContainerStrategy) Requires() []string { return []string{"oc"} }

func (s *ContainerStrategy) Execute(ctx context.Context, t *task.Task, rc *runcontext.RunContext) (*task.Step, error) {
	if err := extractClusterConfig(ctx, rc); err != nil {
		return nil, err
	}

	image := os.Getenv("CONSOLE_IMAGE")
	if image == "" {
		image = s.Image
	}
	platform := os.Getenv("CONSOLE_IMAGE_PLATFORM")
	if platform == "" {
		platform = s.Platform
	}

	rt := detectRuntime()
	if rt == "" {
		return nil, fmt.Errorf("no container runtime found — install podman or docker")
	}
	hostRef := hostReference(rt)
	plugins := discoverPlugins(rc)
	proxies := discoverProxies(rc)

	env := buildEnv(rc, plugins, proxies, hostRef)

	args := []string{"run", "--pull", "always", "--platform", platform, "--rm"}
	args = append(args, networkArgs(rt)...)

	for _, k := range slices.Sorted(maps.Keys(env)) {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, env[k]))
	}

	args = append(args, image)

	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleLongRunning,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    "console",
			Command: rt,
			Args:    args,
			Ports:   []int{9000},
		}},
	}, nil
}

func extractClusterConfig(ctx context.Context, rc *runcontext.RunContext) error {
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if user := runOC(checkCtx, "whoami"); user == "" {
		return fmt.Errorf("cluster unreachable or session expired — run 'oc login' first")
	}

	apiServer := runOC(checkCtx, "whoami", "--show-server")
	if apiServer == "" {
		return fmt.Errorf("could not determine cluster API server")
	}

	bearerToken := runOC(checkCtx, "whoami", "--show-token")
	if bearerToken == "" {
		fmt.Fprintln(os.Stderr, "warning: could not retrieve bearer token — console may not authenticate to the cluster")
	}
	thanosURL := runOC(checkCtx, "-n", "openshift-config-managed", "get", "configmap",
		"monitoring-shared-config", "-o", "jsonpath={.data.thanosPublicURL}")
	alertmanagerURL := runOC(checkCtx, "-n", "openshift-config-managed", "get", "configmap",
		"monitoring-shared-config", "-o", "jsonpath={.data.alertmanagerPublicURL}")

	rc.Set("console", "api-server", apiServer)
	rc.Set("console", "bearer-token", bearerToken)
	rc.Set("console", "thanos-url", thanosURL)
	rc.Set("console", "alertmanager-url", alertmanagerURL)

	return nil
}

func buildEnv(rc *runcontext.RunContext, plugins []pluginInfo, proxies []proxyInfo, hostRef string) map[string]string {
	env := map[string]string{
		"BRIDGE_USER_AUTH":                             "disabled",
		"BRIDGE_K8S_MODE":                              "off-cluster",
		"BRIDGE_K8S_AUTH":                               "bearer-token",
		"BRIDGE_K8S_MODE_OFF_CLUSTER_SKIP_VERIFY_TLS":  "true",
		"BRIDGE_USER_SETTINGS_LOCATION":                 "localstorage",
	}

	if v := rc.Get("console", "api-server"); v != "" {
		env["BRIDGE_K8S_MODE_OFF_CLUSTER_ENDPOINT"] = v
	}
	if v := rc.Get("console", "bearer-token"); v != "" {
		env["BRIDGE_K8S_AUTH_BEARER_TOKEN"] = v
	}
	if v := rc.Get("console", "thanos-url"); v != "" {
		env["BRIDGE_K8S_MODE_OFF_CLUSTER_THANOS"] = v
	}
	if v := rc.Get("console", "alertmanager-url"); v != "" {
		env["BRIDGE_K8S_MODE_OFF_CLUSTER_ALERTMANAGER"] = v
	}

	if len(plugins) > 0 {
		var bridgePlugins []string
		var i18nNamespaces []string
		for _, p := range plugins {
			bridgePlugins = append(bridgePlugins, fmt.Sprintf("%s=http://%s:%s", p.name, hostRef, p.port))
			i18nNamespaces = append(i18nNamespaces, fmt.Sprintf("plugin__%s", p.name))
		}
		env["BRIDGE_PLUGINS"] = strings.Join(bridgePlugins, ",")
		env["BRIDGE_I18N_NAMESPACES"] = strings.Join(i18nNamespaces, ",")
	}

	if len(proxies) > 0 {
		var proxyServices []map[string]any
		for _, p := range proxies {
			proxyServices = append(proxyServices, map[string]any{
				"consoleAPIPath": p.path,
				"endpoint":      fmt.Sprintf("http://%s:%s", hostRef, p.port),
				"authorize":     true,
			})
		}
		proxyJSON, err := json.Marshal(map[string]any{"services": proxyServices})
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to marshal proxy config: %v\n", err)
		} else {
			env["BRIDGE_PLUGIN_PROXY"] = string(proxyJSON)
		}
	}

	return env
}

func networkArgs(rt string) []string {
	if rt == "podman" && runtime.GOOS == "linux" {
		return []string{"--network=host"}
	}
	return []string{"-p", "9000:9000"}
}

func hostReference(rt string) string {
	if rt == "podman" {
		if runtime.GOOS == "linux" {
			return "localhost"
		}
		return "host.containers.internal"
	}
	return "host.docker.internal"
}

func detectRuntime() string {
	if _, err := exec.LookPath("podman"); err == nil {
		return "podman"
	}
	if _, err := exec.LookPath("docker"); err == nil {
		return "docker"
	}
	return ""
}

type pluginInfo struct {
	name string
	port string
}

type proxyInfo struct {
	path string
	port string
}

func discoverPlugins(rc *runcontext.RunContext) []pluginInfo {
	all := task.All()
	var plugins []pluginInfo
	for _, t := range all {
		if !rc.HasTask(t.Name) {
			continue
		}
		pluginName := t.Labels["console-plugin"]
		if pluginName == "" {
			continue
		}
		port := rc.Get(t.Name, "port")
		if port == "" && len(t.Ports) > 0 {
			port = strconv.Itoa(t.Ports[0])
		}
		if port != "" {
			plugins = append(plugins, pluginInfo{name: pluginName, port: port})
		}
	}
	return plugins
}

func discoverProxies(rc *runcontext.RunContext) []proxyInfo {
	all := task.All()
	var proxies []proxyInfo
	for _, t := range all {
		if !rc.HasTask(t.Name) {
			continue
		}
		path := t.Labels["console-proxy-path"]
		if path == "" {
			continue
		}
		port := t.Labels["console-proxy-port"]
		if port == "" {
			port = "8080"
		}
		proxies = append(proxies, proxyInfo{path: path, port: port})
	}
	return proxies
}

func runOC(ctx context.Context, args ...string) string {
	out, err := exec.CommandContext(ctx, "oc", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
