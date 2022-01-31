terraform {
  backend "s3" {
    region = "eu-central-1"
    bucket = "terraform-eks"
    key    = "fourkeys/terraform.tfstate"
  }

  required_providers {
    aws = "~> 3.3"
  }
}


provider "aws" {
  region = var.region

  assume_role {
    role_arn     = var.assume_role_arn
    session_name = "terraform-watchops"
  }
}
