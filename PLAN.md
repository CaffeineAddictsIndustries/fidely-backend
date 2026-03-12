# Fidely Backend — Project Plan

## Current Status

### Done

- Go module initialized (`fidely-backend`).
- API server bootstrap implemented with Echo (`cmd/api/main.go`).
- Health endpoint implemented at `GET /health`.
- Environment/config loader implemented (`DATABASE_URL`, `SERVER_PORT`).
- PostgreSQL connection pool implemented with `pgx`.
- Migrations created and validated locally with PostgreSQL:
  - 14 up / 14 down migration files.
  - 13 application tables created.
  - 14 foreign keys created.
  - Migration version reached `14` successfully.

### In Scope Now

- Admin-only authentication (`store_admin` and `fidely_admin`).
- Secure login/session flow for the existing login page.
- No customer login in this phase.
- Landing page after login will be implemented in a later step.

---

## Admin Auth Plan

1. Add auth migrations:
    - Convert admin password columns to `password_hash` semantics.
    - Add unique index/constraint for `fidely_admin.username`.
    - Create `admin_sessions` table with:
      - `id`, `admin_type`, `admin_id`, `session_token_hash`, `expires_at`, `revoked_at`, `created_at`, `last_seen_at`.
    - Add indexes for token lookup, admin session lookup, and expiry cleanup.

2. Add auth configuration:
    - Session cookie name.
    - Session TTL.
    - Secure cookie toggle.
    - SameSite policy.
    - Auth secret/pepper if required by implementation.

3. Implement password module:
    - Hash and verify password (Argon2id preferred).
    - Never store or return plaintext passwords.

4. Implement session module:
    - Generate cryptographically random opaque session token.
    - Store only token hash in DB.
    - Validate expiration and revocation.
    - Revoke on logout.

5. Implement admin auth endpoints:
    - `POST /admin/auth/login`
    - `POST /admin/auth/logout`
    - `GET /admin/auth/me`

6. Add middleware and authorization:
    - Session validation middleware.
    - Role/context injection (`store_admin` vs `fidely_admin`).
    - Route guards by admin type.

7. Add tests:
    - Password hashing/verification.
    - Session lifecycle.
    - Login/logout/me flow.
    - Authorization boundaries.

---

## Login Response Messaging Policy

For the login page integration in this phase:

- On successful login, return a general success message.
- On failed login, return a general failure message.
- Also return a helper reason message for the UI when applicable:
  - `username does not exist`
  - `incorrect password`

Suggested response contract:

- `success`: `true|false`
- `message`: general message (`login successful` or `login failed`)
- `reason`: optional helper reason string for failed logins

Notes:

- Keep HTTP status usage consistent (`200` for success, `401` for auth failure).
- Do not include sensitive details (hashes, internal IDs, stack traces).
- Landing page redirection logic will be added in a later step.

---

## Notes

- Database name is standardized as `fidely`.
- Runtime requires `DATABASE_URL`.
- Current backend is bootstrap-ready and migration-ready; admin auth is the next implementation block.
