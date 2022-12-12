variable "project_id" {
  type        = string
  description = "project to deploy four keys resources to"
}

variable "bigquery_region" {
  type        = string
  description = "Region to deploy BigQuery resources in."
}

variable "parsers" {
  type        = list(string)
  description = "List of data parsers to configure. Acceptable values are: 'github', 'gitlab', 'circleci', 'opsgenie'"
}

