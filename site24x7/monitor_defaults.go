package site24x7

import (
	"errors"

	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/Bonial-International-GmbH/site24x7-go/api"
)

// DefaultLocationProfile fetches the first location profile returned by the
// client. If no location profiles are configured, DefaultLocationProfile will
// return an error.
func DefaultLocationProfile(client site24x7.Client) (*api.LocationProfile, error) {
	profiles, err := client.LocationProfiles().List()
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, errors.New("no location profiles configured")
	}

	return profiles[0], nil
}
