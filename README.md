terraform-provider-site24x7
===========================

[![Build Status](https://travis-ci.org/Bonial-International-GmbH/terraform-provider-site24x7.svg?branch=master)](https://travis-ci.org/Bonial-International-GmbH/terraform-provider-site24x7)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bonial-International-GmbH/terraform-provider-site24x7?style=flat)](https://goreportcard.com/report/github.com/Bonial-International-GmbH/terraform-provider-site24x7)
[![GoDoc](https://godoc.org/github.com/Bonial-International-GmbH/terraform-provider-site24x7?status.svg)](https://godoc.org/github.com/Bonial-International-GmbH/terraform-provider-site24x7)

A terraform provider for managing Site24x7 monitors which currently supports
the following resources:

- `site24x7_action` ([Site24x7 IT Automation API doc](https://www.site24x7.com/help/api/#it-automation))
- `site24x7_monitor_group` ([Site24x7 Monitor Group API doc](https://www.site24x7.com/help/api/#monitor-groups))
- `site24x7_website_monitor` ([Site24x7 Monitor API doc](https://www.site24x7.com/help/api/#website))

Installation
------------

Clone the repository and build the provider:

```sh
git clone git@github.com:Bonial-International-GmbH/terraform-provider-site24x7
cd terraform-provider-site24x7
make install
```

This will build the `terraform-provider-site24x7` binary and install it into
the `$HOME/.terraform.d/plugins` directory.

Development
-----------

You can run the tests via:

```sh
make test
```

For a full list of available `make` targets, just run `make` without arguments.

Usage example
-------------

Refer to the [_examples/](_examples/) directory for a fully documented usage example.

This is a quick example of the provider configuration:

```terraform
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
}
```

Usage example with docker
-------------------------

Building
```sh
docker build -t terraform-provider-site24x7 .
```

Using (terraform plan)
```sh
docker run -it --rm -v $PWD:/data/src -e SITE24X7_OAUTH2_CLIENT_ID=$SITE24X7_OAUTH2_CLIENT_ID -e SITE24X7_OAUTH2_CLIENT_SECRET=$SITE24X7_OAUTH2_CLIENT_SECRET -e SITE24X7_OAUTH2_REFRESH_TOKEN=$SITE24X7_OAUTH2_REFRESH_TOKEN -w /data/src terraform-provider-site24x7:latest plan
```
