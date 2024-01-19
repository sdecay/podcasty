-- name: CreatePost :one
insert into posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostUniques :many
select url from posts;

-- name: GetUsersPosts :many
select posts.* from posts
join follows on posts.feed_id = follows.feed_id
where follows.user_id=$1
order by posts.published_at desc
limit $2;
