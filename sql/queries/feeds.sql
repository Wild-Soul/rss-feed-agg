-- name: CreatFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetNextFeedsToFetch :many
select * from feeds
where last_fetched_at < $1
order by last_fetched_at asc
limit $2;

-- name: MarkFeedAsFetched :one
update feeds
set last_fetched_at = now(),
updated_at = now()
where id = $1
returning *;

-- name: GetTotalFeeds :one
select count(id) from feeds;
