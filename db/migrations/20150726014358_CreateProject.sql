
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "projects" (
  "id" integer PRIMARY KEY AUTOINCREMENT NOT NULL,
  "uuid" varchar(255) NOT NULL,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime
);

CREATE INDEX "index_projects_on_uuid" ON "projects" ("uuid");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "projects";