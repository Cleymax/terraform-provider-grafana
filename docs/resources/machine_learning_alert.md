---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "grafana_machine_learning_alert Resource - terraform-provider-grafana"
subcategory: "Machine Learning"
description: |-
  
---

# grafana_machine_learning_alert (Resource)



## Example Usage

```terraform
resource "grafana_machine_learning_job" "test_alert_job" {
  name            = "Test Job"
  metric          = "tf_test_alert_job"
  datasource_type = "prometheus"
  datasource_uid  = "abcd12345"
  query_params = {
    expr = "grafanacloud_grafana_instance_active_user_count"
  }
}

resource "grafana_machine_learning_alert" "test_job_alert" {
  job_id            = grafana_machine_learning_job.test_alert_job.id
  title             = "Test Alert"
  anomaly_condition = "any"
  threshold         = ">0.8"
  window            = "15m"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) The title of the alert.

### Optional

- `annotations` (Map of String) Annotations to add to the alert generated in Grafana.
- `anomaly_condition` (String) The condition for when to consider a point as anomalous.
- `for` (String) How long values must be anomalous before firing an alert.
- `job_id` (String) The forecast this alert belongs to.
- `labels` (Map of String) Labels to add to the alert generated in Grafana.
- `no_data_state` (String) How the alert should be processed when no data is returned by the underlying series
- `outlier_id` (String) The forecast this alert belongs to.
- `threshold` (String) The threshold of points over the window that need to be anomalous to alert.
- `window` (String) How much time to average values over

### Read-Only

- `id` (String) The ID of the alert.

## Import

Import is supported using the following syntax:

```shell
terraform import grafana_machine_learning_alert.name "{{ id }}"
```