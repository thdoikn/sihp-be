
-- +migrate Up
ALTER TABLE "sihp"."komoditas"
    ADD COLUMN IF NOT EXISTS "gambar_url" TEXT NULL;

-- +migrate Down
ALTER TABLE "sihp"."komoditas"
    DROP COLUMN IF EXISTS "gambar_url";
