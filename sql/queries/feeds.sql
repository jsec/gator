-- name: CreateFeed :one
insert into feeds (id, name, url, user_id, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetAllFeeds :many
select
  f.name,
  f.url,
  u.name as user_name
from feeds f
join users u on f.user_id = u.id;

-- name: GetFeedByURL :one
select * from feeds where url = $1;

-- name: MarkFeedFetched :exec
-- id, timestamp
update feeds
set last_fetched_at = $2,
    updated_at = $2
where id = $1;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at nulls first
limit 1;
