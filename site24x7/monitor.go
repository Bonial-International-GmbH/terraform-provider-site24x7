package site24x7

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*
 * shared type and functions
 * TODO: create common class all monitors inherit from
 */

type Status int

const (
	Down           Status = 0
	Up             Status = 1
	Trouble        Status = 2
	Suspended      Status = 5
	Maintenance    Status = 7
	Discovery      Status = 9
	DiscoveryError Status = 10
)

type MonitorID struct {
	Monitor_id string `json:"monitor_id"`
}

type ValueAndSeverity struct {
	Value    string `json:"value"`
	Severity Status `json:"severity"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ActionRef struct {
	ActionID  string `json:"action_id"`
	AlertType Status `json:"alert_type"`
}

func monitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*http.Client)

	req, err := http.NewRequest(http.MethodDelete, "https://www.site24x7.com/api/monitors/"+d.Id(), nil)
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

func monitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return fetchMonitorExists(meta.(*http.Client), d.Id())
}

func fetchMonitorExists(client *http.Client, id string) (bool, error) {
	var apiResp struct {
		Data []MonitorID `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/monitors", &apiResp); err != nil {
		return false, err
	}
	for _, v := range apiResp.Data {
		if v.Monitor_id == id {
			return true, nil
		}
	}

	return false, nil
}

func defaultLocationProfile(client *http.Client) (string, error) {
	var apiResp struct {
		Data []struct {
			ProfileID string `json:"profile_id"`
		} `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/location_profiles", &apiResp); err != nil {
		return "", err
	}
	return apiResp.Data[0].ProfileID, nil
}

func defaultNotificationProfile(client *http.Client) (string, error) {
	var apiResp struct {
		Data []struct {
			ProfileID string `json:"profile_id"`
		} `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/notification_profiles", &apiResp); err != nil {
		return "", err
	}
	return apiResp.Data[0].ProfileID, nil
}

func defaultThresholdProfile(client *http.Client) (string, error) {
	var apiResp struct {
		Data []struct {
			ProfileID string `json:"profile_id"`
		} `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/threshold_profiles", &apiResp); err != nil {
		return "", err
	}
	return apiResp.Data[0].ProfileID, nil
}

func defaultUserGroup(client *http.Client) (string, error) {
	var apiResp struct {
		Data []struct {
			UserGroupID string `json:"user_group_id"`
		} `json:"data"`
	}
	if err := doGetRequest(client, "https://www.site24x7.com/api/user_groups", &apiResp); err != nil {
		return "", err
	}
	return apiResp.Data[0].UserGroupID, nil
}

func doGetRequest(client *http.Client, url string, data interface{}) error {
	log.Printf("url: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("body: %s", body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("URL: %s failed with: %d. Body: %s", url, resp.StatusCode, body)
	}
	jsonerr := json.Unmarshal(body, &data)

	return jsonerr
}

func parseAPIError(r io.Reader, info_optional ...string) error {
	info := "Not provided"

	if 0 < len(info_optional) {
		info = info_optional[0]
	}

	var apiErr struct {
		ErrorCode int             `json:"error_code"`
		Message   string          `json:"message"`
		ErrorInfo json.RawMessage `json:"error_info"`
	}
	if err := json.NewDecoder(r).Decode(&apiErr); err != nil {
		return fmt.Errorf("json decoding of error failed: (%s). More info: %s", err, info)
	}
	if len(apiErr.ErrorInfo) != 0 {
		return fmt.Errorf("%s (%s). More info: %s", apiErr.Message, string(apiErr.ErrorInfo), info)
	}
	return fmt.Errorf("%s. More info: %s", apiErr.Message, info)
}
