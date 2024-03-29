---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "defectdojo Provider"
subcategory: ""
description: |-
  
---

# defectdojo Provider



## Example Usage

```terraform
provider "defectdojo" {
  host     = "https://defectdojo.example.com"
  username = "admin"
  password = "password"
  token    = "token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `host` (String) The host of the defectdojo instance
- `http_proxy` (String) The HTTP proxy to use for requests to the defectdojo API
- `password` (String, Sensitive) The password of the defectdojo user (required if token is not set)
- `tls_insecure_skip_verify` (Boolean) Whether to insecurely skip verifying the server's certificate chain and host name
- `token` (String, Sensitive) The token of the defectdojo user (required if username and password are not set)
- `username` (String) The username of the defectdojo user (required if token is not set)
