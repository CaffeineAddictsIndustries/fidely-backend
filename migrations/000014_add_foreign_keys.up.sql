-- store_admin -> stores
ALTER TABLE "store_admin"
    ADD CONSTRAINT "fk_store_admin_store_id"
    FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

-- card_types -> stores
ALTER TABLE "card_types"
    ADD CONSTRAINT "fk_card_types_store_id"
    FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

-- cards -> users
ALTER TABLE "cards"
    ADD CONSTRAINT "fk_cards_user_id"
    FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- cards -> card_types
ALTER TABLE "cards"
    ADD CONSTRAINT "fk_cards_card_type_id"
    FOREIGN KEY ("card_type_id") REFERENCES "card_types" ("id");

-- campaigns -> card_types
ALTER TABLE "campaigns"
    ADD CONSTRAINT "fk_campaigns_card_type_id"
    FOREIGN KEY ("card_type_id") REFERENCES "card_types" ("id");

-- campaign_redeems -> campaigns
ALTER TABLE "campaign_redeems"
    ADD CONSTRAINT "fk_campaign_redeems_campaign_id"
    FOREIGN KEY ("campaign_id") REFERENCES "campaigns" ("id");

-- campaign_redeems -> cards
ALTER TABLE "campaign_redeems"
    ADD CONSTRAINT "fk_campaign_redeems_card_id"
    FOREIGN KEY ("card_id") REFERENCES "cards" ("id");

-- purchase_types -> card_types
ALTER TABLE "purchase_types"
    ADD CONSTRAINT "fk_purchase_types_card_type_id"
    FOREIGN KEY ("card_type_id") REFERENCES "card_types" ("id");

-- purchases -> purchase_types
ALTER TABLE "purchases"
    ADD CONSTRAINT "fk_purchases_purchase_type_id"
    FOREIGN KEY ("purchase_type_id") REFERENCES "purchase_types" ("id");

-- purchases -> cards
ALTER TABLE "purchases"
    ADD CONSTRAINT "fk_purchases_card_id"
    FOREIGN KEY ("card_id") REFERENCES "cards" ("id");

-- store_theme -> stores
ALTER TABLE "store_theme"
    ADD CONSTRAINT "fk_store_theme_store_id"
    FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

-- push_notifications -> card_types
ALTER TABLE "push_notifications"
    ADD CONSTRAINT "fk_push_notifications_card_type_id"
    FOREIGN KEY ("card_type_id") REFERENCES "card_types" ("id");

-- push_notification_cards -> push_notifications
ALTER TABLE "push_notification_cards"
    ADD CONSTRAINT "fk_push_notification_cards_push_notification_id"
    FOREIGN KEY ("push_notification_id") REFERENCES "push_notifications" ("id");

-- push_notification_cards -> cards
ALTER TABLE "push_notification_cards"
    ADD CONSTRAINT "fk_push_notification_cards_card_id"
    FOREIGN KEY ("card_id") REFERENCES "cards" ("id");
