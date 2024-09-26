-- +goose Up
create table feeds (
  id UUID primary key,
  created_at timestamp not null,
  updated_at timestamp not null default now(),
  name text not null,
  url text unique not null,
  user_id uuid not null references users(id) on delete cascade
);

-- +goose Down
drop table feeds;

