locals {
  redis_plan = {
    development = "rediscloud"
    staging     = "rediscloud:100"
    production  = "rediscloud:500"
  }
}
