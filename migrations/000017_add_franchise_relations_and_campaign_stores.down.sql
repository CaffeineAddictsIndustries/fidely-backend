ALTER TABLE "campaign_stores"
    DROP CONSTRAINT IF EXISTS "fk_campaign_stores_store_id";

ALTER TABLE "campaign_stores"
    DROP CONSTRAINT IF EXISTS "fk_campaign_stores_campaign_id";

DROP TABLE IF EXISTS "campaign_stores";

ALTER TABLE "franchise_theme"
    DROP CONSTRAINT IF EXISTS "fk_franchise_theme_franchise_id";

DROP TABLE IF EXISTS "franchise_theme";

ALTER TABLE "card_types"
    DROP CONSTRAINT IF EXISTS "fk_card_types_franchise_id";

ALTER TABLE "card_types"
    DROP COLUMN IF EXISTS "franchise_id";

ALTER TABLE "card_types"
    ADD COLUMN "store_id" INTEGER NOT NULL;

ALTER TABLE "card_types"
    ADD CONSTRAINT "fk_card_types_store_id"
    FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "store_admin"
    DROP CONSTRAINT IF EXISTS "fk_store_admin_franchise_id";

ALTER TABLE "store_admin"
    DROP COLUMN IF EXISTS "franchise_id";

ALTER TABLE "stores"
    DROP CONSTRAINT IF EXISTS "fk_stores_franchise_id";

ALTER TABLE "stores"
    DROP COLUMN IF EXISTS "franchise_id";

DROP TABLE IF EXISTS "franchises";
