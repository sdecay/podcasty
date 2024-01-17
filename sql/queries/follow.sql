-- name: CreateFollow :one
insert into follows (id, created_at, updated_at, user_id, feed_id)
values($1, $2, $3, $4, $5)
returning *;

-- name: GetFollowed :many
select * from follows where user_id=$1;

-- name: DeleteFollowed :exec
delete from follows where id=$1 and user_id=$2;
