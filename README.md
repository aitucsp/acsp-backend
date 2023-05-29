# Development of a corporate self-study portal

# This project which contains 4 sections:
- Project-Based Learning (projects, courses)
- Search for teammates
- Materials and Articles
- Contests

## Technological stack:
- Go
- PostgreSQL, Redis 
- Docker
- Prometheus and Grafana (in the nearest future)

## Some of used libraries and packages:
- Fiber (web router)
- Viper (for application configs)
- Swag (generating OpenAPI Documentation)
- Bcrypt (encrypting, generating hashes)
- Sqlx (database)
- Zap (well-designed and structured logging)
- golang-jwt (for auth via jwt)
- go-playground/validator (validating struct fields)
- mock (mocking different layers)
- AWS SDK (s3 bucket)

## How to configure
1. Set up your own `.env` file (`example.env` will be given for example)
2. Set up your Postgres and Redis databases (migrations given in `./migrations` folder)

## How to launch
1. Set up your ports in `base.env`, `Dockerfile`, `docker-compose.yml`
2. Build a container by docker-compose using `docker-compose up` command in the root directory
3. Run `localhost:<your-port>/swagger` to see the documentation of application (Swagger UI)

