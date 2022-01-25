-- migrate:up
create table if not exists public.user (
    id SERIAL primary key,
    username VARCHAR(240) unique not null
);
create table if not exists public.group (
    id SERIAL primary key,
    name VARCHAR(240)
);
create table if not exists user_group (
    id SERIAL primary key,
    user_id int references public.user(id) not null,
    group_id int references public.group(id) not null,
    UNIQUE (user_id, group_id)
);
create table if not exists message (
    id SERIAL primary key,
    re_id int references message(id),
    sender int references public.user(id) not null,
    recipient int references public.user(id) not null,
    subject text,
    body text,
    sent_at timestamp without time zone not null
)
-- migrate:down