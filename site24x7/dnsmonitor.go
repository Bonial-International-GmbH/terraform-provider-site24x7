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

func resourceSite24x7DnsMonitor() *schema.Resource {
	return &schema.Resource{
		Create: dnsMonitorCreate,
		Read:   dnsMonitorRead,
		Update: dnsMonitorUpdate,
		Delete: monitorDelete,
		Exists: monitorExists,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"dns_host": {
				Type:     schema.TypeString,
				Required: true,
			},

			"dns_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  53,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"check_frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
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

			"user_group_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

			"actions": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

type DnsMonitor struct {
	MonitorID             string      `json:"monitor_id,omitempty"`
	DisplayName           string      `json:"display_name"`
	Type                  string      `json:"type"`
	DnsHost               string      `json:"dns_host"`
	DnsPort               string      `json:"dns_port"`
	DomainName            string      `json:"domain_name"`
	CheckFrequency        string      `json:"check_frequency"`
	Timeout               int         `json:"timeout"`
	LocationProfileID     string      `json:"location_profile_id"`
	NotificationProfileID string      `json:"notification_profile_id"`
	ThresholdProfileID    string      `json:"threshold_profile_id"`
	UserGroupIDs          []string    `json:"user_group_ids"`
	MonitorGroups         []string    `json:"monitor_groups"`
	ActionIDs             []ActionRef `json:"action_ids,omitempty"`
}

func dnsMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	return dnsMonitorCreateOrUpdate(http.MethodPost, "https://www.site24x7.com/api/monitors", http.StatusCreated, d, meta)
}

func dnsMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	return dnsMonitorCreateOrUpdate(http.MethodPut, "https://www.site24x7.com/api/monitors/"+d.Id(), http.StatusOK, d, meta)
}

func dnsMonitorCreateOrUpdate(method, url string, expectedResponseStatus int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

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

	m := &DnsMonitor{
		DisplayName:           d.Get("display_name").(string),
		Type:                  "DNS",
		DnsHost:               d.Get("dns_host").(string),
		DnsPort:               strconv.Itoa(d.Get("dns_port").(int)),
		DomainName:            d.Get("domain_name").(string),
		CheckFrequency:        strconv.Itoa(d.Get("check_frequency").(int)),
		Timeout:               d.Get("timeout").(int),
		LocationProfileID:     d.Get("location_profile_id").(string),
		NotificationProfileID: d.Get("notification_profile_id").(string),
		ThresholdProfileID:    d.Get("threshold_profile_id").(string),
		UserGroupIDs:          userGroupIDs,
		MonitorGroups:         monitorGroups,
		ActionIDs:             actionRefs,
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

func dnsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	var apiResp struct {
		Data DnsMonitor `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/monitors/"+d.Id(), &apiResp); err != nil {
		return err
	}
	updateDnsMonitorResourceData(d, &apiResp.Data)

	return nil
}

func updateDnsMonitorResourceData(d *schema.ResourceData, m *DnsMonitor) {
	d.Set("display_name", m.DisplayName)
	d.Set("dns_host", m.DnsHost)
	d.Set("dns_port", m.DnsPort)
	d.Set("check_frequency", m.CheckFrequency)
	d.Set("timeout", m.Timeout)
	d.Set("location_profile_id", m.LocationProfileID)
	d.Set("notification_profile_id", m.NotificationProfileID)
	d.Set("threshold_profile_id", m.ThresholdProfileID)
	d.Set("user_group_ids", m.UserGroupIDs)
	d.Set("monitor_groups", m.MonitorGroups)
	actions := make(map[string]string)
	for _, r := range m.ActionIDs {
		actions[fmt.Sprintf("%d", r.AlertType)] = r.ActionID
	}
	d.Set("actions", actions)
}
