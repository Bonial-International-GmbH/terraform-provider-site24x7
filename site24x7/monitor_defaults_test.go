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
