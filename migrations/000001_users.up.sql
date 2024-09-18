CREATE TABLE users (
    id           text  unique not null,
    name         text  not null,
    email        text  not null,
    role         int   not null,
    passwordhash bytea not null
);
