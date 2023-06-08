terraform {
  required_providers {
    byteset = {
      source = "seal-io/byteset"
    }
    alicloud = {
      source = "aliyun/alicloud"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

#####################
# Setup AliCloud RDS
#####################

provider "random" {}

resource "random_string" "rds" {
  length  = 16
  special = false
}

provider "alicloud" {}

locals {
  db_engine                = "PostgreSQL"
  db_engine_version        = "14.0"
  db_instance_type         = "pg.n2.2c.1m"
  db_instance_storage_type = "cloud_essd"
  db_name                  = "byteset"
  db_port                  = 5432
  db_username              = "root"
  db_password              = random_string.rds.result
}

data "alicloud_db_zones" "rds" {
  category                 = "Basic"
  engine                   = local.db_engine
  engine_version           = local.db_engine_version
  db_instance_storage_type = local.db_instance_storage_type
}

locals {
  vpc_cidr_block = "10.0.0.0/16"
  vpc_azs        = try(data.alicloud_db_zones.rds.ids, [])
  vpc_subnets    = [for k, v in local.vpc_azs : cidrsubnet(local.vpc_cidr_block, 8, k)]
}

resource "alicloud_vpc" "rds" {
  vpc_name   = "byteset"
  cidr_block = local.vpc_cidr_block
}

resource "alicloud_vswitch" "rds" {
  count = length(local.vpc_subnets)

  vpc_id     = alicloud_vpc.rds.id
  zone_id    = element(local.vpc_azs, count.index)
  cidr_block = element(local.vpc_subnets, count.index)
}

resource "alicloud_db_instance" "rds" {
  # Must allow publicly accessing.
  security_ips = ["0.0.0.0/0"]
  category     = "Basic"

  instance_name  = "byteset"
  engine         = local.db_engine
  engine_version = local.db_engine_version
  instance_type  = local.db_instance_type

  db_instance_storage_type = local.db_instance_storage_type
  instance_storage         = "20"
  storage_upper_bound      = "100"
  storage_auto_scale       = "Enable"
  storage_threshold        = "20"

  deletion_protection = false

  vpc_id     = alicloud_vpc.rds.id
  vswitch_id = local.vpc_subnets.0
}

resource "alicloud_db_database" "default" {
  name        = local.db_name
  instance_id = alicloud_db_instance.rds.id
}

resource "alicloud_rds_account" "rds" {
  account_type     = "Super"
  db_instance_id   = alicloud_db_instance.rds.id
  account_name     = local.db_username
  account_password = local.db_password
}

resource "alicloud_db_connection" "rds" {
  instance_id = alicloud_db_instance.rds.id
  port        = local.db_port
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
    address = "postgres://${local.db_username}:${local.db_password}@${alicloud_db_connection.rds.connection_string}:${local.db_port}/${local.db_name}"
    salt    = alicloud_db_instance.rds.id
  }
}
