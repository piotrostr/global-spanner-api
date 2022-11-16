resource "google_spanner_instance" "global-spanner-api-instance" {
  config           = "us-central1"
  display_name     = "global-spanner-api-instance"
  processing_units = 100
}

resource "google_spanner_database" "global-spanner-api-db" {
  instance = google_spanner_instance.global-spanner-api-instance.name
  name     = "global-spanner-api-db"
}

resource "google_cloud_run_service" "global-spanner-api" {
  name     = "global-spanner-api"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-central1-docker.pkg.dev/${var.project}/global-spanner-api-repository/global-spanner-api:latest"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
