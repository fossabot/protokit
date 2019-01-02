provider "google" {
  project = "${var.gcp_project_id}"
  region  = "${var.gcp_region}"
}

provider "http" {}

data "http" "workstation_external_ip" {
  url = "http://ipv4.icanhazip.com"
}

provider "namecheap" {
  ip = "${data.http.workstation_external_ip.body}"
}

resource "google_project_service" "dns" {
  service = "dns.googleapis.com"
  project = "${var.gcp_project_id}"
}

resource "google_dns_managed_zone" "protokit_io" {
  name       = "protokit-io"
  dns_name   = "protokit.io."
  depends_on = ["google_project_service.dns"]
}

resource "google_dns_record_set" "protokit_io_apex" {
  name         = "${google_dns_managed_zone.protokit_io.dns_name}"
  managed_zone = "${google_dns_managed_zone.protokit_io.name}"
  type         = "A"
  ttl          = 300
  rrdatas      = ["185.199.108.153", "185.199.109.153", "185.199.110.153", "185.199.111.153"] # github pages
}

resource "google_dns_record_set" "protokit_io_www" {
  name         = "www.${google_dns_managed_zone.protokit_io.dns_name}"
  managed_zone = "${google_dns_managed_zone.protokit_io.name}"
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["${google_dns_managed_zone.protokit_io.dns_name}"]
}

resource "namecheap_ns" "protokit_io_nameservers" {
  domain  = "protokit.io"
  servers = ["${google_dns_managed_zone.protokit_io.name_servers}"]
}
