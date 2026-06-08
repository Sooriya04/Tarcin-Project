# Environment Variables

The backend relies on environment variables loaded via a `.env` file at the root of the `backend/` directory.

## Required Variables

### `DB_URL`
The connection string to your PostgreSQL database.
**Format:** `postgres://[user]:[password]@[host]:[port]/[database]?sslmode=disable`
**Example:** `postgres://postgres:localhost@localhost:5432/tarcin?sslmode=disable`

### `PORT`
*(Legacy/Optional)* The HTTP port for the legacy REST router. Defaults to `8080`.
