-- +goose Up
-- +goose StatementBegin
CREATE TYPE "account_status" AS ENUM (
  'active',
  'suspended',
  'deleted'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(75) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "phone_number" varchar(15),
  "avatar" varchar,
  "status" account_status NOT NULL DEFAULT 'active',
  "is_staff" bool NOT NULL DEFAULT false,
  "join_date" timestamptz NOT NULL DEFAULT NOW(),
  "last_login" timestamptz
);

CREATE INDEX ON "users" ("status");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS account_status;
-- +goose StatementEnd
