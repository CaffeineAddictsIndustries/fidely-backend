CREATE TABLE "push_notification_cards" (
    "push_notification_id" INTEGER NOT NULL,
    "card_id" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    PRIMARY KEY ("push_notification_id", "card_id")
);
