resource "google_spanner_instance" "global_spanner_api-instance" {
  config           = "us-central1"
  display_name     = "global_spanner_api-instance"
  processing_units = 100
}

resource "google_spanner_database" "global_spanner_api-db" {
  instance = google_spanner_instance.global_spanner_api-instance.name
  name     = "global_spanner_api-db"
}
