package advanced

import "github.com/DaiYuANg/arcgo/dix"

func newProvider(
	label string,
	output dix.ServiceRef,
	register func(*dix.Container),
	deps ...dix.ServiceRef,
) dix.ProviderFunc {
	return dix.NewProviderFunc(register, dix.ProviderMetadata{
		Label:        label,
		Output:       output,
		Dependencies: dix.ServiceRefs(deps...),
	})
}

func newSetup(
	label string,
	run func(*dix.Container) error,
	dependencies []dix.ServiceRef,
	provides []dix.ServiceRef,
	overrides []dix.ServiceRef,
) dix.SetupFunc {
	return dix.NewSetupFunc(func(c *dix.Container, _ dix.Lifecycle) error {
		return run(c)
	}, dix.SetupMetadata{
		Label:         label,
		Dependencies:  dix.ServiceRefs(dependencies...),
		Provides:      dix.ServiceRefs(provides...),
		Overrides:     dix.ServiceRefs(overrides...),
		GraphMutation: false,
		Raw:           false,
	})
}
