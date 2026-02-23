-- Drop foreign keys in reverse order

ALTER TABLE "push_notification_cards"
    DROP CONSTRAINT IF EXISTS "fk_push_notification_cards_card_id";

ALTER TABLE "push_notification_cards"
    DROP CONSTRAINT IF EXISTS "fk_push_notification_cards_push_notification_id";

ALTER TABLE "push_notifications"
    DROP CONSTRAINT IF EXISTS "fk_push_notifications_card_type_id";

ALTER TABLE "store_theme"
    DROP CONSTRAINT IF EXISTS "fk_store_theme_store_id";

ALTER TABLE "purchases"
    DROP CONSTRAINT IF EXISTS "fk_purchases_card_id";

ALTER TABLE "purchases"
    DROP CONSTRAINT IF EXISTS "fk_purchases_purchase_type_id";

ALTER TABLE "purchase_types"
    DROP CONSTRAINT IF EXISTS "fk_purchase_types_card_type_id";

ALTER TABLE "campaign_redeems"
    DROP CONSTRAINT IF EXISTS "fk_campaign_redeems_card_id";

ALTER TABLE "campaign_redeems"
    DROP CONSTRAINT IF EXISTS "fk_campaign_redeems_campaign_id";

ALTER TABLE "campaigns"
    DROP CONSTRAINT IF EXISTS "fk_campaigns_card_type_id";

ALTER TABLE "cards"
    DROP CONSTRAINT IF EXISTS "fk_cards_card_type_id";

ALTER TABLE "cards"
    DROP CONSTRAINT IF EXISTS "fk_cards_user_id";

ALTER TABLE "card_types"
    DROP CONSTRAINT IF EXISTS "fk_card_types_store_id";

ALTER TABLE "store_admin"
    DROP CONSTRAINT IF EXISTS "fk_store_admin_store_id";
