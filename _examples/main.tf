// Authentication API doc: https://www.site24x7.com/help/api/#authentication
provider "site24x7" {
  // The client ID will be looked up in the SITE24X7_OAUTH2_CLIENT_ID
  // environment variable if the attribute is empty or omitted.
  oauth2_client_id = "${var.oauth2_client_id}"

  // The client secret will be looked up in the SITE24X7_OAUTH2_CLIENT_SECRET
  // environment variable if the attribute is empty or omitted.
  oauth2_client_secret = "${var.oauth2_client_secret}"

  // The refresh token will be looked up in the SITE24X7_OAUTH2_REFRESH_TOKEN
  // environment variable if the attribute is empty or omitted.
  oauth2_refresh_token = "${var.oauth2_refresh_token}"

  // The minimum time to wait in seconds before retrying failed Site24x7 API requests.
  retry_min_wait = 1

  // The maximum time to wait in seconds before retrying failed Site24x7 API
  // requests. This is the upper limit for the wait duration with exponential
  // backoff.
  retry_max_wait = 30

  // Maximum number of Site24x7 API request retries to perform until giving up.
  max_retries = 4

  // Site24x7 API base URL to use. Used to select the desired data center.
  // See https://www.site24x7.com/help/api/#introduction
  api_base_url = "https://www.site24x7.com/api"

  // Site24x7 token URL to use. Used to select the desired data center.
  // See https://www.site24x7.com/help/api/#authentication
  // NOTE: This needs to be configured to match the API base URL domain.
  token_url = "https://accounts.zoho.com/oauth/v2/token"
}

// IT Automation API doc: https://www.site24x7.com/help/api/#it-automation
resource "site24x7_action" "action" {
  // (Required) Display name for the action.
  name = "mywebhook"

  // (Required) URL to be invoked for action execution.
  url = "https://foo.bar/webhook"

  // (Required) The type of the action. See
  // https://www.site24x7.com/help/api/#it-automation-type-constants for allowed values.
  type = 2

  // (Optional) HTTP Method to access the URL. Default: "P". See
  // https://www.site24x7.com/help/api/#http_methods for allowed values.
  method = "G"

  // (Optional) If send_custom_parameters is set as true. Custom parameters to
  // be passed while accessing the URL.
  custom_parameters = "param=value"

  // (Optional) Configuration to send custom parameters while executing the action.
  send_custom_parameters = true

  // (Optional) Configuration to enable json format for post parameters.
  send_in_json_format = true

  // (Optional) Configuration to send incident parameters while executing the action.
  send_incident_parameters = true

  // (Optional) The amount of time a connection waits to time out. Range 1 - 90. Default: 30.
  timeout = 10
}

// Monitor Group API doc: https://www.site24x7.com/help/api/#monitor-groups
resource "site24x7_monitor_group" "monitor_group" {
  // (Required) Display Name for the Monitor Group.
  display_name = "mygroup"

  // (Optional) Description for the Monitor Group.
  description = "This is the description of the group"
}

// Website Monitor API doc: https://www.site24x7.com/help/api/#website
resource "site24x7_website_monitor" "website_monitor" {
  // (Required) Name for the monitor.
  display_name = "mymonitor"

  // (Required) Website address to monitor.
  website = "https://foo.bar"

  // (Optional) Check interval for monitoring. Default: 1. See
  // https://www.site24x7.com/help/api/#check-interval for all supported
  // values.
  check_frequency = 1

  // (Optional) HTTP Method to be used for accessing the website. Default: "G".
  // See https://www.site24x7.com/help/api/#http_methods for allowed values.
  http_method = "P"

  // (Optional) Authentication user name to access the website.
  auth_user = "theuser"

  // (Optional) Authentication password to access the website.
  auth_pass = "thepasswd"

  // (Optional) Check for the keyword in the website response.
  matching_keyword_value = "foo"

  // (Optional) Alert type to match on. See
  // https://www.site24x7.com/help/api/#alert-type-constants for available
  // values.
  matching_keyword_severity = 2

  // (Optional) Check for non existence of keyword in the website response.
  unmatching_keyword_value = "error"

  // (Optional) Alert type to match on. See
  // https://www.site24x7.com/help/api/#alert-type-constants for available
  // values.
  unmatching_keyword_severity = 2

  // (Optional) Match the regular expression in the website response.
  match_regex_value = ".*imprint.*"

  // (Optional) Alert type to match on. See
  // https://www.site24x7.com/help/api/#alert-type-constants for available
  // values.
  match_regex_severity = 2

  // (Optional) Perform case sensitive keyword search or not. Default: false.
  match_case = true

  // (Optional) User Agent to be used while monitoring the website.
  user_agent = "some user agent string"

  // (Optional) Map of custom HTTP headers to send.
  custom_headers = {
    "Accept" = "application/json"
  }

  // (Optional) Timeout for connecting to website. Range 1 - 45. Default: 10
  timeout = 10

  // (Optional) Location Profile to be associated with the monitor. If omitted,
  // the first profile returned by the /api/location_profiles endpoint
  // (https://www.site24x7.com/help/api/#list-of-all-location-profiles) will be
  // used.
  location_profile_id = "123"

  // (Optional) Notification profile to be associated with the monitor. If
  // omitted, the first profile returned by the /api/notification_profiles
  // endpoint (https://www.site24x7.com/help/api/#list-notification-profiles)
  // will be used.
  notification_profile_id = "123"

  // (Optional) Threshold profile to be associated with the monitor. If
  // omitted, the first profile returned by the /api/threshold_profiles
  // endpoint (https://www.site24x7.com/help/api/#list-threshold-profiles) will
  // be used.
  threshold_profile_id = "123"

  // (Optional) List of monitor group IDs to associate the monitor to.
  monitor_groups = [
    "${site24x7_monitor_group.monitor_group.id}",
  ]

  // (Optional) List if user group IDs to be notified on down. If omitted, the
  // first user group returned by the /api/user_groups endpoint
  // (https://www.site24x7.com/help/api/#list-of-all-user-groups) will be used.
  user_group_ids = [
    "123",
  ]

  // (Optional) Map of status to actions that should be performed on monitor
  // status changes. See
  // https://www.site24x7.com/help/api/#action-rule-constants for all available
  // status values.
  actions = {
    "1" = "${site24x7_action.action.id}"
  }

  // (Optional) Resolve the IP address using Domain Name Server. Default: true.
  use_name_server = false

  // (Optional) Provide a comma-separated list of HTTP status codes that indicate a successful response. You can specify individual status codes, as well as ranges separated with a colon. Default: ""
  up_status_codes = "200,404"
}
