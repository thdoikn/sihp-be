
-- +migrate Up
ALTER TABLE "sihp"."pasar"
    ADD COLUMN IF NOT EXISTS "longitude" DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS "latitude" DOUBLE PRECISION NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE "sihp"."pasar"
    DROP COLUMN IF EXISTS "longitude",
    DROP COLUMN IF EXISTS "latitude";
