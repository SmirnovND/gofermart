create table "order"
(
    id      serial
        constraint order_pk
            primary key,
    number  varchar(255) not null,
    user_id integer      not null
        constraint order_user_id_fk
            references "user"
);

create unique index order_number_uindex
    on "order" (number);

create index order_user_id_index
    on "order" (user_id);

