CREATE DATABASE el_api;

CREATE TABLE mst_customer (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	phone_number VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL
)

CREATE TABLE mst_employee (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	phone_number VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL
)

CREATE TABLE mst_product (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	price INT NOT NULL,
	unit VARCHAR(255) NOT NULL
)

CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
	bill_date VARCHAR(255) NOT NULL,
	entry_date VARCHAR(255) NOT NULL,
	finish_date VARCHAR(255) NOT NULL,
	employe_id INT NOT NULL,
	customer_id INT NOT NULL,
	FOREIGN KEY(employe_id) REFERENCES mst_employee(id),
	FOREIGN KEY(customer_id) REFERENCES mst_customer(id)
)

CREATE TABLE trx_detail (
	id SERIAL PRIMARY KEY,
	trx_id INT NOT NULL,
	product_id INT NOT NULL,
	product_price INT NOT NULL,
	qty INT NOT NULL,
	FOREIGN KEY(trx_id) REFERENCES transactions(id),
	FOREIGN KEY(product_id) REFERENCES mst_product(id)
)