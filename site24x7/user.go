package site24x7

import (
	site24x7 "github.com/Bonial-International-GmbH/site24x7-go"
	"github.com/Bonial-International-GmbH/site24x7-go/api"
	apierrors "github.com/Bonial-International-GmbH/site24x7-go/api/errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var UserSchema = map[string]*schema.Schema{
	"display_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"email_address": {
		Type:     schema.TypeString,
		Required: true,
	},
	"role": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"notify_medium": {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
	},
	"selection_type": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"user_groups": {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"status_iq_role": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"cloudspend_role": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"mobile_settings": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mobile_number": {
					Type:     schema.TypeString,
					Required: true,
				},
				"country_code": {
					Type:     schema.TypeString,
					Required: true,
				},
				"voice_call": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sms": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			},
		},
	},
	"alert_settings": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"email_format": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"dont_alert_on_days": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeInt},
				},
				"alert_start_time": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"alert_end_time": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"down": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"trouble": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"up": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"applogs": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			},
		},
	},
}

func resourceSite24x7User() *schema.Resource {
	return &schema.Resource{
		Create: userCreate,
		Read:   userRead,
		Update: userUpdate,
		Delete: userDelete,
		Exists: userExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: UserSchema,
	}
}

func userCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	user := resourceDataToUser(d)

	user, err := client.Users().Create(user)
	if err != nil {
		return err
	}

	d.SetId(user.UserID)

	return nil
}

func userRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	user, err := client.Users().Get(d.Id())
	if err != nil {
		return err
	}

	updateUserResourceData(d, user)

	return nil
}

func userUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	user := resourceDataToUser(d)

	user, err := client.Users().Update(user)
	if err != nil {
		return err
	}

	d.SetId(user.UserID)

	return nil
}

func userDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(site24x7.Client)

	err := client.Users().Delete(d.Id())
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func userExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(site24x7.Client)

	_, err := client.Users().Get(d.Id())
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDataToUser(d *schema.ResourceData) *api.User {
	user := &api.User{
		UserID:         d.Id(),
		DisplayName:    d.Get("display_name").(string),
		EmailAddress:   d.Get("email_address").(string),
		Role:           d.Get("role").(int),
		SelectionType:  d.Get("selection_type").(int),
		StatusIQRole:   d.Get("status_iq_role").(int),
		CloudspendRole: d.Get("cloudspend_role").(int),
	}

	if v, ok := d.GetOk("notify_medium"); ok {
		raw := v.([]interface{})
		mediums := make([]int, len(raw))
		for i, m := range raw {
			mediums[i] = m.(int)
		}
		user.NotifyMedium = mediums
	}

	if v, ok := d.GetOk("user_groups"); ok {
		raw := v.([]interface{})
		groups := make([]string, len(raw))
		for i, g := range raw {
			groups[i] = g.(string)
		}
		user.UserGroups = groups
	}

	if v, ok := d.GetOk("mobile_settings"); ok {
		list := v.([]interface{})
		if len(list) > 0 {
			m := list[0].(map[string]interface{})
			user.MobileSettings = &api.MobileSettings{
				MobileNumber: m["mobile_number"].(string),
				CountryCode:  m["country_code"].(string),
				VoiceCall:    m["voice_call"].(bool),
				SMS:          m["sms"].(bool),
			}
		}
	}

	if v, ok := d.GetOk("alert_settings"); ok {
		list := v.([]interface{})
		if len(list) > 0 {
			a := list[0].(map[string]interface{})

			var dontAlertOnDays []int
			if raw, ok := a["dont_alert_on_days"].([]interface{}); ok {
				dontAlertOnDays = make([]int, len(raw))
				for i, d := range raw {
					dontAlertOnDays[i] = d.(int)
				}
			}

			user.AlertSettings = &api.AlertSettings{
				EmailFormat:     a["email_format"].(int),
				DontAlertOnDays: dontAlertOnDays,
				AlertStartTime:  a["alert_start_time"].(string),
				AlertEndTime:    a["alert_end_time"].(string),
				Down:            a["down"].(bool),
				Trouble:         a["trouble"].(bool),
				Up:              a["up"].(bool),
				AppLogs:         a["applogs"].(bool),
			}
		}
	}

	return user
}

func updateUserResourceData(d *schema.ResourceData, user *api.User) {
	d.Set("display_name", user.DisplayName)       //nolint:errcheck
	d.Set("email_address", user.EmailAddress)     //nolint:errcheck
	d.Set("role", user.Role)                      //nolint:errcheck
	d.Set("notify_medium", user.NotifyMedium)     //nolint:errcheck
	d.Set("selection_type", user.SelectionType)   //nolint:errcheck
	d.Set("user_groups", user.UserGroups)         //nolint:errcheck
	d.Set("status_iq_role", user.StatusIQRole)    //nolint:errcheck
	d.Set("cloudspend_role", user.CloudspendRole) //nolint:errcheck

	if user.MobileSettings != nil {
		d.Set("mobile_settings", []interface{}{ //nolint:errcheck
			map[string]interface{}{
				"mobile_number": user.MobileSettings.MobileNumber,
				"country_code":  user.MobileSettings.CountryCode,
				"voice_call":    user.MobileSettings.VoiceCall,
				"sms":           user.MobileSettings.SMS,
			},
		})
	}

	if user.AlertSettings != nil {
		d.Set("alert_settings", []interface{}{ //nolint:errcheck
			map[string]interface{}{
				"email_format":       user.AlertSettings.EmailFormat,
				"dont_alert_on_days": user.AlertSettings.DontAlertOnDays,
				"alert_start_time":   user.AlertSettings.AlertStartTime,
				"alert_end_time":     user.AlertSettings.AlertEndTime,
				"down":               user.AlertSettings.Down,
				"trouble":            user.AlertSettings.Trouble,
				"up":                 user.AlertSettings.Up,
				"applogs":            user.AlertSettings.AppLogs,
			},
		})
	}
}
