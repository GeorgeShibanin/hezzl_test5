-- +goose Up
-- +goose StatementBegin
CREATE TABLE default.items (
                       id int,
                       campaign_id int,
                       name text,
                       description text,
                       priority int,
                       removed boolean,
                       created_at date
                   ) ENGINE = MergeTree()
	                 ORDER BY (id, campaign_id, name)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE default.items
-- +goose StatementEnd
