CREATE TABLE users
(
    id serial primary key,
    username varchar(255) not null unique,
    password_hash varchar(255) not null unique
);

CREATE TABLE lists
(
    id serial primary key,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE users_lists
(
    user_id int not null,
    list_id int not null,
    primary key (user_id, list_id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (list_id) references lists(id) on delete cascade
);

CREATE TABLE books
(
    id serial primary key,
    title varchar(255) not null,
    author varchar(100) not null,
    publisher varchar(100),
    publication_year int,
    page_count int
);

CREATE TABLE lists_books
(
    list_id int not null,
    book_id int not null,
    primary key (list_id, book_id),
    foreign key (list_id) references lists(id) on delete cascade,
    foreign key (book_id) references books(id) on delete cascade
);