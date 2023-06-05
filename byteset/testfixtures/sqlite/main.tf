terraform {
  required_providers {
    byteset = {
      source = "seal-io/byteset"
    }
  }
}

provider "byteset" {}

resource "byteset_pipeline" "file_to_sqlite" {
  source = {
    address = "file://${path.cwd}/sqlite.sql"
  }

  destination = {
    address = "sqlite://${path.cwd}/sqlite.db?_pragma=foreign_keys(1)"
  }
}
