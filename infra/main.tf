terraform {
  required_version = ">= 1.5"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0"
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
  }
}

provider "neon" {}   # lee NEON_API_KEY
provider "render" {} # lee RENDER_API_KEY y RENDER_OWNER_ID

# --- El proyecto de Neon de PREPROD: no existe hasta que esto se aplica ---

resource "neon_project" "preprod" {
  name = "gastos-preprod"

  # Sin esta línea el primer apply falla: el plan gratuito tiene tope de 6 horas.
  history_retention_seconds = 21600
  region_id                 = "aws-us-east-2"
}

resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id # default_branch_id, NO branch_id
  name       = "app_preprod_owner"
}

resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "app_preprod"
  owner_name = neon_role.preprod.name
}

resource "render_web_service" "api" {
  name   = "gastos-api-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_api # SIN la etiqueta pegada
      tag       = var.image_tag      # la etiqueta va en su propio campo
    }
  }

  # La conexión se DERIVA de lo que devolvió Neon: nadie la copia y pega.
  # No declares COMMIT_SHA: la imagen ya lo trae y una variable del servicio la pisaría.
  env_vars = {
    DB_HOST     = { value = neon_project.preprod.database_host }
    DB_PORT     = { value = "5432" }
    DB_NAME     = { value = neon_database.preprod.name }
    DB_USER     = { value = neon_role.preprod.name }
    DB_PASSWORD = { value = neon_role.preprod.password }
    DB_SSLMODE  = { value = "require" }
    PORT        = { value = "8080" }
  }
}

resource "render_web_service" "front" {
  name   = "gastos-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  env_vars = {
    # La URL sale de la api de Terraform, no de un copy-paste. trimsuffix evita la barra final.
    BACKEND_URL  = { value = trimsuffix(render_web_service.api.url, "/") }
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}