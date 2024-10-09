-- +goose Up
create table posts (
    id uuid primary key,
    title text not null,
    url text not null unique,
    description text,
    feed_id uuid not null references feeds(id) on delete cascade,
    published_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table posts;

