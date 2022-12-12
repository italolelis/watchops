data "google_project" "project" {
  project_id = var.project_id
}

resource "google_pubsub_topic" "pagerduty" {
  project                    = var.project_id
  name                       = "watchops-pagerduty"
  message_retention_duration = 7
}

resource "google_pubsub_topic_iam_member" "service_account_editor" {
  project = var.project_id
  topic   = google_pubsub_topic.pagerduty.id
  role    = "roles/editor"
  member  = "serviceAccount:${var.watchops_service_account_email}"
}

resource "google_pubsub_subscription" "pagerduty" {
  project = var.project_id
  name    = "watchops-pagerduty"
  topic   = google_pubsub_topic.pagerduty.id
}
