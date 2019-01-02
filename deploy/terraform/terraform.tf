terraform {
  backend "gcs" {
    bucket = "protokit-tf-state"
    prefix = "terraform/state"
  }
}

data "terraform_remote_state" "gcs" {
  backend = "gcs"

  config {
    bucket = "protokit-tf-state"
    prefix = "terraform/state"
  }
}
