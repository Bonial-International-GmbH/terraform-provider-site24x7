package site24x7

import (
	"errors"
	"testing"

	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
	"github.com/Bonial-International-GmbH/site24x7-go/fake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWebsiteMonitorCreate(t *testing.T) {
	tests := []struct {
		name                 string
		setup                func(t *testing.T, c *fake.Client)
		resourceDataProvider func(t *testing.T) *schema.ResourceData
		expectedErr          error
		validate             func(t *testing.T, c *fake.Client)
	}{
		{
			name: "create simple monitor",
			setup: func(t *testing.T, c *fake.Client) {
				a := &api.Monitor{
					DisplayName:           "foo",
					Type:                  "URL",
					Website:               "www.test.tld",
					CheckFrequency:        "1",
					HTTPMethod:            "G",
					Timeout:               10,
					LocationProfileID:     "456",
					NotificationProfileID: "789",
					ThresholdProfileID:    "012",
					UseNameServer:         true,
					UserGroupIDs:          []string{"123"},
					CustomHeaders:         []api.Header{},
					ActionIDs:             []api.ActionRef{},
				}

				c.FakeMonitors.On("Create", a).Return(a, nil).Once()
			},
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
		},
		{
			name: "passes through create monitor error",
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeMonitors.On("Create", mock.Anything).Return(nil, apierrors.NewStatusError(500, "server error")).Once()
			},
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
			expectedErr: apierrors.NewStatusError(500, "server error"),
		},
		{
			name: "somebody tampered with the statefile",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"actions": map[string]interface{}{
						"this-will-cause-an-error": "123action",
					},
					"user_group_ids": []interface{}{"123"},
				})
			},
			validate: func(t *testing.T, c *fake.Client) {
				assert.Len(t, c.FakeMonitors.Calls, 0)
			},
			expectedErr: errors.New(`strconv.Atoi: parsing "this-will-cause-an-error": invalid syntax`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := fake.NewClient()
			d := test.resourceDataProvider(t)

			if test.setup != nil {
				test.setup(t, c)
			}

			err := websiteMonitorCreate(d, c)
			if test.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if test.validate != nil {
				test.validate(t, c)
			}
		})
	}
}

func TestWebsiteMonitorUpdate(t *testing.T) {
	tests := []struct {
		name                 string
		setup                func(t *testing.T, c *fake.Client)
		resourceDataProvider func(t *testing.T) *schema.ResourceData
		expectedErr          error
		validate             func(t *testing.T, c *fake.Client)
	}{
		{
			name: "updates simple monitor",
			setup: func(t *testing.T, c *fake.Client) {
				a := &api.Monitor{
					MonitorID:             "123",
					DisplayName:           "foo",
					Type:                  "URL",
					Website:               "www.test.tld",
					CheckFrequency:        "1",
					HTTPMethod:            "G",
					Timeout:               10,
					LocationProfileID:     "456",
					NotificationProfileID: "789",
					ThresholdProfileID:    "012",
					UseNameServer:         true,
					UserGroupIDs:          []string{"123"},
					CustomHeaders: []api.Header{
						{
							Name:  "Accept",
							Value: "application/json",
						},
						{
							Name:  "Cache-Control",
							Value: "nocache",
						},
					},
					ActionIDs: []api.ActionRef{
						{
							ActionID:  "123action",
							AlertType: 1,
						},
						{
							ActionID:  "234action",
							AlertType: 5,
						},
					},
					UnmatchingKeyword: &api.ValueAndSeverity{
						Value:    "foo",
						Severity: 2,
					},
					MatchingKeyword: &api.ValueAndSeverity{
						Value:    "bar",
						Severity: 2,
					},
					MatchRegex: &api.ValueAndSeverity{
						Value:    ".*",
						Severity: 2,
					},
				}

				c.FakeMonitors.On("Update", a).Return(a, nil).Once()
			},
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				rd := schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
					"custom_headers": map[string]interface{}{
						"Accept":        "application/json",
						"Cache-Control": "nocache",
					},
					"actions": map[string]interface{}{
						"1": "123action",
						"5": "234action",
					},
					"unmatching_keyword_value": "foo",
					"matching_keyword_value":   "bar",
					"match_regex_value":        ".*",
				})

				rd.SetId("123")

				return rd
			},
		},
		{
			name: "passes through create monitor error",
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeMonitors.On("Update", mock.Anything).Return(nil, apierrors.NewStatusError(500, "server error")).Once()
			},
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				rd := schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})

				rd.SetId("123")

				return rd
			},
			expectedErr: apierrors.NewStatusError(500, "server error"),
		},
		{
			name: "somebody tampered with the statefile",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"actions": map[string]interface{}{
						"this-will-cause-an-error": "123action",
					},
					"user_group_ids": []interface{}{"123"},
				})
			},
			validate: func(t *testing.T, c *fake.Client) {
				assert.Len(t, c.FakeMonitors.Calls, 0)
			},
			expectedErr: errors.New(`strconv.Atoi: parsing "this-will-cause-an-error": invalid syntax`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := fake.NewClient()
			d := test.resourceDataProvider(t)

			if test.setup != nil {
				test.setup(t, c)
			}

			err := websiteMonitorUpdate(d, c)
			if test.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if test.validate != nil {
				test.validate(t, c)
			}
		})
	}
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

	c.FakeMonitors.On("Delete", "123").Return(apierrors.NewStatusError(404, "not found")).Once()

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
			"Header Name": "testheader",
			"cache":       "nocache",
		},
		"timeout":                 120,
		"location_profile_id":     "456",
		"notification_profile_id": "789",
		"threshold_profile_id":    "012",
		"monitor_groups": []interface{}{
			"234",
			"567",
		},
		"user_group_ids": []interface{}{
			"123",
			"456",
		},
		"use_name_server": true,
		"actions": map[string]interface{}{
			"1": "123action",
			"5": "234action",
		},
	})
}

