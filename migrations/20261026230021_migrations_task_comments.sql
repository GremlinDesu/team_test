-- +goose Up
create table task_comments (
    id serial primary key,
    task_id bigint unsigned not null references tasks(id) on delete cascade,
    user_id bigint unsigned not null references users(id) on delete cascade,
    content text not null,
    created_at timestamp default current_timestamp
);

-- +goose Down
drop table if exists task_comments;
