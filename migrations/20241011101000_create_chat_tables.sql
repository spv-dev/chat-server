-- +goose Up
-- +goose StatementBegin
create table chats (
    id bigserial primary key,
    title text not null,
    state smallint not null default 1,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp
);

create table messages (
    id bigserial primary key,
    chat_id bigint not null,
    user_id bigint not null,
    body text not null,
    type_id smallint not null default 10,
    state smallint not null default 1,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp
);

create table chat_users (
    id bigserial primary key,
    chat_id bigint not null,
    user_id bigint not null,
    state smallint not null default 1,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chats;

drop table messages;

drop table chat_users;
-- +goose StatementEnd