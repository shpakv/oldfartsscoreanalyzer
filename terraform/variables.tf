variable "namespace" {
  type        = string
  description = "Namespace, который должен быть использован как префикс для любых сущностей проекта"
  default     = "oldfarts"
}

variable "aws_region" {
  type        = string
  description = "Регион аккаунта, базируемся в Ирландии"
  default     = "eu-west-1"
}