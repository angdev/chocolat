
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "api_keys" (
  "id" integer PRIMARY KEY AUTOINCREMENT NOT NULL,
  "project_id" integer,
  "value" varchar(255),
  "scope" integer,
  "revoked" boolean DEFAULT 'f',
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "api_keys";