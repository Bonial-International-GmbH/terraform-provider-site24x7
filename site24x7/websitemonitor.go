package site24x7

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSite24x7WebsiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: websiteMonitorCreate,
		Read:   websiteMonitorRead,
		Update: websiteMonitorUpdate,
		Delete: monitorDelete,
		Exists: monitorExists,

		Schema: map[string]*schema.Schema{
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
		},
	}
}

type WebsiteMonitor struct {
	MonitorID             string            `json:"monitor_id,omitempty"`
	DisplayName           string            `json:"display_name"`
	Type                  string            `json:"type"`
	Website               string            `json:"website"`
	CheckFrequency        string            `json:"check_frequency"`
	HTTPMethod            string            `json:"http_method"`
	AuthUser              string            `json:"auth_user"`
	AuthPass              string            `json:"auth_pass"`
	MatchingKeyword       *ValueAndSeverity `json:"matching_keyword,omitempty"`
	UnmatchingKeyword     *ValueAndSeverity `json:"unmatching_keyword,omitempty"`
	MatchRegex            *ValueAndSeverity `json:"match_regex,omitempty"`
	MatchCase             bool              `json:"match_case"`
	UserAgent             string            `json:"user_agent"`
	CustomHeaders         []Header          `json:"custom_headers"`
	Timeout               int               `json:"timeout"`
	LocationProfileID     string            `json:"location_profile_id"`
	NotificationProfileID string            `json:"notification_profile_id"`
	ThresholdProfileID    string            `json:"threshold_profile_id"`
	MonitorGroups         []string          `json:"monitor_groups"`
	UserGroupIDs          []string          `json:"user_group_ids"`
	ActionIDs             []ActionRef       `json:"action_ids,omitempty"`
	UseNameServer         bool              `json:"use_name_server"`
}

func websiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	return websiteMonitorCreateOrUpdate(http.MethodPost, "https://www.site24x7.com/api/monitors", http.StatusCreated, d, meta)
}

func websiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	return websiteMonitorCreateOrUpdate(http.MethodPut, "https://www.site24x7.com/api/monitors/"+d.Id(), http.StatusOK, d, meta)
}

