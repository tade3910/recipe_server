# recipe_server

A modern web-based recipe serving platform built with Node.js and Express.js, enhanced with Docker for containerization and PostgreSQL for database management.

## Quick Start Guide

### Docker Setup
1. Install Docker if not already installed
2. Clone the repository
3. Create and run Docker Compose file:
```bash
docker-compose up
```

### Development Configuration

#### Postgres Database
- `DATABASE_URL`: Connection string for your Postgres database (e.g., `postgres://username:password@localhost:5432/RECIPE_DEV`)
- `TEST_DATABASE_DSN`: Connection string for test database (e.g., `host=localhost user=test_user password=test_pass dbname=RECIPE_TEST port=5400 sslmode=disable`)
- Ensure your production database uses proper security settings

#### Frontends
- `FRONTEND_URLS`: List of frontend URLs your application should be accessible to

#### Environment Variables
- Use `.env` file to store sensitive information:
```bash
echo "Postgres DB: ${DATABASE_URL}" > .env
echo "Test DB: ${TEST_DATABASE_DSN}" >> .env
```