package site24x7

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSite24x7Action() *schema.Resource {
	return &schema.Resource{
		Create: actionCreate,
		Read:   actionRead,
		Update: actionUpdate,
		Delete: actionDelete,
		Exists: actionExists,

		Schema: map[string]*schema.Schema{
			"custom_parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "P",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"requires_authentication": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"send_custom_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"send_in_json_format": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"send_incident_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type Action struct {
	Action_method            string `json:"action_method"`
	Action_name              string `json:"action_name"`
	Action_timeout           int    `json:"action_timeout"`
	Action_type              int    `json:"action_type"`
	Action_url               string `json:"action_url"`
	Custom_parameters        string `json:"custom_parameters"`
	Requires_authentication  bool   `json:"requires_authentication"`
	Send_custom_parameters   bool   `json:"send_custom_parameters"`
	Send_in_json_format      bool   `json:"send_in_json_format"`
	Send_incident_parameters bool   `json:"send_incident_parameters"`
}

type ActionID struct {
	Action_id string `json:"action_id"`
}

func actionCreate(d *schema.ResourceData, meta interface{}) error {
	return actionCreateOrUpdate(http.MethodPost, "https://www.site24x7.com/api/it_automation", http.StatusCreated, d, meta)
}

func actionUpdate(d *schema.ResourceData, meta interface{}) error {
	return actionCreateOrUpdate(http.MethodPut, "https://www.site24x7.com/api/it_automation/"+d.Id(), http.StatusOK, d, meta)
}

func actionCreateOrUpdate(method, url string, expectedResponseStatus int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	m := &Action{
		Action_method:            d.Get("method").(string),
		Action_name:              d.Get("name").(string),
		Action_timeout:           d.Get("timeout").(int),
		Action_type:              d.Get("type").(int),
		Action_url:               d.Get("url").(string),
		Custom_parameters:        d.Get("custom_parameters").(string),
		Requires_authentication:  d.Get("requires_authentication").(bool),
		Send_custom_parameters:   d.Get("send_custom_parameters").(bool),
		Send_in_json_format:      d.Get("send_in_json_format").(bool),
		Send_incident_parameters: d.Get("send_incident_parameters").(bool),
	}

	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.Marshal(%#v) failed: %s", m, err)
	}
	log.Printf("action body: %s", body)

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest(%s, %s, %s) failed: %s", method, url, body, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do(%s, %s, %s) failed: %s", method, url, body, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseStatus {
		return parseAPIError(resp.Body, fmt.Sprintf("request(%s, %s, %s)", method, url, body))
	}

	var apiResp struct {
		Data struct {
			ActionID string `json:"action_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("json decoding of response failed: (%s)", err)
	}
	d.SetId(apiResp.Data.ActionID)

	return nil
}

func actionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	var apiResp struct {
		Data Action `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/it_automation/"+d.Id(), &apiResp); err != nil {
		return err
	}
	updateActionResourceData(d, &apiResp.Data)

	return nil
}

func updateActionResourceData(d *schema.ResourceData, m *Action) {
	d.Set("method", m.Action_method)
	d.Set("name", m.Action_name)
	d.Set("timeout", m.Action_timeout)
	d.Set("type", m.Action_type)
	d.Set("url", m.Action_url)
	d.Set("custom_parameters", m.Custom_parameters)
	d.Set("requires_authentication", m.Requires_authentication)
	d.Set("send_custom_parameters", m.Send_custom_parameters)
	d.Set("send_in_json_format", m.Send_in_json_format)
	d.Set("send_incident_parameters", m.Send_incident_parameters)
}

func actionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	req, err := http.NewRequest(http.MethodDelete, "https://www.site24x7.com/api/it_automation/"+d.Id(), nil)
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

func actionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return fetchActionExists(meta.(*http.Client), d.Id())
}

func fetchActionExists(client *http.Client, id string) (bool, error) {
	var apiResp struct {
		Data []ActionID `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/it_automation", &apiResp); err != nil {
		return false, err
	}
	for _, v := range apiResp.Data {
		if v.Action_id == id {
			return true, nil
		}
	}

	return false, nil
}
