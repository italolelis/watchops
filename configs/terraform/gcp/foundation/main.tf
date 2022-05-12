resource "google_service_account" "watchops" {
  project      = var.project_id
  account_id   = "watchops"
  display_name = "Service Account for WatchOps resources"
}

resource "google_project_iam_member" "bigquery_user" {
  project = var.project_id
  role    = "roles/bigquery.user"
  member  = "serviceAccount:${google_service_account.watchops.email}"
  depends_on = [
    google_service_account.watchops
  ]
}
