---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "grafana_oncall_user Data Source - terraform-provider-grafana"
subcategory: ""
description: |-
  HTTP API https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/users/
---

# grafana_oncall_user (Data Source)

* [HTTP API](https://grafana.com/docs/grafana-cloud/oncall/oncall-api-reference/users/)

## Example Usage

```terraform
data "grafana_oncall_user" "alex" {
  username = "alex"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `username` (String) The username of the user.

### Optional

- `id` (String) The ID of this resource.

### Read-Only

- `email` (String) The email of the user.
- `role` (String) The role of the user.

