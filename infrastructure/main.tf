data "heroku_pipeline" "pipeline" {
  name = "order-book"
}

# --- Server ---
resource "heroku_app" "server" {
  name   = "ob-server-${var.environment_name}"
  region = "us"

  organization {
    name   = "orbs"
    locked = true
  }

  config_vars = {
    "ENVIRONMENT"           = var.environment_name
    "LOG_LEVEL"             = var.log_level
    "RPC_URL"               = var.rpc_url
    "REACTOR_ADDRESS" =     var.reactor_address
  }

  lifecycle {
    ignore_changes = [
      config_vars
    ]
  }
}

resource "heroku_pipeline_coupling" "server" {
  app_id   = heroku_app.server.id
  pipeline = data.heroku_pipeline.pipeline.id
  stage    = var.environment_name
}

resource "heroku_formation" "server" {
  app_id   = heroku_app.server.id
  type     = "web"
  quantity = var.environment_name == "development" ? 1 : 2
  size     = var.environment_name == "development" ? "Basic" : "Standard-2X"
}

resource "heroku_drain" "server" {
  app_id = heroku_app.server.id
  url    = "https://http-intake.logs.datadoghq.eu/api/v2/logs?dd-api-key=${var.dd_api_key}&ddsource=heroku&env=${var.environment_name}&service=ob-server-${var.environment_name}&host=order-book"
}

# --- Swaps Tracker ---
resource "heroku_app" "swaps-tracker" {
  name   = "ob-swaps-tracker-${var.environment_name}"
  region = "us"

  organization {
    name   = "orbs"
    locked = true
  }

  config_vars = {
    "ENVIRONMENT" = var.environment_name
    "LOG_LEVEL"   = var.log_level
    "RPC_URL"     = var.rpc_url
  }

  lifecycle {
    ignore_changes = [
      config_vars
    ]
  }
}

resource "heroku_pipeline_coupling" "swaps-tracker" {
  app_id   = heroku_app.swaps-tracker.id
  pipeline = data.heroku_pipeline.pipeline.id
  stage    = var.environment_name
}

resource "heroku_formation" "swaps-tracker" {
  app_id   = heroku_app.swaps-tracker.id
  type     = "worker"
  quantity = 1
  size     = var.environment_name == "development" ? "Basic" : "Standard-2X"
}

resource "heroku_drain" "swaps-tracker" {
  app_id = heroku_app.swaps-tracker.id
  url    = "https://http-intake.logs.datadoghq.eu/api/v2/logs?dd-api-key=${var.dd_api_key}&ddsource=heroku&env=${var.environment_name}&service=ob-swaps-tracker-${var.environment_name}&host=order-book"
}

# --- Maker Mock ---
resource "heroku_app" "maker-mock" {
  count  = var.environment_name == "development" ? 1 : 0
  name   = "ob-maker-mock"
  region = "us"

  organization {
    name   = "orbs"
    locked = true
  }

  config_vars = {
    "ENVIRONMENT" = var.environment_name
    "LOG_LEVEL"   = var.log_level
    "API_KEY"     = var.maker_mock_api_key
    "BASE_URL"    = heroku_app.server.web_url
    "IS_DISABLED" = var.maker_mock_is_disabled
    "PRIVATE_KEY" = var.maker_mock_private_key
  }

  lifecycle {
    ignore_changes = [
      config_vars
    ]
  }
}

resource "heroku_pipeline_coupling" "maker-mock" {
  count    = var.environment_name == "development" ? 1 : 0
  app_id   = heroku_app.maker-mock[0].id
  pipeline = data.heroku_pipeline.pipeline.id
  stage    = var.environment_name
}

resource "heroku_formation" "maker-mock" {
  count    = var.environment_name == "development" ? 1 : 0
  app_id   = heroku_app.maker-mock[0].id
  type     = "worker"
  quantity = 1
  size     = var.environment_name == "development" ? "Standard-1X" : "Standard-2X"
}

# --- Redis ---
# https://elements.heroku.com/addons/rediscloud
resource "heroku_addon" "redis" {
  app_id = heroku_app.server.id
  plan   = local.redis_plan[var.environment_name]
}

resource "heroku_addon_attachment" "redis" {
  app_id   = heroku_app.swaps-tracker.id
  addon_id = heroku_addon.redis.id
}
