-- +migrate Up
CREATE TYPE "sihp"."satuan_stok_enum" AS ENUM ('kg', 'gram', 'ons', 'ton', 'liter', 'ml');
CREATE TYPE "sihp"."satuan_periode_enum" AS ENUM ('hari', 'minggu', 'bulan');
CREATE TYPE "sihp"."kelas_komoditas_enum" AS ENUM ('besar', 'menengah', 'kecil');

ALTER TABLE "sihp"."komoditas_dijual"
    ADD COLUMN "harga_normal"               DOUBLE PRECISION              NOT NULL DEFAULT 0,
    ADD COLUMN "harga_mahal"                DOUBLE PRECISION              NOT NULL DEFAULT 0,
    ADD COLUMN "satuan_stok"                "sihp"."satuan_stok_enum"     NOT NULL DEFAULT 'kg',
    ADD COLUMN "nilai_stok"                 NUMERIC(18, 4)                NOT NULL DEFAULT 0,
    ADD COLUMN "satuan_periode"             "sihp"."satuan_periode_enum"  NOT NULL DEFAULT 'minggu',
    ADD COLUMN "nilai_periode"              INTEGER                       NOT NULL DEFAULT 1,
    ADD COLUMN "lokasi_supplier"            VARCHAR(255)                  NOT NULL DEFAULT '',
    ADD COLUMN "pola_distribusi"            VARCHAR(50)                   NULL,
    ADD COLUMN "standardized_stock_periode" NUMERIC(18, 4)                NOT NULL DEFAULT 0,
    ADD COLUMN "kelas_komoditas"            "sihp"."kelas_komoditas_enum" NULL;

ALTER TABLE "sihp"."komoditas_dijual"
    ADD COLUMN "harga_avg" DOUBLE PRECISION GENERATED ALWAYS AS (("harga_normal" + "harga_mahal") / 2.0) STORED;

-- +migrate Down
ALTER TABLE "sihp"."komoditas_dijual"
    DROP COLUMN "harga_avg",
    DROP COLUMN "kelas_komoditas",
    DROP COLUMN "standardized_stock_periode",
    DROP COLUMN "pola_distribusi",
    DROP COLUMN "lokasi_supplier",
    DROP COLUMN "nilai_periode",
    DROP COLUMN "satuan_periode",
    DROP COLUMN "nilai_stok",
    DROP COLUMN "satuan_stok",
    DROP COLUMN "harga_mahal",
    DROP COLUMN "harga_normal";

DROP TYPE "sihp"."kelas_komoditas_enum";
DROP TYPE "sihp"."satuan_periode_enum";
DROP TYPE "sihp"."satuan_stok_enum";