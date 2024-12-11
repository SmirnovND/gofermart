CREATE SEQUENCE balance_id_seq;
CREATE SEQUENCE transaktion_id_seq;
ALTER TABLE transaktion ALTER COLUMN id SET DEFAULT nextval('transaktion_id_seq');
ALTER TABLE balance ALTER COLUMN id SET DEFAULT nextval('balance_id_seq');
