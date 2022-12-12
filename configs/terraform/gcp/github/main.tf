data "google_project" "project" {
  project_id = var.project_id
}

resource "google_pubsub_topic" "github" {
  project = var.project_id
  name    = "watchops-github"
}

resource "google_pubsub_topic_iam_member" "service_account_editor" {
  project = var.project_id
  topic   = google_pubsub_topic.github.id
  role    = "roles/editor"
  member  = "serviceAccount:${var.watchops_service_account_email}"
}

resource "google_pubsub_subscription" "github" {
  project = var.project_id
  name    = "watchops-github"
  topic   = google_pubsub_topic.github.id

  retry_policy {
    maximum_backoff = "600s"
    minimum_backoff = "10s"
  }
}
