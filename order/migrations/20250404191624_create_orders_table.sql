-- +goose Up
create table orders (
    id serial primary key,
    order_uuid uuid not null unique default uuid_generate_v4(),
    user_uuid uuid not null,
    part_uuids uuid[] not null,
    total_price double precision not null,
    transaction_uuid uuid,
    payment_method text,
    order_status text not null
);

-- +goose Down
drop table if exists orders;