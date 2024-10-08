---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "grafana_machine_learning_outlier_detector Resource - terraform-provider-grafana"
subcategory: "Machine Learning"
description: |-
  An outlier detector monitors the results of a query and reports when its values are outside normal bands.
  The normal band is configured by choice of algorithm, its sensitivity and other configuration.
  Visit https://grafana.com/docs/grafana-cloud/machine-learning/outlier-detection/ for more details.
---

# grafana_machine_learning_outlier_detector (Resource)

An outlier detector monitors the results of a query and reports when its values are outside normal bands.

The normal band is configured by choice of algorithm, its sensitivity and other configuration.

Visit https://grafana.com/docs/grafana-cloud/machine-learning/outlier-detection/ for more details.

## Example Usage

### DBSCAN Outlier Detector

This outlier detector uses the DBSCAN algorithm to detect outliers.

```terraform
resource "grafana_machine_learning_outlier_detector" "my_dbscan_outlier_detector" {
  name        = "My DBSCAN outlier detector"
  description = "My DBSCAN Outlier Detector"

  metric          = "tf_test_dbscan_job"
  datasource_type = "prometheus"
  datasource_uid  = "AbCd12345"
  query_params = {
    expr = "grafanacloud_grafana_instance_active_user_count"
  }
  interval = 300

  algorithm {
    name        = "dbscan"
    sensitivity = 0.5
    config {
      epsilon = 1.0
    }
  }
}
```

### MAD Outlier Detector

This outlier detector uses the Median Absolute Deviation (MAD) algorithm to detect outliers.

```terraform
resource "grafana_machine_learning_outlier_detector" "my_mad_outlier_detector" {
  name        = "My MAD outlier detector"
  description = "My MAD Outlier Detector"

  metric          = "tf_test_mad_job"
  datasource_type = "prometheus"
  datasource_uid  = "AbCd12345"
  query_params = {
    expr = "grafanacloud_grafana_instance_active_user_count"
  }
  interval = 300

  algorithm {
    name        = "mad"
    sensitivity = 0.7
  }
}
```

## Import

Import is supported using the following syntax:

```shell
terraform import grafana_machine_learning_outlier_detector.name "{{ id }}"
```
