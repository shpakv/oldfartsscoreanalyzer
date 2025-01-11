locals {
  S3GamesFolders = {
    logs        = "logs"
    screenshots = "screenshots"
  }
}

module "s3_bucket_games" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = join("-", [var.namespace, "games"])

  versioning = {
    enabled = true
  }
  tags = local.defaultTags
}

resource "aws_s3_object" "s3_bucket_games_folders" {
  for_each = local.S3GamesFolders
  bucket   = module.s3_bucket_games.s3_bucket_id
  key      = each.value
}