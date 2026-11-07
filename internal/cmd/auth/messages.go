package auth

import "fmt"

const (
	noProfilesConfiguredMessage = "No profiles configured. Run 'rdflow auth login' to add one."
	noActiveProfileMessage      = "No active profile selected. Run 'rdflow auth switch' to select one."
)

func profileNotFoundError(name string) error {
	return fmt.Errorf("profile %q not found. Run 'rdflow auth list' to see available profiles", name)
}
