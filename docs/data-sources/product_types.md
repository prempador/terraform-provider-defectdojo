---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "defectdojo_product_types Data Source - terraform-provider-defectdojo"
subcategory: ""
description: |-
  
---

# defectdojo_product_types (Data Source)



## Example Usage

```terraform
data "defectdojo_product_types" "test" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `product_types` (Attributes List) Product Types (see [below for nested schema](#nestedatt--product_types))

<a id="nestedatt--product_types"></a>
### Nested Schema for `product_types`

Read-Only:

- `authorization_groups` (List of Number) The authorization groups of the product type
- `created` (String) The date the product type was created
- `critical_product` (Boolean) Whether the product type is a critical product
- `description` (String) The description of the product type
- `id` (Number) The unique identifier for the product type
- `key_product` (Boolean) Whether the product type is a key product
- `members` (List of Number) The members of the product type
- `name` (String) The name of the product type
- `updated` (String) The date the product type was last updated
