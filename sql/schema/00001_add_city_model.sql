-- +goose Up
-- +goose StatementBegin
CREATE TABLE "cities" (
  "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name_en" varchar(75) NOT NULL,
  "name_ar" varchar(75) NOT NULL,
  "is_active" bool NOT NULL DEFAULT true
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "cities";
-- +goose StatementEnd
