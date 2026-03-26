package dix

// SetupFunc describes a typed setup registration.
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

// Setup registers a typed setup callback.
func Setup(fn func(*Container, Lifecycle) error) SetupFunc {
	return SetupWithMetadata(fn, SetupMetadata{
		Label: "Setup",
	})
}

// SetupWithMetadata registers a typed setup callback with metadata.
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

// RawSetup registers an untyped setup callback.
func RawSetup(fn func(*Container, Lifecycle) error) SetupFunc {
	return RawSetupWithMetadata(fn, SetupMetadata{
		Label: "RawSetup",
	})
}

// RawSetupWithMetadata registers an untyped setup callback with metadata.
func RawSetupWithMetadata(fn func(*Container, Lifecycle) error, meta SetupMetadata) SetupFunc {
	meta.Raw = true
	return NewSetupFunc(fn, meta)
}
