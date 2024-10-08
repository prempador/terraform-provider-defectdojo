---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "defectdojo_users Data Source - terraform-provider-defectdojo"
subcategory: ""
description: |-
  
---

# defectdojo_users (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `users` (Attributes List) List of users (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `configuration_permissions` (List of Number) Configuration permissions of the user
- `date_joined` (String) The date the user joined
- `email` (String) The email of the user
- `first_name` (String) The first name of the user
- `id` (Number) The unique identifier for the user
- `is_active` (Boolean) Whether the user is active
- `is_superuser` (Boolean) Whether the user is a superuser
- `last_login` (String) The last login date of the user
- `last_name` (String) The last name of the user
- `username` (String) The username of the user
