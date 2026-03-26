package dix

type SetupFunc struct {
	run  func(*Container, Lifecycle) error
	meta SetupMetadata
}

func (s SetupFunc) apply(c *Container, lc Lifecycle) error {
	if s.run == nil {
		return nil
	}
	return s.run(c, lc)
}

func Setup(fn func(*Container, Lifecycle) error) SetupFunc {
	return SetupWithMetadata(fn, SetupMetadata{
		Label: "Setup",
	})
}

func SetupWithMetadata(fn func(*Container, Lifecycle) error, meta SetupMetadata) SetupFunc {
	return NewSetupFunc(fn, SetupMetadata{
		Label:         meta.Label,
		Dependencies:  meta.Dependencies,
		Provides:      meta.Provides,
		Overrides:     meta.Overrides,
		GraphMutation: meta.GraphMutation,
		Raw:           meta.Raw,
	})
}

func RawSetup(fn func(*Container, Lifecycle) error) SetupFunc {
	return RawSetupWithMetadata(fn, SetupMetadata{
		Label: "RawSetup",
	})
}

func RawSetupWithMetadata(fn func(*Container, Lifecycle) error, meta SetupMetadata) SetupFunc {
	meta.Raw = true
	return NewSetupFunc(fn, meta)
}
