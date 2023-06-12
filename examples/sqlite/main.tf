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
