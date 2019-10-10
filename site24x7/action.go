package site24x7

import (
	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
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
			// @TODO(mohmann): requires_authentication is no valid field in the
			// ITAutomations API anymore and is thus ignored. We should remove
			// it completely from the resource in the future. This is just here
			// for backwards compatibility.
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

func actionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	automation := resourceDataToAction(d)

	automation, err := client.ITAutomations().Create(automation)
	if err != nil {
		return err
	}

	d.SetId(automation.ActionID)

	return nil
}

func actionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	automation, err := client.ITAutomations().Get(d.Id())
	if err != nil {
		return err
	}

	updateActionResourceData(d, automation)

	return nil
}

func actionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	automation := resourceDataToAction(d)

	automation, err := client.ITAutomations().Update(automation)
	if err != nil {
		return err
	}

	d.SetId(automation.ActionID)

	return nil
}

func actionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	return client.ITAutomations().Delete(d.Id())
}

func actionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(site24x7.Client)

	_, err := client.ITAutomations().Get(d.Id())
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDataToAction(d *schema.ResourceData) *api.ITAutomation {
	return &api.ITAutomation{
		ActionID:               d.Id(),
		ActionMethod:           d.Get("method").(string),
		ActionName:             d.Get("name").(string),
		ActionTimeout:          d.Get("timeout").(int),
		ActionType:             d.Get("type").(int),
		ActionUrl:              d.Get("url").(string),
		CustomParameters:       d.Get("custom_parameters").(string),
		SendCustomParameters:   d.Get("send_custom_parameters").(bool),
		SendInJsonFormat:       d.Get("send_in_json_format").(bool),
		SendIncidentParameters: d.Get("send_incident_parameters").(bool),
	}
}

func updateActionResourceData(d *schema.ResourceData, automation *api.ITAutomation) {
	d.Set("method", automation.ActionMethod)
	d.Set("name", automation.ActionName)
	d.Set("timeout", automation.ActionTimeout)
	d.Set("type", automation.ActionType)
	d.Set("url", automation.ActionUrl)
	d.Set("custom_parameters", automation.CustomParameters)
	d.Set("send_custom_parameters", automation.SendCustomParameters)
	d.Set("send_in_json_format", automation.SendInJsonFormat)
	d.Set("send_incident_parameters", automation.SendIncidentParameters)
}
