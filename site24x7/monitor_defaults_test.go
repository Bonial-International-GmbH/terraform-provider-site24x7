package site24x7

import (
	"errors"
	"testing"

	"github.com/Bonial-International-GmbH/site24x7-go/api"
	"github.com/Bonial-International-GmbH/site24x7-go/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultLocationProfile(t *testing.T) {
	client := fake.NewClient()

	client.FakeLocationProfiles.On("List").Return(nil, errors.New("an error occurred")).Once()

	_, err := DefaultLocationProfile(client)

	require.Equal(t, errors.New("an error occurred"), err)

	client.FakeLocationProfiles.On("List").Return(nil, nil).Once()

	_, err = DefaultLocationProfile(client)

	require.Equal(t, errors.New("no location profiles configured"), err)

	client.FakeLocationProfiles.On("List").Return([]*api.LocationProfile{
		{ProfileID: "456"},
		{ProfileID: "123"},
	}, nil).Once()

	profile, err := DefaultLocationProfile(client)

	require.NoError(t, err)
	assert.Equal(t, &api.LocationProfile{ProfileID: "456"}, profile)
}

func TestDefaultNotificationProfile(t *testing.T) {
	client := fake.NewClient()

	client.FakeNotificationProfiles.On("List").Return(nil, errors.New("an error occurred")).Once()

	_, err := DefaultNotificationProfile(client)

	require.Equal(t, errors.New("an error occurred"), err)

	client.FakeNotificationProfiles.On("List").Return(nil, nil).Once()

	_, err = DefaultNotificationProfile(client)

	require.Equal(t, errors.New("no notification profiles configured"), err)

	client.FakeNotificationProfiles.On("List").Return([]*api.NotificationProfile{
		{ProfileID: "456"},
		{ProfileID: "123"},
	}, nil).Once()

	profile, err := DefaultNotificationProfile(client)

	require.NoError(t, err)
	assert.Equal(t, &api.NotificationProfile{ProfileID: "456"}, profile)
}

func TestDefaultThresholdProfile(t *testing.T) {
	client := fake.NewClient()

	client.FakeThresholdProfiles.On("List").Return(nil, errors.New("an error occurred")).Once()

	_, err := DefaultThresholdProfile(client)

	require.Equal(t, errors.New("an error occurred"), err)

	client.FakeThresholdProfiles.On("List").Return(nil, nil).Once()

	_, err = DefaultThresholdProfile(client)

	require.Equal(t, errors.New("no threshold profiles configured"), err)

	client.FakeThresholdProfiles.On("List").Return([]*api.ThresholdProfile{
		{ProfileID: "456"},
		{ProfileID: "123"},
	}, nil).Once()

	profile, err := DefaultThresholdProfile(client)

	require.NoError(t, err)
	assert.Equal(t, &api.ThresholdProfile{ProfileID: "456"}, profile)
}

func TestDefaultUserGroup(t *testing.T) {
	client := fake.NewClient()

	client.FakeUserGroups.On("List").Return(nil, errors.New("an error occurred")).Once()

	_, err := DefaultUserGroup(client)

	require.Equal(t, errors.New("an error occurred"), err)

	client.FakeUserGroups.On("List").Return(nil, nil).Once()

	_, err = DefaultUserGroup(client)

	require.Equal(t, errors.New("no usergroups configured"), err)

	client.FakeUserGroups.On("List").Return([]*api.UserGroup{
		{UserGroupID: "456"},
		{UserGroupID: "123"},
	}, nil).Once()

	userGroup, err := DefaultUserGroup(client)

	require.NoError(t, err)
	assert.Equal(t, &api.UserGroup{UserGroupID: "456"}, userGroup)
}