func websiteMonitorCreateOrUpdate(method, url string, expectedResponseStatus int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	customHeaders := []Header{}
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders = append(customHeaders, Header{Name: k, Value: v.(string)})
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
	actionRefs := make([]ActionRef, len(d.Get("actions").(map[string]interface{})))
	for k, v := range d.Get("actions").(map[string]interface{}) {
		tmp, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		alertType := Status(tmp)
		actionRefs[i] = ActionRef{ActionID: v.(string), AlertType: alertType}
		i++
	}

	m := &WebsiteMonitor{
		DisplayName:           d.Get("display_name").(string),
		Type:                  "URL",
		Website:               d.Get("website").(string),
		CheckFrequency:        strconv.Itoa(d.Get("check_frequency").(int)),
		HTTPMethod:            d.Get("http_method").(string),
		AuthUser:              d.Get("auth_user").(string),
		AuthPass:              d.Get("auth_pass").(string),
		MatchingKeyword:       new(ValueAndSeverity),
		UnmatchingKeyword:     new(ValueAndSeverity),
		MatchRegex:            new(ValueAndSeverity),
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
	}

	if _, ok := d.GetOk("match_regex_value"); ok {
		m.MatchRegex.Value = d.Get("match_regex_value").(string)
		m.MatchRegex.Severity = Status(d.Get("match_regex_severity").(int))
	} else {
		m.MatchRegex = nil
	}
	if _, ok := d.GetOk("unmatching_keyword_value"); ok {
		m.UnmatchingKeyword.Value = d.Get("unmatching_keyword_value").(string)
		m.UnmatchingKeyword.Severity = Status(d.Get("unmatching_keyword_severity").(int))
	} else {
		m.UnmatchingKeyword = nil
	}
	if _, ok := d.GetOk("matching_keyword_value"); ok {
		m.MatchingKeyword.Value = d.Get("matching_keyword_value").(string)
		m.MatchingKeyword.Severity = Status(d.Get("matching_keyword_severity").(int))
	} else {
		m.MatchingKeyword = nil
	}

	if m.LocationProfileID == "" {
		id, err := defaultLocationProfile(client)
		if err != nil {
			return err
		}
		m.LocationProfileID = id
		d.Set("location_profile_id", id)
	}
	if m.NotificationProfileID == "" {
		id, err := defaultNotificationProfile(client)
		if err != nil {
			return err
		}
		m.NotificationProfileID = id
		d.Set("notification_profile_id", id)
	}
	if m.ThresholdProfileID == "" {
		id, err := defaultThresholdProfile(client)
		if err != nil {
			return err
		}
		m.ThresholdProfileID = id
		d.Set("threshold_profile_id", id)
	}
	if len(m.UserGroupIDs) == 0 {
		id, err := defaultUserGroup(client)
		if err != nil {
			return err
		}
		m.UserGroupIDs = []string{id}
		d.Set("user_group_ids", []string{id})
	}

	log.Printf("monitor: %+v", m)
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseStatus {
		return parseAPIError(resp.Body)
	}

	var apiResp struct {
		Data struct {
			MonitorID string `json:"monitor_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}
	d.SetId(apiResp.Data.MonitorID)
	// can't update the rest of the data here, because the response format is broken

	return nil
}

func websiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	var apiResp struct {
		Data WebsiteMonitor `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/monitors/"+d.Id(), &apiResp); err != nil {
		return err
	}
	updateWebsiteMonitorResourceData(d, &apiResp.Data)

	return nil
}

func updateWebsiteMonitorResourceData(d *schema.ResourceData, m *WebsiteMonitor) {
	d.Set("display_name", m.DisplayName)
	d.Set("website", m.Website)
	d.Set("check_frequency", m.CheckFrequency)
	d.Set("timeout", m.Timeout)
	d.Set("http_method", m.HTTPMethod)
	d.Set("auth_user", m.AuthUser)
	d.Set("auth_pass", m.AuthPass)
	if m.MatchingKeyword != nil {
		d.Set("matching_keyword_value", m.MatchingKeyword.Value)
		d.Set("matching_keyword_severity", m.MatchingKeyword.Severity)
	}
	if m.UnmatchingKeyword != nil {
		d.Set("unmatching_keyword_value", m.UnmatchingKeyword.Value)
		d.Set("unmatching_keyword_severity", m.UnmatchingKeyword.Severity)
	}
	if m.MatchRegex != nil {
		d.Set("match_regex_value", m.MatchRegex.Value)
		d.Set("match_regex_severity", m.MatchRegex.Severity)
	}
	d.Set("match_case", m.MatchCase)
	d.Set("user_agent", m.UserAgent)
	customHeaders := make(map[string]interface{})
	for _, h := range m.CustomHeaders {
		if h.Name == "" {
			continue
		}
		customHeaders[h.Name] = h.Value
	}
	d.Set("custom_headers", customHeaders)
	d.Set("location_profile_id", m.LocationProfileID)
	d.Set("notification_profile_id", m.NotificationProfileID)
	d.Set("threshold_profile_id", m.ThresholdProfileID)
	d.Set("monitor_groups", m.MonitorGroups)
	d.Set("user_group_ids", m.UserGroupIDs)
	actions := make(map[string]string)
	for _, r := range m.ActionIDs {
		actions[fmt.Sprintf("%d", r.AlertType)] = r.ActionID
	}
	d.Set("actions", actions)
	d.Set("use_name_server", m.UseNameServer)
}

func fetchWebsiteMonitorExists(client *http.Client, id string) (bool, error) {
	resp, err := client.Get("https://www.site24x7.com/api/monitors/" + id)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, parseAPIError(resp.Body)
	}
}
