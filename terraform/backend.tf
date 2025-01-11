// S3 стейт
module "s3_bucket_state" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = join("-", [var.namespace, "terraform-state"])

  versioning = {
    enabled = true
  }
  tags = local.defaultTags
}

// DynamoDB лок-таблица
module "dynamo_db_state_lock" {
  source   = "terraform-aws-modules/dynamodb-table/aws"
  name = join("-", [var.namespace, "terraform-locks"])
  hash_key = "LockID"
  attributes = [
    {
      name = "LockID"
      type = "S"
    }
  ]
  tags = local.defaultTags
}