DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM "store_admin" WHERE "store_id" IS NULL) THEN
        RAISE EXCEPTION 'cannot set store_admin.store_id back to NOT NULL while NULL values exist';
    END IF;
END $$;

ALTER TABLE "store_admin"
    ALTER COLUMN "store_id" SET NOT NULL;
