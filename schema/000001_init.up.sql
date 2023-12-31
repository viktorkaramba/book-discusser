CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    role          varchar(255) not null
);

CREATE TABLE books
(
    id     serial       not null unique,
    name   varchar(255) not null,
    author varchar(255) not null,
    description varchar not null,
    imagebook varchar not null
);

CREATE TABLE users_books
(
    id      serial                              not null unique,
    user_id int references users (id) on delete cascade not null,
    book_id int references books (id) on delete cascade not null
);

CREATE TABLE comments
(
    id        serial  not null unique,
    message   varchar not null
);


CREATE TABLE books_comments
(
    id         serial                                 not null unique,
    user_id    int references users (id) on delete cascade not null,
    comment_id int references comments (id) on delete cascade not null,
    book_id    int references books (id) on delete cascade not null
);