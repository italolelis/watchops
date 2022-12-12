data "google_project" "project" {
  project_id = var.project_id
}

resource "google_pubsub_topic" "circleci" {
  project = data.project.project_id
  name    = "watchops-circleci"
}

resource "google_pubsub_topic_iam_member" "service_account_editor" {
  project = data.project.project_id
  topic   = google_pubsub_topic.circleci.id
  role    = "roles/editor"
  member  = "serviceAccount:${var.watchops_service_account_email}"
}

resource "google_pubsub_subscription" "circleci" {
  project = data.project.project_id
  name    = "watchops-circleci"
  topic   = google_pubsub_topic.circleci.id

  retry_policy {
    maximum_backoff = "600s"
    minimum_backoff = "10s"
  }
}
