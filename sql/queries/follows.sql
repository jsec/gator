-- name: CreateFollow :one
with follow as (
  insert into feed_follows (id, user_id, feed_id, created_at, updated_at)
  values ($1, $2, $3, $4, $5)
  returning *
)
select 
  follow.*,
  f.name as feed_name,
  u.name as user_name
from follow
join feeds f on follow.feed_id = f.id
join users u on follow.user_id = u.id;

-- name: GetFollowsForUser :many
select 
  ff.id,
  f.name
from feed_follows ff
join feeds f on ff.feed_id = f.id
where ff.user_id = $1;

-- name: DeleteFollow :exec
delete from feed_follows
where id in (
  select ff.id 
  from feed_follows ff
  join feeds f on ff.feed_id = f.id
  where ff.user_id = $1
  and f.url = $2
);
