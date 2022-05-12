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
}

resource "google_project_iam_member" "pubsub_service_account_token_creator" {
  project = var.project_id
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  role    = "roles/iam.serviceAccountTokenCreator"
}
