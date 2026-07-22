package site24x7

import (
	"testing"

	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
	"github.com/Bonial-International-GmbH/site24x7-go/fake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCreate(t *testing.T) {
	d := userTestResourceData(t)

	c := fake.NewClient()

	u := &api.User{
		DisplayName:   "Jane Doe",
		EmailAddress:  "jane.doe@example.com",
		Role:          2,
		NotifyMedium:  []int{1, 3},
		SelectionType: 0,
		UserGroups:    []string{"123", "456"},
	}

	c.FakeUsers.On("Create", u).Return(u, nil).Once()

	require.NoError(t, userCreate(d, c))

	c.FakeUsers.On("Create", u).Return(u, apierrors.NewStatusError(500, "error")).Once()

	err := userCreate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestUserUpdate(t *testing.T) {
	d := userTestResourceData(t)
	d.SetId("112233")

	c := fake.NewClient()

	u := &api.User{
		UserID:        "112233",
		DisplayName:   "Jane Doe",
		EmailAddress:  "jane.doe@example.com",
		Role:          2,
		NotifyMedium:  []int{1, 3},
		SelectionType: 0,
		UserGroups:    []string{"123", "456"},
	}

	c.FakeUsers.On("Update", u).Return(u, nil).Once()

	require.NoError(t, userUpdate(d, c))

	c.FakeUsers.On("Update", u).Return(u, apierrors.NewStatusError(500, "error")).Once()

	err := userUpdate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestUserRead(t *testing.T) {
	d := userTestResourceData(t)
	d.SetId("112233")

	c := fake.NewClient()

	c.FakeUsers.On("Get", "112233").Return(&api.User{}, nil).Once()

	require.NoError(t, userRead(d, c))

	c.FakeUsers.On("Get", "112233").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	err := userRead(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestUserDelete(t *testing.T) {
	d := userTestResourceData(t)
	d.SetId("112233")

	c := fake.NewClient()

	c.FakeUsers.On("Delete", "112233").Return(nil).Once()

	require.NoError(t, userDelete(d, c))

	c.FakeUsers.On("Delete", "112233").Return(apierrors.NewStatusError(404, "not found")).Once()

	require.NoError(t, userDelete(d, c))
}

func TestUserExists(t *testing.T) {
	d := userTestResourceData(t)
	d.SetId("112233")

	c := fake.NewClient()

	c.FakeUsers.On("Get", "112233").Return(&api.User{}, nil).Once()

	exists, err := userExists(d, c)

	require.NoError(t, err)
	assert.True(t, exists)

	c.FakeUsers.On("Get", "112233").Return(nil, apierrors.NewStatusError(404, "not found")).Once()

	exists, err = userExists(d, c)

	require.NoError(t, err)
	assert.False(t, exists)

	c.FakeUsers.On("Get", "112233").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	exists, err = userExists(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
	assert.False(t, exists)
}

func userTestResourceData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, UserSchema, map[string]interface{}{
		"display_name":    "Jane Doe",
		"email_address":   "jane.doe@example.com",
		"role":            2,
		"notify_medium":   []interface{}{1, 3},
		"selection_type":  0,
		"user_groups":     []interface{}{"123", "456"},
		"status_iq_role":  0,
		"cloudspend_role": 0,
	})
}
