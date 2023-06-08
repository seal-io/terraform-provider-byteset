resource "byteset_pipeline" "example" {
  source = {
    address = "..."
  }

  destination = {
    address = "..."
    salt    = "..."
  }
}