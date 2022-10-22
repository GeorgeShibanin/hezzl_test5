create table campaigns (
    id serial primary key,
    name text
);

create table items  (
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

create index on campaigns (id);
create index on items (id, campaign_id, name);

CREATE SEQUENCE books_sequence
  start 1
  increment 1;
