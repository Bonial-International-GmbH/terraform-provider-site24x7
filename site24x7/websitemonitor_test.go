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

func TestWebsiteMonitorCreate(t *testing.T) {
	d := monitorTestResourceData(t)

	c := fake.NewClient()

	// TODO switch monitorType
	monitorType := "URL"

	a := &api.Monitor{
		DisplayName:    "foo",
		Type:           monitorType,
		Website:        "www.test.tld",
		CheckFrequency: "60",
		HTTPMethod:     "P",
		AuthUser:       "username",
		AuthPass:       "password",
		MatchCase:      true,
		UserAgent:      "firefox",
		CustomHeaders: [{
			"name":  "Header Name",
			"value": "testheader",
		},
		{
			"name": "cache",
			"value": "nocache",
		}
		],
		Timeout:               120,
		LocationProfileID:     "456",
		NotificationProfileID: "789",
		ThresholdProfileID:    "012",
		//		MonitorGroups: []string{
		//			"234",
		//			"567",
		//		},
		//		UserGroupIDs: []string{
		//			"123",
		//			"456",
		//		},
		UseNameServer: true,
		//	ActionIDs:     []api.ActionRef{},
	}

	c.FakeMonitors.On("Create", a).Return(a, nil).Once()

	require.NoError(t, websiteMonitorCreate(d, c))

	c.FakeMonitors.On("Create", a).Return(a, apierrors.NewStatusError(500, "error")).Once()

	err := websiteMonitorCreate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestWebsiteMonitorRead(t *testing.T) {
	d := monitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeMonitors.On("Get", "123").Return(&api.Monitor{}, nil).Once()

	require.NoError(t, websiteMonitorRead(d, c))

	c.FakeMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	err := websiteMonitorRead(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestWebsiteMonitorDelete(t *testing.T) {
	d := monitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeMonitors.On("Delete", "123").Return(nil).Once()

	require.NoError(t, websiteMonitorDelete(d, c))
}

func TestWebsiteMonitorExists(t *testing.T) {
	d := monitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeMonitors.On("Get", "123").Return(&api.Monitor{}, nil).Once()

	exists, err := websiteMonitorExists(d, c)

	require.NoError(t, err)
	assert.True(t, exists)

	c.FakeMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(404, "not found")).Once()

	exists, err = websiteMonitorExists(d, c)

	require.NoError(t, err)
	assert.False(t, exists)

	c.FakeMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	exists, err = websiteMonitorExists(d, c)

	require.Equal(t, apierrors.NewStatusError(500, "error"), err)
	assert.False(t, exists)
}

func monitorTestResourceData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
		"display_name":    "foo",
		"type":            "URL",
		"website":         "www.test.tld",
		"check_frequency": "60",
		"http_method":     "P",
		"auth_user":       "username",
		"auth_pass":       "password",
		"match_case":      true,
		"user_agent":      "firefox",
		"custom_headers": map[string]interface{}{
			"name":  "Accept-Encoding",
			"value": "gzip",
		},
		"timeout":                 120,
		"location_profile_id":     "456",
		"notification_profile_id": "789",
		"threshold_profile_id":    "012",
		"monitor_groups": []string{
			"234",
			"567",
		},
		"user_group_ids": []string{
			"123",
			"456",
		},
		"use_name_server": true,
		"action_ids": map[string]interface{}{
			"action_id":  "123",
			"alert_type": "5",
		},
	})
}
