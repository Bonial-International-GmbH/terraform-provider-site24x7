package site24x7

import (
	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
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
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func monitorGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	monitorGroup := resourceDataToMonitorGroup(d)

	monitorGroup, err := client.MonitorGroups().Create(monitorGroup)
	if err != nil {
		return err
	}

	d.SetId(monitorGroup.GroupID)

	return nil
}

func monitorGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	monitorGroup, err := client.MonitorGroups().Get(d.Id())
	if err != nil {
		return err
	}

	updateMonitorGroupResourceData(d, monitorGroup)

	return nil
}

func monitorGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	monitorGroup := resourceDataToMonitorGroup(d)

	monitorGroup, err := client.MonitorGroups().Update(monitorGroup)
	if err != nil {
		return err
	}

	d.SetId(monitorGroup.GroupID)

	return nil
}

func monitorGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	return client.MonitorGroups().Delete(d.Id())
}

func monitorGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(site24x7.Client)

	_, err := client.MonitorGroups().Get(d.Id())
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDataToMonitorGroup(d *schema.ResourceData) *api.MonitorGroup {
	return &api.MonitorGroup{
		GroupID:     d.Id(),
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
	}
}

func updateMonitorGroupResourceData(d *schema.ResourceData, monitorGroup *api.MonitorGroup) {
	d.Set("display_name", monitorGroup.DisplayName)
	d.Set("description", monitorGroup.Description)
}
