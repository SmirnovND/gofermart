CREATE SEQUENCE balance_id_seq;
CREATE SEQUENCE transaction_id_seq;
ALTER TABLE transaction ALTER COLUMN id SET DEFAULT nextval('transaction_id_seq');
ALTER TABLE balance ALTER COLUMN id SET DEFAULT nextval('balance_id_seq');
