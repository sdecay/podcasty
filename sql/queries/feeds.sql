-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: FetchNextFeeds :many
select * from feeds order by fetched_at asc nulls first limit $1;

-- name: MarkFetched :one
update feeds set fetched_at=now(), updated_at=now() where id=$1
returning *;
