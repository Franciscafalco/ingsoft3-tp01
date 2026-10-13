variable "image_repo_api" {
  default = "ghcr.io/franciscafalco/ingsoft3-tp01-backend" # SIN la etiqueta
}

variable "image_repo_front" {
  default = "ghcr.io/franciscafalco/ingsoft3-tp01-frontend" # SIN la etiqueta
}

variable "image_tag" {
  default = "sha-6f5c3a831b838b7c5a4fce0d0e1bcd169b48da46" # el commit de v7.0.0
}

variable "region" {
  default = "oregon"
}