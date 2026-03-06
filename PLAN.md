# Fidely Backend — Project Plan

## Current Status

### Done

- Go module initialized (`fidely-backend`).
- Basic API server implemented with Echo (`cmd/api/main.go`).
- Health endpoint available at `GET /health`.
- Environment config implemented (`DATABASE_URL` required, `SERVER_PORT` default `8080`).
- PostgreSQL connection pool implemented with `pgx`.
- Migration structure created with `golang-migrate`.
- 14 migration pairs created (`14 up` + `14 down`).
- Schema from `Fidely.sql` mapped into migrations:
    - 13 tables created.
    - 14 foreign keys applied in dedicated FK migration.

### Missing / Pending

- Run migrations in a live PostgreSQL instance and confirm successful `up`.
- Validate resulting DB schema against `Fidely.sql` in a real database.
- Add first domain routes/handlers (stores, users, cards, etc.).
- Implement repository queries and service logic.
- Add auth flow (store admin / fidely admin) and middleware.
- Add tests for config, DB connection, and HTTP handlers.

---

## Migration Inventory

1. `000001_create_stores_table`
2. `000002_create_users_table`
3. `000003_create_store_admin_table`
4. `000004_create_card_types_table`
5. `000005_create_cards_table`
6. `000006_create_campaigns_table`
7. `000007_create_campaign_redeems_table`
8. `000008_create_purchase_types_table`
9. `000009_create_purchases_table`
10. `000010_create_store_theme_table`
11. `000011_create_push_notifications_table`
12. `000012_create_push_notification_cards_table`
13. `000013_create_fidely_admin_table`
14. `000014_add_foreign_keys`

---

## Notes

- Database name is standardized as `fidely`.
- Runtime requires `DATABASE_URL` (no implicit fallback).
- Current app scope is bootstrap-only (health route + infra wiring).
