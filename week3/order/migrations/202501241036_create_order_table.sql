-- +goose Up
create extension if not exists "pgcrypto";

create table orders(
    id serial primary key,
    order_uuid uuid not null unique default gen_random_uuid(),
    user_uuid uuid not null,
    part_uuids uuid[] not null,
    total_price decimal(10, 2) not null,
    transaction_uuid uuid unique default gen_random_uuid(),
    payment_method integer,
    status integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp default now(),
    deleted_at timestamp
);

create index idx_orders_user_uuid on orders(user_uuid);
create index idx_orders_status on orders(status);
create index idx_orders_created_at on orders(created_at);

-- +goose Down
drop table if exists orders;
