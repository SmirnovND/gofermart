create table "transaction"
(
    id      integer        not null
        constraint "transaction _pk"
            primary key,
    user_id integer        not null
        constraint "transaction _user_id_fk"
            references "user",
    value   numeric(12, 2) not null
);

create index "transaction _user_id_index"
    on "transaction" (user_id);

CREATE TYPE operation_enum AS ENUM ('DEBITING', 'ACCRUAL', 'ERROR');
ALTER TABLE "transaction" ADD COLUMN operation operation_enum NOT NULL default 'ERROR';