package site24x7

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSite24x7MonitorGroup() *schema.Resource {
	return &schema.Resource{
		Create: monitorGroupCreate,
		Read:   monitorGroupRead,
		Update: monitorGroupUpdate,
		Delete: monitorGroupDelete,
		Exists: monitorGroupExists,

		Schema: map[string]*schema.Schema{
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type MonitorGroup struct {
	MonitorGroupID string `json:"group_id,omitempty"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
}

func monitorGroupCreate(d *schema.ResourceData, meta interface{}) error {
	return monitorGroupCreateOrUpdate(http.MethodPost, "https://www.site24x7.com/api/monitor_groups", http.StatusCreated, d, meta)
}

func monitorGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return monitorGroupCreateOrUpdate(http.MethodPut, "https://www.site24x7.com/api/monitor_groups/"+d.Id(), http.StatusOK, d, meta)
}

func monitorGroupCreateOrUpdate(method, url string, expectedResponseStatus int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	m := &MonitorGroup{
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
	}

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	log.Printf("group body: %s", body)

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
			GroupID string `json:"group_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}
	d.SetId(apiResp.Data.GroupID)

	return nil
}

func monitorGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	var apiResp struct {
		Data MonitorGroup `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/monitor_groups/"+d.Id(), &apiResp); err != nil {
		return err
	}
	updateMonitorGroupResourceData(d, &apiResp.Data)

	return nil
}

func updateMonitorGroupResourceData(d *schema.ResourceData, m *MonitorGroup) {
	d.Set("display_name", m.DisplayName)
	d.Set("description", m.Description)
}

func monitorGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	req, err := http.NewRequest(http.MethodDelete, "https://www.site24x7.com/api/monitor_groups/"+d.Id(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp.Body)
	}

	return nil
}

func monitorGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return fetchMonitorGroupExists(meta.(*http.Client), d.Id())
}

func fetchMonitorGroupExists(client *http.Client, id string) (bool, error) {
	var apiResp struct {
		Data []MonitorGroup `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/monitor_groups", &apiResp); err != nil {
		return false, err
	}
	for _, v := range apiResp.Data {
		if v.MonitorGroupID == id {
			return true, nil
		}
	}

	return false, nil
}
