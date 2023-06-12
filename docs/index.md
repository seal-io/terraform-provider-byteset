---
layout: ""
page_title: "ByteSet Provider"
description: The ByteSet provider for Terraform is a plugin that seed your database for Development/Testing.
---

# ByteSet Provider

The ByteSet provider for Terraform is a plugin that seed your database for Development/Testing.

## Why this tool?

- Seed your development database with real data.
- Warm up your testing dataset.

## Example Usage

### SQLite

```terraform
terraform {
  required_providers {
    byteset = {
      source = "seal-io/byteset"
    }
  }
}

provider "byteset" {}

resource "byteset_pipeline" "local_file_to_sqlite" {
  # Source to indicate the SQLite SQL file.
  source = {
    address = "file:///path/to/sqlite.sql"
  }

  # Destination to load the SQLite SQL file.
  destination = {
    address       = "sqlite:///path/to/sqlite.db?_pragma=foreign_keys(1)"
    conn_max_open = 1
    conn_max_idle = 1
  }
}
```

### AWS RDS

```terraform
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
```

### AliCloud RDS

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema
