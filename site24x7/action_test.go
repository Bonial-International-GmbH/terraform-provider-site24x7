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

func TestActionCreate(t *testing.T) {
	d := actionTestResourceData(t)

	c := fake.NewClient()

	a := &api.ITAutomation{
		ActionName:             "foobar",
		ActionMethod:           "P",
		CustomParameters:       "foobarbaz",
		SendCustomParameters:   true,
		SendInJsonFormat:       true,
		SendIncidentParameters: false,
		ActionTimeout:          30,
		ActionUrl:              "https://example.com",
		ActionType:             1,
	}

	c.FakeITAutomations.On("Create", a).Return(a, nil).Once()

	require.NoError(t, actionCreate(d, c))

	c.FakeITAutomations.On("Create", a).Return(a, apierrors.NewStatusError(500, "error")).Once()

	err := actionCreate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestActionUpdate(t *testing.T) {
	d := actionTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	a := &api.ITAutomation{
		ActionID:               "123",
		ActionName:             "foobar",
		ActionMethod:           "P",
		CustomParameters:       "foobarbaz",
		SendCustomParameters:   true,
		SendInJsonFormat:       true,
		SendIncidentParameters: false,
		ActionTimeout:          30,
		ActionUrl:              "https://example.com",
		ActionType:             1,
	}

	c.FakeITAutomations.On("Update", a).Return(a, nil).Once()

	require.NoError(t, actionUpdate(d, c))

	c.FakeITAutomations.On("Update", a).Return(a, apierrors.NewStatusError(500, "error")).Once()

	err := actionUpdate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestActionRead(t *testing.T) {
	d := actionTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeITAutomations.On("Get", "123").Return(&api.ITAutomation{}, nil).Once()

	require.NoError(t, actionRead(d, c))

	c.FakeITAutomations.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	err := actionRead(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestActionDelete(t *testing.T) {
	d := actionTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeITAutomations.On("Delete", "123").Return(nil).Once()

	require.NoError(t, actionDelete(d, c))

	c.FakeITAutomations.On("Delete", "123").Return(apierrors.NewStatusError(404, "not found")).Once()

	require.NoError(t, actionDelete(d, c))
}

func TestActionExists(t *testing.T) {
	d := actionTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeITAutomations.On("Get", "123").Return(&api.ITAutomation{}, nil).Once()

	exists, err := actionExists(d, c)

	require.NoError(t, err)
	assert.True(t, exists)

	c.FakeITAutomations.On("Get", "123").Return(nil, apierrors.NewStatusError(404, "not found")).Once()

	exists, err = actionExists(d, c)

	require.NoError(t, err)
	assert.False(t, exists)

	c.FakeITAutomations.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	exists, err = actionExists(d, c)

	require.Equal(t, apierrors.NewStatusError(500, "error"), err)
	assert.False(t, exists)
}

func actionTestResourceData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, ActionSchema, map[string]interface{}{
		"name":                     "foobar",
		"method":                   "P",
		"custom_parameters":        "foobarbaz",
		"send_custom_parameters":   true,
		"send_in_json_format":      true,
		"send_incident_parameters": false,
		"timeout":                  30,
		"url":                      "https://example.com",
		"type":                     1,
	})
}
