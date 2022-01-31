module "kinesis_github" {
  source = "modules/kinesis"

  name                      = "watchops-github"
  shard_count               = 2
  retention_period          = 24
  shard_level_metrics       = ["IncomingBytes", "OutgoingBytes"]
  enforce_consumer_deletion = false
  tags = {
    Name      = "watchops_github"
    Workspace = terraform.workspace
    ManagedBy = "Terraform"
  }
}

module "kinesis_opsgenie" {
  source = "modules/kinesis"

  name                      = "watchops-opsgenie"
  shard_count               = 1
  retention_period          = 24
  shard_level_metrics       = ["IncomingBytes", "OutgoingBytes"]
  enforce_consumer_deletion = false
  tags = {
    Name      = "watchops_opsgenie"
    Workspace = terraform.workspace
    ManagedBy = "Terraform"
  }
}

# ============================
# IRSA policy
# This role can be assumed by a service account in Kubernetes to give access to Kinesis.
# ============================
data "aws_eks_cluster" "cluster" {
  name = var.eks_cluster_name
}

module "watchops_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "~> 3.0"

  create_role = true

  role_name = "watchops-role"

  tags = {
    ManagedBy  = "Terraform"
    Enviroment = terraform.workspace
  }

  provider_url = data.aws_eks_cluster.cluster.identity.0.oidc.0.issuer

  role_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
    module.kinesis_stream_iam_policy_read_only_arn,
    module.kinesis_stream_iam_policy_write_only_arn
  ]
  number_of_role_policy_arns = 3

  oidc_fully_qualified_subjects = ["system:serviceaccount:default:watchops"]
}
