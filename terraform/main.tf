provider "aws" {
  region = var.aws_region
}

locals {
  defaultTags = {
    "Manually"  = "false"
    "Terraform" = "true"
    "Group"     = "OldFarts"
    "Project"   = "OldFartsANALyzer"
  }
}

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.83.1"
    }
  }

  backend "s3" {
    bucket         = "oldfarts-terraform-state"
    key            = "terraform/sharedstate.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "oldfarts-terraform-locks"
    encrypt        = true
  }
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_partition" "current" {}