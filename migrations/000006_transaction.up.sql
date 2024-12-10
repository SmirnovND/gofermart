create table "transaktion"
(
    id      integer        not null
        constraint "transaktion _pk"
            primary key,
    user_id integer        not null
        constraint "transaktion _user_id_fk"
            references "user",
    value   numeric(12, 2) not null
);

create index "transaktion _user_id_index"
    on "transaktion" (user_id);

CREATE TYPE operation_enum AS ENUM ('DEBITING', 'ACCRUAL', 'ERROR');
ALTER TABLE "transaktion" ADD COLUMN operation operation_enum NOT NULL default 'ERROR';