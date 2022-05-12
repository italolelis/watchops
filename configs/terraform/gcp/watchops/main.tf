module "foundation" {
  source                      = "../foundation"
  project_id                  = var.project_id
}

module "bigquery" {
  source                         = "../bigquery"
  project_id                     = var.project_id
  bigquery_region                = var.region
  watchops_service_account_email = module.foundation.watchops_service_account_email
  depends_on = [
    module.foundation
  ]
}

module "github" {
  source = "../github"
  count = contains(var.parsers, "github") ? 1 : 0
  project_id  = var.project_id
  region  = var.region
  watchops_service_account_email = module.foundation.watchops_service_account_email
}

module "gitlab" {
  source = "../gitlab"
  count = contains(var.parsers, "gitlab") ? 1 : 0
  project_id  = var.project_id
  region  = var.region
  watchops_service_account_email = module.foundation.watchops_service_account_email
}

module "opsgenie" {
  source = "../opsgenie"
  count = contains(var.parsers, "opsgenie") ? 1 : 0
  project_id  = var.project_id
  region  = var.region
  watchops_service_account_email = module.foundation.watchops_service_account_email
}

module "circleci" {
  source = "../circleci"
  count = contains(var.parsers, "circleci") ? 1 : 0
  project_id  = var.project_id
  region  = var.region
  watchops_service_account_email = module.foundation.watchops_service_account_email
}
