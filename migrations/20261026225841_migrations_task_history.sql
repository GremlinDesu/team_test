-- +goose Up
CREATE TABLE task_history (
    id serial primary key,
    task_id bigint unsigned not null references tasks(id) on delete cascade,
    changed_by bigint unsigned not null references users(id) on delete cascade,
    changes text not null,
    created_at timestamp default current_timestamp
);

-- +goose Down
drop table if exists task_history;
