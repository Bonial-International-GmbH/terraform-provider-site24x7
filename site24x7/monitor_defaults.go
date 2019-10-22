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

// DefaultNotificationProfile fetches the first notification profile returned by the
// client. If no notification profiles are configured, DefaultNotificationProfile will
// return an error.
func DefaultNotificationProfile(client site24x7.Client) (*api.NotificationProfile, error) {
	profiles, err := client.NotificationProfiles().List()
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, errors.New("no notification profiles configured")
	}

	return profiles[0], nil
}

// DefaultThresholdProfile fetches the first threshold profile returned by the
// client. If no threshold profiles are configured, DefaultThresholdProfile will
// return an error.
func DefaultThresholdProfile(client site24x7.Client) (*api.ThresholdProfile, error) {
	profiles, err := client.ThresholdProfiles().List()
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, errors.New("no threshold profiles configured")
	}

	return profiles[0], nil
}

// DefaultUserGroup fetches the first user group returned by the
// client. If no user groups are configured, DefaultUserGroup will
// return an error.
func DefaultUserGroup(client site24x7.Client) (*api.UserGroup, error) {
	userGroups, err := client.UserGroups().List()
	if err != nil {
		return nil, err
	}

	if len(userGroups) == 0 {
		return nil, errors.New("no user groups configured")
	}

	return userGroups[0], nil
}
