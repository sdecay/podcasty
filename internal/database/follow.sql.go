// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: follow.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFollow = `-- name: CreateFollow :one
insert into follows (id, created_at, updated_at, user_id, feed_id)
values($1, $2, $3, $4, $5)
returning id, created_at, updated_at, user_id, feed_id
`

type CreateFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFollow(ctx context.Context, arg CreateFollowParams) (Follow, error) {
	row := q.db.QueryRowContext(ctx, createFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i Follow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFollowed = `-- name: DeleteFollowed :exec
delete from follows where id=$1 and user_id=$2
`

type DeleteFollowedParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFollowed(ctx context.Context, arg DeleteFollowedParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollowed, arg.ID, arg.UserID)
	return err
}

const getFollowed = `-- name: GetFollowed :many
select id, created_at, updated_at, user_id, feed_id from follows where user_id=$1
`

func (q *Queries) GetFollowed(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowed, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
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
