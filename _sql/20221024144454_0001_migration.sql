-- +goose Up
-- +goose StatementBegin
CREATE TABLE campaigns (
                           id serial primary key,
                           name text
);

CREATE TABLE items  (
                        id serial,
                        campaign_id int,
                        name text,
                        description text,
                        priority serial,
                        removed boolean,
                        created_at date,
                        primary key (id, campaign_id),
                        foreign key (campaign_id) references campaigns (id)
);

CREATE INDEX on campaigns (id);
CREATE INDEX on items (id, campaign_id, name);
INSERT INTO campaigns (name) VALUES ('первая запись');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE campaigns;
-- +goose StatementEnd
