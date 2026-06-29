package recipe

var providers = make(map[string]StepProvider)

func RegisterProvider(p StepProvider) {
	providers[p.Name()] = p
}

func GetProvider(name string) (StepProvider, bool) {
	p, ok := providers[name]
	return p, ok
}
