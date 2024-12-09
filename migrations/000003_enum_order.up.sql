CREATE TYPE status_enum AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
ALTER TABLE "order" ADD COLUMN status status_enum NOT NULL DEFAULT 'NEW';