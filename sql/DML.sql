INSERT INTO mst_customer (name, phone_number, address)
VALUES ('Jilong', '082121212121', 'Riau');
INSERT INTO mst_customer (name, phone_number, address)
VALUES ('yujong', '083131313131', 'Riau');
INSERT INTO mst_customer (name, phone_number, address)
VALUES ('nana', '0841414141414', 'Riau');

INSERT INTO mst_employee (name, phone_number, address)
VALUES ('Montoon', '085151515151', 'Riau');
INSERT INTO mst_employee (name, phone_number, address)
VALUES ('Garena', '08616161616161', 'Riau');
INSERT INTO mst_employee (name, phone_number, address)
VALUES ('Tencent', '087171717171', 'Riau');

INSERT INTO mst_product (name, price, unit)
VALUES ('Cuci bersih', 3000, 'KG');
INSERT INTO mst_product (name, price, unit)
VALUES ('Cuci + Gosok', 5000, 'KG');
INSERT INTO mst_product (name, price, unit)
VALUES ('Cuci Karpet', 15000, 'Pcs');

INSERT INTO transactions (bill_date, entry_date, finish_date, employe_id, customer_id)
VALUES ('2024-04-05', '2024-04-05', '2024-04-07', 1, 2);
INSERT INTO transactions (bill_date, entry_date, finish_date, employe_id, customer_id)
VALUES ('2024-04-07', '2024-04-07', '2024-04-10', 1, 3);

INSERT INTO trx_detail (trx_id, product_id, product_price, qty)
VALUES (1, 1, 3000, 3);
INSERT INTO trx_detail (trx_id, product_id, product_price, qty)
VALUES (1, 2, 5000, 2);
INSERT INTO trx_detail (trx_id, product_id, product_price, qty)
VALUES (2, 3, 15000, 2);