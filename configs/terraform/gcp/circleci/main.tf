data "google_project" "project" {
  project_id = var.project_id
}

resource "google_pubsub_topic" "circleci" {
  project = var.project_id
  name    = "watchops-circleci"
}

resource "google_pubsub_topic_iam_member" "service_account_editor" {
  project = var.project_id
  topic   = google_pubsub_topic.circleci.id
  role    = "roles/editor"
  member  = "serviceAccount:${var.watchops_service_account_email}"
}

resource "google_pubsub_subscription" "circleci" {
  project = var.project_id
  name    = "watchops-circleci"
  topic   = google_pubsub_topic.circleci.id
}