func TestResourceDataToWebsiteMonitor(t *testing.T) {
	tests := []struct {
		name                 string
		setup                func(t *testing.T, c *fake.Client)
		resourceDataProvider func(t *testing.T) *schema.ResourceData
		expected             *api.Monitor
		expectedErr          error
	}{
		{
			name: "simple",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "456",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
			expected: &api.Monitor{
				DisplayName:           "foo",
				Type:                  "URL",
				Website:               "www.test.tld",
				CheckFrequency:        "1",
				HTTPMethod:            "G",
				Timeout:               10,
				LocationProfileID:     "456",
				NotificationProfileID: "789",
				ThresholdProfileID:    "012",
				UseNameServer:         true,
				UserGroupIDs:          []string{"123"},
				CustomHeaders:         []api.Header{},
				ActionIDs:             []api.ActionRef{},
			},
		},
		{
			name: "fetches default location profile if not set",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeLocationProfiles.On("List").Return([]*api.LocationProfile{
					{ProfileID: "345"},
				}, nil)
			},
			expected: &api.Monitor{
				DisplayName:           "foo",
				Type:                  "URL",
				Website:               "www.test.tld",
				CheckFrequency:        "1",
				HTTPMethod:            "G",
				Timeout:               10,
				LocationProfileID:     "345",
				NotificationProfileID: "789",
				ThresholdProfileID:    "012",
				UseNameServer:         true,
				UserGroupIDs:          []string{"123"},
				CustomHeaders:         []api.Header{},
				ActionIDs:             []api.ActionRef{},
			},
		},
		{
			name: "returns error if lookup of default location profile fails",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"notification_profile_id": "789",
					"threshold_profile_id":    "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeLocationProfiles.On("List").Return(nil, apierrors.NewStatusError(503, "service unavailable"))
			},
			expectedErr: apierrors.NewStatusError(503, "service unavailable"),
		},
		{
			name: "fetches default notification profile if not set",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":         "foo",
					"type":                 "URL",
					"website":              "www.test.tld",
					"location_profile_id":  "789",
					"threshold_profile_id": "012",
					"user_group_ids":       []interface{}{"123"},
				})
			},
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeNotificationProfiles.On("List").Return([]*api.NotificationProfile{
					{ProfileID: "345"},
				}, nil)
			},
			expected: &api.Monitor{
				DisplayName:           "foo",
				Type:                  "URL",
				Website:               "www.test.tld",
				CheckFrequency:        "1",
				HTTPMethod:            "G",
				Timeout:               10,
				LocationProfileID:     "789",
				NotificationProfileID: "345",
				ThresholdProfileID:    "012",
				UseNameServer:         true,
				UserGroupIDs:          []string{"123"},
				CustomHeaders:         []api.Header{},
				ActionIDs:             []api.ActionRef{},
			},
		},
		{
			name: "fetches default threshold profile if not set",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "789",
					"notification_profile_id": "012",
					"user_group_ids":          []interface{}{"123"},
				})
			},
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeThresholdProfiles.On("List").Return([]*api.ThresholdProfile{
					{ProfileID: "345"},
				}, nil)
			},
			expected: &api.Monitor{
				DisplayName:           "foo",
				Type:                  "URL",
				Website:               "www.test.tld",
				CheckFrequency:        "1",
				HTTPMethod:            "G",
				Timeout:               10,
				LocationProfileID:     "789",
				NotificationProfileID: "012",
				ThresholdProfileID:    "345",
				UseNameServer:         true,
				UserGroupIDs:          []string{"123"},
				CustomHeaders:         []api.Header{},
				ActionIDs:             []api.ActionRef{},
			},
		},
		{
			name: "fetches default threshold profile if not set",
			resourceDataProvider: func(t *testing.T) *schema.ResourceData {
				return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
					"display_name":            "foo",
					"type":                    "URL",
					"website":                 "www.test.tld",
					"location_profile_id":     "789",
					"notification_profile_id": "012",
					"threshold_profile_id":    "345",
				})
			},
			setup: func(t *testing.T, c *fake.Client) {
				c.FakeUserGroups.On("List").Return([]*api.UserGroup{
					{UserGroupID: "345"},
				}, nil)
			},
			expected: &api.Monitor{
				DisplayName:           "foo",
				Type:                  "URL",
				Website:               "www.test.tld",
				CheckFrequency:        "1",
				HTTPMethod:            "G",
				Timeout:               10,
				LocationProfileID:     "789",
				NotificationProfileID: "012",
				ThresholdProfileID:    "345",
				UseNameServer:         true,
				UserGroupIDs:          []string{"345"},
				CustomHeaders:         []api.Header{},
				ActionIDs:             []api.ActionRef{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := fake.NewClient()
			d := test.resourceDataProvider(t)

			if test.setup != nil {
				test.setup(t, c)
			}

			monitor, err := resourceDataToWebsiteMonitor(d, c)
			if test.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, monitor)
			}
		})
	}
}
