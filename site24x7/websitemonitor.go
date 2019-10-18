package site24x7

import (
	"strconv"

	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var WebsiteMonitorSchema = map[string]*schema.Schema{

	"display_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"website": {
		Type:     schema.TypeString,
		Required: true,
	},
	"check_frequency": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	},
	"http_method": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "G",
	},
	"auth_user": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"auth_pass": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"matching_keyword_value": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "", // do not auto detect
	},
	"matching_keyword_severity": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  2,
	},
	"unmatching_keyword_value": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "", // do not auto detect
	},
	"unmatching_keyword_severity": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  2,
	},
	"match_regex_value": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"match_regex_severity": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  2,
	},
	"match_case": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"user_agent": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"custom_headers": {
		Type:     schema.TypeMap,
		Optional: true,
	},
	"timeout": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  10,
	},
	"location_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"notification_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"threshold_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"monitor_groups": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"user_group_ids": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
		Computed: true,
	},
	"actions": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     schema.TypeString,
	},
	"use_name_server": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
}

func resourceSite24x7WebsiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: websiteMonitorCreate,
		Read:   websiteMonitorRead,
		Update: websiteMonitorUpdate,
		Delete: websiteMonitorDelete,
		Exists: websiteMonitorExists,

		Schema: WebsiteMonitorSchema,
	}
}

func websiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	websiteMonitor, err := resourceDataToWebsiteMonitor(d)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("match_regex_value"); ok {
		websiteMonitor.MatchRegex.Value = d.Get("match_regex_value").(string)
		websiteMonitor.MatchRegex.Severity = api.Status(d.Get("match_regex_severity").(int))
	}

	if _, ok := d.GetOk("unmatching_keyword_value"); ok {
		websiteMonitor.UnmatchingKeyword.Value = d.Get("unmatching_keyword_value").(string)
		websiteMonitor.UnmatchingKeyword.Severity = api.Status(d.Get("unmatching_keyword_severity").(int))
	}

	if _, ok := d.GetOk("matching_keyword_value"); ok {
		websiteMonitor.MatchingKeyword.Value = d.Get("matching_keyword_value").(string)
		websiteMonitor.MatchingKeyword.Severity = api.Status(d.Get("matching_keyword_severity").(int))
	}

	if websiteMonitor.LocationProfileID == "" {
		profile, err := DefaultLocationProfile(client)
		if err != nil {
			return err
		}
		websiteMonitor.LocationProfileID = profile.ProfileID
		d.Set("location_profile_id", profile.ProfileID)
	}

	if websiteMonitor.NotificationProfileID == "" {
		profile, err := DefaultNotificationProfile(client)
		if err != nil {
			return err
		}
		websiteMonitor.NotificationProfileID = profile.ProfileID
		d.Set("notification_profile_id", profile.ProfileID)
	}

	if websiteMonitor.ThresholdProfileID == "" {
		profile, err := DefaultThresholdProfile(client)
		if err != nil {
			return err
		}
		websiteMonitor.ThresholdProfileID = profile.ProfileID
		d.Set("threshold_profile_id", profile)
	}

	if len(websiteMonitor.UserGroupIDs) == 0 {
		userGroup, err := DefaultUserGroup(client)
		if err != nil {
			return err
		}
		websiteMonitor.UserGroupIDs = []string{userGroup.UserGroupID}
		d.Set("user_group_ids", []string{userGroup.UserGroupID})
	}

	websiteMonitor, err = client.Monitors().Create(websiteMonitor)
	if err != nil {
		return err
	}

	d.SetId(websiteMonitor.MonitorID)

	return nil
}

func websiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	websiteMonitor, err := client.Monitors().Get(d.Id())
	if err != nil {
		return err
	}

	updateWebsiteMonitorResourceData(d, websiteMonitor)

	return nil
}

func websiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	websiteMonitor, err := resourceDataToWebsiteMonitor(d)
	if err != nil {
		return err
	}

	websiteMonitor, err = client.Monitors().Update(websiteMonitor)
	if err != nil {
		return err
	}

	d.SetId(websiteMonitor.MonitorID)

	return nil
}

func websiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	return client.Monitors().Delete(d.Id())
}

func websiteMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(site24x7.Client)

	_, err := client.Monitors().Get(d.Id())
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDataToWebsiteMonitor(d *schema.ResourceData) (*api.Monitor, error) {
	customHeaders := []api.Header{}
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders = append(customHeaders, api.Header{Name: k, Value: v.(string)})
	}

	var userGroupIDs []string
	for _, id := range d.Get("user_group_ids").([]interface{}) {
		userGroupIDs = append(userGroupIDs, id.(string))
	}

	var monitorGroups []string
	for _, group := range d.Get("monitor_groups").([]interface{}) {
		monitorGroups = append(monitorGroups, group.(string))
	}

	i := 0
	actionRefs := make([]api.ActionRef, len(d.Get("action_ids").(map[string]interface{})))
	for k, v := range d.Get("action_ids").(map[string]interface{}) {
		tmp, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		alertType := api.Status(tmp)
		actionRefs[i] = api.ActionRef{ActionID: v.(string), AlertType: alertType}
		i++
	}

	return &api.Monitor{
		MonitorID:             d.Id(),
		DisplayName:           d.Get("display_name").(string),
		Type:                  "URL",
		Website:               d.Get("website").(string),
		CheckFrequency:        strconv.Itoa(d.Get("check_frequency").(int)),
		HTTPMethod:            d.Get("http_method").(string),
		AuthUser:              d.Get("auth_user").(string),
		AuthPass:              d.Get("auth_pass").(string),
		MatchCase:             d.Get("match_case").(bool),
		UserAgent:             d.Get("user_agent").(string),
		CustomHeaders:         customHeaders,
		Timeout:               d.Get("timeout").(int),
		LocationProfileID:     d.Get("location_profile_id").(string),
		NotificationProfileID: d.Get("notification_profile_id").(string),
		ThresholdProfileID:    d.Get("threshold_profile_id").(string),
		MonitorGroups:         monitorGroups,
		UserGroupIDs:          userGroupIDs,
		ActionIDs:             actionRefs,
		UseNameServer:         d.Get("use_name_server").(bool),
	}, nil
}

func updateWebsiteMonitorResourceData(d *schema.ResourceData, monitor *api.Monitor) {
	d.Set("display_name", monitor.DisplayName)
	d.Set("type", monitor.Type)
	d.Set("auth_user", monitor.AuthUser)
	d.Set("auth_pass", monitor.AuthPass)
	if monitor.MatchingKeyword != nil {
		d.Set("matching_keyword_value", monitor.MatchingKeyword.Value)
		d.Set("matching_keyword_severity", monitor.MatchingKeyword.Severity)
	}
	if monitor.UnmatchingKeyword != nil {
		d.Set("unmatching_keyword_value", monitor.UnmatchingKeyword.Value)
		d.Set("unmatching_keyword_severity", monitor.UnmatchingKeyword.Severity)
	}
	if monitor.MatchRegex != nil {
		d.Set("match_regex_value", monitor.MatchRegex.Value)
		d.Set("match_regex_severity", monitor.MatchRegex.Severity)
	}
	d.Set("match_case", monitor.MatchCase)
	d.Set("user_agent", monitor.UserAgent)
	d.Set("custom_headers", monitor.CustomHeaders)
	d.Set("timeout", monitor.Timeout)
	d.Set("location_profile_id", monitor.LocationProfileID)
	d.Set("notification_profile_id", monitor.NotificationProfileID)
	d.Set("threshold_profile_id", monitor.ThresholdProfileID)
	d.Set("monitor_groups", monitor.MonitorGroups)
	d.Set("user_group_ids", monitor.UserGroupIDs)
	d.Set("action_ids", monitor.ActionIDs)
	d.Set("use_name_server", monitor.UseNameServer)
}
