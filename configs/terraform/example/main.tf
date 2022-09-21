terraform {
  required_version = ">= 0.15"
  required_providers {
    google = {
      version = "~> 3.86.0"
    }
  }
}

module "watchops" {
  source                      = "../gcp/watchops"
  project_id                  = "urbansportsclub-dev"
  region                      = "europe-west3"
  bigquery_region             = "EU"
  parsers                     = ["github"]
}
