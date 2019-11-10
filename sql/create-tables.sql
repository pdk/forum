-- create-tables.sql

-- init script for sqlite3.
-- use: .read create-tables.sql

create table if not exists users (
    id integer primary key autoincrement,
    joined_at timestamp not null,
    name varchar not null unique
);

create table if not exists topics (
    id integer primary key autoincrement,
    created_by_id int not null references users(id),
    name varchar not null unique
);

create table if not exists threads (
    id integer primary key autoincrement,
    topic_id int not null references topic(id),
    created_by_id int not null references users(id),
    subject varchar not null
);

create table if not exists posts (
    id integer primary key autoincrement,
    thread_id int not null references threads(id),
    posted_by_id int not null references users(id),
    posted_at timestamp not null,
    body varchar not null
);
