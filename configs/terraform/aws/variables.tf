variable "assume_role_arn" {
  type        = string
  description = "ARN of the role to be assumed by terraform."
}

variable "region" {
  type        = string
  description = "The AWS region to apply the resources. Defaults to eu-central-1."
  default     = "eu-central-1"
}

variable "kms_deletion_window_in_days" {
  description = "The KMS deletion window in days. Defaults to 10 days."
  type        = number
  default     = 10
}

variable "eks_cluster_name" {
  description = "The name of the EKS cluster."
  type        = string
}


