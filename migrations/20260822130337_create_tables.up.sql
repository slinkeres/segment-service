CREATE TABLE IF NOT EXISTS segments (
    id bigserial primary key,
    name varchar(100) not null unique,
    created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS user_segments(
    user_id bigint not null references users(id) on delete cascade,
    segment_id bigint not null references segments(id) on delete cascade,
    created_at timestamp default now(),
    primary key(user_id, segment_id)
);