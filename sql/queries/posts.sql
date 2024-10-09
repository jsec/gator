-- name: CreatePost :one
insert into posts (
  id,
  title,
  url,
  description,
  feed_id,
  published_at,
  created_at,
  updated_at
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostsForUser :many
select 
  p.id,
  p.title,
  p.url,
  p.description
from posts p
join feed_follows ff on ff.feed_id = p.feed_id
where ff.user_id = $1
order by p.published_at desc
limit $2;
