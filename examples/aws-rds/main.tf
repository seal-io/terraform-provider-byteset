terraform {
  required_providers {
    byteset = {
      source = "seal-io/byteset"
    }
    aws = {
      source = "hashicorp/aws"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

################
# Setup AWS RDS
################

provider "random" {}

resource "random_string" "rds" {
  length  = 16
  special = false
}

provider "aws" {}

locals {
  db_engine         = "mysql"
  db_engine_version = "5.7"
  db_instance_class = "db.t3.medium"
  db_name           = "byteset"
  db_port           = 3306
  db_username       = "root"
  db_password       = random_string.rds.result
}

data "aws_availability_zones" "rds" {
  state = "available"
}

locals {
  vpc_cidr_block = "10.0.0.0/16"
  vpc_azs        = try(data.aws_availability_zones.rds.names, [])
  vpc_subnets    = [for k, v in local.vpc_azs : cidrsubnet(local.vpc_cidr_block, 8, k)]
}

resource "aws_vpc" "rds" {
  cidr_block           = local.vpc_cidr_block
  instance_tenancy     = "default"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_subnet" "rds" {
  count = length(local.vpc_subnets)

  vpc_id            = aws_vpc.rds.id
  availability_zone = element(local.vpc_azs, count.index)
  cidr_block        = element(local.vpc_subnets, count.index)
}

resource "aws_internet_gateway" "rds" {
  vpc_id = aws_vpc.rds.id
}

resource "aws_route" "rds_internet_gateway" {
  route_table_id         = aws_vpc.rds.default_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.rds.id

  timeouts {
    create = "5m"
  }
}

resource "aws_security_group" "rds" {
  name   = "byteset"
  vpc_id = aws_vpc.rds.id
}

resource "aws_security_group_rule" "rds_ingress" {
  security_group_id = aws_security_group.rds.0.id
  type              = "ingress"
  from_port         = local.db_port
  to_port           = local.db_port
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_db_subnet_group" "rds" {
  name       = "byteset"
  subnet_ids = local.vpc_subnets
}

resource "aws_db_instance" "rds" {
  # Must allow publicly accessing.
  publicly_accessible = true
  multi_az            = false

  identifier     = "byteset"
  engine         = local.db_engine
  engine_version = local.db_engine_version
  instance_class = local.db_instance_class

  db_name  = local.db_name
  port     = local.db_port
  username = local.db_username
  password = local.db_password

  allocated_storage     = 20
  max_allocated_storage = 100

  skip_final_snapshot     = true
  backup_retention_period = 0
  apply_immediately       = true

  db_subnet_group_name   = aws_db_subnet_group.rds.id
  vpc_security_group_ids = [aws_security_group.rds.id]
}

################
# Setup ByteSet
################

provider "byteset" {}

resource "byteset_pipeline" "remote_file_to_mysql" {
  # Source to indicate the MySQL SQL file.
  source = {
    address = "https://somewhere.to.download/mysql.sql"
  }

  # Destination to load the MySQL SQL file.
  destination = {
    address = "mysql://${local.db_username}:${local.db_password}@tcp(${aws_db_instance.rds.address}:${local.db_port})/${local.db_name}"
  }
}