DROP INDEX IF EXISTS "idx_admin_sessions_expires_at";
DROP INDEX IF EXISTS "idx_admin_sessions_admin_lookup";
DROP INDEX IF EXISTS "uq_admin_sessions_session_token_hash";

DROP TABLE IF EXISTS "admin_sessions";

DROP INDEX IF EXISTS "uq_fidely_admin_username";

ALTER TABLE "fidely_admin"
    RENAME COLUMN "password_hash" TO "password";

ALTER TABLE "store_admin"
    RENAME COLUMN "password_hash" TO "password";
