// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
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
returning id, title, url, description, feed_id, published_at, created_at, updated_at
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	FeedID      uuid.UUID
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.FeedID,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.FeedID,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many
select 
  p.id,
  p.title,
  p.url,
  p.description
from posts p
join feed_follows ff on ff.feed_id = p.feed_id
where ff.user_id = $1
order by p.published_at desc
limit $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostsForUserRow struct {
	ID          uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
