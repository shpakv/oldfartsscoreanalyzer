module "dynamo_db_old_farts_people" {
  source   = "terraform-aws-modules/dynamodb-table/aws"
  name = join("-", [var.namespace, "old-farts-people"])
  hash_key = "steamId"
  attributes = [
    {
      name = "steamId"
      type = "S"
    }
  ]
  tags = local.defaultTags
}