-- +goose Up
alter table feeds add column last_fetched_at timestamp not null default now();

-- +goose Down
alter table feeds drop column last_fetched_at;
