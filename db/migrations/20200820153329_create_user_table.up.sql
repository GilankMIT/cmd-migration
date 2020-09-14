create table if not exists users
(
    id          int auto_increment
        primary key,
    first_name  varchar(255) null,
    middle_name varchar(255) null,
    last_name   varchar(255) null,
    email       varchar(255) null,
    username    varchar(255) null,
    password    varchar(255) null,
    created_at  timestamp    null,
    updated_at  timestamp    null,
    constraint users_email_uindex
        unique (email),
    constraint users_username_uindex
        unique (username)
);
