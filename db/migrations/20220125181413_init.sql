-- migrate:up
create table if not exists public.user (
    id SERIAL primary key,
    username VARCHAR(240) unique not null
);
create table if not exists public.group (
    id SERIAL primary key,
    groupname VARCHAR(240)
);
create table if not exists user_group (
    id SERIAL primary key,
    user_id int references public.user(id) not null,
    group_id int references public.group(id) not null,
    CONSTRAINT user_group_unique_user_group UNIQUE (user_id, group_id)
);
create table if not exists message (
    id SERIAL primary key,
    re_id int references message(id),
    sender_id int references public.user(id) not null,
    recipient_id int references public.user(id) null,
    group_id int references public.group(id) null,
    subject text,
    body text,
    sent_at timestamp without time zone not null,
    CHECK (recipient_id is not null OR group_id is not null)
)
-- migrate:down