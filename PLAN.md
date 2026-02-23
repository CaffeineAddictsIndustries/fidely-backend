# Fidely Backend — Project Plan

## Overview

Go backend for the Fidely loyalty card platform, using **golang-migrate** for PostgreSQL database migrations.

- **Go version**: 1.23
- **Database**: PostgreSQL (`fidely`)
- **Migration tool**: [golang-migrate](https://github.com/golang-migrate/migrate)

---

## Project Structure

```
fidely-backend/
├── Fidely.sql                          # Original schema reference
├── PLAN.md
├── go.mod
├── go.sum
├── cmd/
│   └── api/
│       └── main.go                     # Entrypoint (placeholder)
└── migrations/
    ├── 000001_create_stores_table.up.sql
    ├── 000001_create_stores_table.down.sql
    ├── 000002_create_users_table.up.sql
    ├── 000002_create_users_table.down.sql
    ├── 000003_create_store_admin_table.up.sql
    ├── 000003_create_store_admin_table.down.sql
    ├── 000004_create_card_types_table.up.sql
    ├── 000004_create_card_types_table.down.sql
    ├── 000005_create_cards_table.up.sql
    ├── 000005_create_cards_table.down.sql
    ├── 000006_create_campaigns_table.up.sql
    ├── 000006_create_campaigns_table.down.sql
    ├── 000007_create_campaign_redeems_table.up.sql
    ├── 000007_create_campaign_redeems_table.down.sql
    ├── 000008_create_purchase_types_table.up.sql
    ├── 000008_create_purchase_types_table.down.sql
    ├── 000009_create_purchases_table.up.sql
    ├── 000009_create_purchases_table.down.sql
    ├── 000010_create_store_theme_table.up.sql
    ├── 000010_create_store_theme_table.down.sql
    ├── 000011_create_push_notifications_table.up.sql
    ├── 000011_create_push_notifications_table.down.sql
    ├── 000012_create_push_notification_cards_table.up.sql
    ├── 000012_create_push_notification_cards_table.down.sql
    ├── 000013_create_fidely_admin_table.up.sql
    ├── 000013_create_fidely_admin_table.down.sql
    ├── 000014_add_foreign_keys.up.sql
    └── 000014_add_foreign_keys.down.sql
```

---

## Migration Order & Rationale

Tables are ordered so independent tables come first, making FK migration clean:

| # | Migration | Dependencies |
|---|-----------|-------------|
| 1 | `stores` | None |
| 2 | `users` | None |
| 3 | `store_admin` | stores |
| 4 | `card_types` | stores |
| 5 | `cards` | users, card_types |
| 6 | `campaigns` | card_types |
| 7 | `campaign_redeems` | campaigns, cards |
| 8 | `purchase_types` | card_types |
| 9 | `purchases` | purchase_types, cards |
| 10 | `store_theme` | stores |
| 11 | `push_notifications` | card_types |
| 12 | `push_notification_cards` | push_notifications, cards |
| 13 | `fidely_admin` | None (standalone) |
| 14 | **Foreign keys** | All tables above |

---

## Schema Decisions

- **All columns are `NOT NULL`** — the API must explicitly populate every field.
- **No `DEFAULT` values** — nothing is auto-filled.
- **`TIMESTAMPTZ`** used instead of `TIMESTAMP` for timezone awareness.
- **No inline foreign keys** in `CREATE TABLE` — all FKs live in migration 14.
- **Named FK constraints** (e.g., `fk_store_admin_store_id`) so they can be cleanly dropped in the `down` migration.

---

## Migration 14 — Foreign Keys

All foreign keys in a single migration for easy debugging:

| Table | Column | References |
|-------|--------|------------|
| `store_admin` | `store_id` | `stores(id)` |
| `campaigns` | `card_type_id` | `card_types(id)` |
| `campaign_redeems` | `campaign_id` | `campaigns(id)` |
| `campaign_redeems` | `card_id` | `cards(id)` |
| `card_types` | `store_id` | `stores(id)` |
| `cards` | `user_id` | `users(id)` |
| `cards` | `card_type_id` | `card_types(id)` |
| `purchases` | `purchase_type_id` | `purchase_types(id)` |
| `purchases` | `card_id` | `cards(id)` |
| `purchase_types` | `card_type_id` | `card_types(id)` |
| `store_theme` | `store_id` | `stores(id)` |
| `push_notifications` | `card_type_id` | `card_types(id)` |
| `push_notification_cards` | `push_notification_id` | `push_notifications(id)` |
| `push_notification_cards` | `card_id` | `cards(id)` |

The `down` migration drops all FK constraints by name in reverse order.

---

## Running Migrations

Using the `migrate` CLI:

```bash
# Apply all migrations
migrate -path migrations -database "postgresql://user:password@localhost:5432/fidely?sslmode=disable" up

# Roll back last migration
migrate -path migrations -database "postgresql://user:password@localhost:5432/fidely?sslmode=disable" down 1

# Roll back everything
migrate -path migrations -database "postgresql://user:password@localhost:5432/fidely?sslmode=disable" down

# Go to a specific version
migrate -path migrations -database "postgresql://user:password@localhost:5432/fidely?sslmode=disable" goto 13
```

---

## Next Steps (after migrations)

1. Set up project config / env loading
2. Database connection layer
3. Repository pattern per entity
4. HTTP handlers & routing
5. Authentication middleware
6. Business logic / services
