resource "google_container_cluster" "primary" {
  name   = "marcellus-wallus"
  region = "us-west1"

  master_auth {
    username = ""
    password = ""
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels {
      managed_by = "Terraform"
    }

    tags = []
  }
}
