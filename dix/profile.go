package dix

import (
	"os"
)

// ProfileManager provides utilities for working with application profiles.
type ProfileManager struct{}

// GetProfileFromEnv retrieves the current profile from an environment variable.
// If the environment variable is not set or contains an invalid value,
// the default profile is returned.
//
// Example:
//
//	profile := ProfileFromEnv("APP_PROFILE", ProfileProd)
func ProfileFromEnv(envVar string, defaultProfile Profile) Profile {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultProfile
	}

	profile := Profile(value)
	switch profile {
	case ProfileDefault, ProfileDev, ProfileTest, ProfileProd:
		return profile
	default:
		return defaultProfile
	}
}

// IsProfile checks if the current profile matches the given profile.
func (pm ProfileManager) IsProfile(current, target Profile) bool {
	return current == target
}

// IsDev checks if the current profile is the development profile.
func (pm ProfileManager) IsDev(profile Profile) bool {
	return profile == ProfileDev
}

// IsTest checks if the current profile is the test profile.
func (pm ProfileManager) IsTest(profile Profile) bool {
	return profile == ProfileTest
}

// IsProd checks if the current profile is the production profile.
func (pm ProfileManager) IsProd(profile Profile) bool {
	return profile == ProfileProd
}

// Profile is a helper for creating profile-aware modules.
//
// Example:
//
//	var DevOnlyModule = dix.NewModule("dev-tools",
//	    dix.WithProviders(ProvideDebugHandler),
//	    dix.WithProfiles(dix.ProfileDev),
//	)
//
//	var ProdOnlyModule = dix.NewModule("monitoring",
//	    dix.WithProviders(ProvideMetrics),
//	    dix.WithExcludeProfiles(dix.ProfileDev, dix.ProfileTest),
//	)

// ProfileFilter provides methods for filtering modules by profile.
type ProfileFilter struct {
	profile Profile
}

// NewProfileFilter creates a new profile filter for the given profile.
func NewProfileFilter(profile Profile) *ProfileFilter {
	return &ProfileFilter{profile: profile}
}

// IsActive checks if a module should be active for the current profile.
func (pf *ProfileFilter) IsActive(mod Module) bool {
	return isActiveForProfile(mod, pf.profile)
}

// FilterModules returns only the modules that are active for the current profile.
func (pf *ProfileFilter) FilterModules(modules []Module) []Module {
	return flattenModules(modules, pf.profile)
}

// Global profile manager instance
var Profiles = ProfileManager{}
