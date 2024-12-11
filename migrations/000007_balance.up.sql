create table balance
(
    id      integer
        constraint balance_pk
            primary key,
    user_id integer        not null
        constraint balance_user_id_fk
            references "user",
    value   numeric(12, 2) not null
);

create unique index balance_user_id_uindex
    on balance (user_id);

