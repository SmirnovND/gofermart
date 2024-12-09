create table "user"
(
    id        serial primary key,
    login     varchar(255) not null,
    pass_hash varchar(255) not null
);

create unique index user_login_uindex
    on "user" (login);

