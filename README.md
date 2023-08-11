# Aplikasi Enigma Laundry

### Deskripsi

Setelah menyelesaikan pembelajaran tentang Go API, Anda ditugaskan oleh manajemen Enigma Laundry (EL) untuk membuat sebuah aplikasi sederhana berbasis API untuk mencatat transaksi di tokonya.

![logo](./aset/Enigma-Laundry.png)

Fitur-fitur yang diminta oleh manajemen EL adalah:

1.  Struktur/Design Database yang memenuhi kaidah normalisasi berdasarkan nota dibawah ini dengan kriteria sbb :

        - Hasil design dalam bentuk file Script DDL Postgre SQL
        - Design database minimal memiliki 2 tabel master dan 1 tabel transaksi
        - Sediakan sample data dalam bentuk Script DML Postgre SQL

2.  Aplikasi berbasis Console menggunakan bahasa pemrograman Golang dengan kriteria sbb :

        - Aplikasi memiliki menu untuk melakukan VIEW, INSERT, UPDATE, dan DELETE pada tabel master
          1. Manajemen Customer
          2. Manajemen Produk
          3. Manajemen Employee
        - Aplikasi memiliki menu untuk melakukan VIEW dan INSERT pada table Transaksi
          1. Manajemen Transaksi
        - Setiap menu master wajib memiliki minimal 2 jenis validasi yang berbeda
        - Setiap transaksi master wajib memiliki minimal 4 jenis validasi yang berbeda

3.  Dokumentasi cara menjalankan aplikasi dan penggunaan aplikasi dalam bentuk readme.md atau dokumen ektensi word atau pdf

- - -

## API Spec

### Customer API

#### Create Customer

Request :

- Method : `POST`
- Endpoint : `/customers`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "customerName": "string",
  "customerPhoneNumber": "string",
  "customerAddress": "string"
}
```

Response :

- Status : 201 Created
- Body :

```json
{
  "errors": "string",
  "data": {
    "customerId": "string",
    "customerName": "string",
    "customerPhoneNumber": "string",
    "customerAddress": "string"
  }
}
```

#### Get Customer

Request :

- Method : GET
- Endpoint : `/customers/:id`
- Header :
  - Accept : application/json

Response :

- Status : 200 OK
- Body :

```json
{
  "errors": "string",
  "data": {
    "customerId": "string",
    "customerName": "string",
    "customerPhoneNumber": "string",
    "customerAddress": "string"
  }
}
```

#### Update Customer

Request :

- Method : PUT
- Endpoint : `/customers/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "customerName": "string",
  "customerPhoneNumber": "string",
  "customerAddress": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "errors": "string",
  "data": {
    "customerId": "string",
    "customerName": "string",
    "customerPhoneNumber": "string",
    "customerAddress": "string"
  }
}
```

#### Delete Customer

Request :

- Method : DELETE
- Endpoint : `/customers/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "errors": "string",
  "data": "OK"
}
```

### Product API

#### Create Product

Request :

- Method : POST
- Endpoint : `/products`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"productName": "string",
    "productPrice": int,
    "productUnit": "string" (satuan product,cth: Buah atau Kg)
}
```

Response :

- Status Code: 201 Created
- Body:

```json
{
	"error": "string",
	"data": {
		"productId": "string",
		"productName": "string",
		"productPrice": int,
		"productUnit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### List Product

Request :

- Method : GET
- Endpoint : `/products`
  - Header :
  - Accept : application/json
- Query Param :
  - productName : string `optional`,

Response :

- Status Code : 200 OK
- Body:

```json
{
	"error": "string",
	"data": [
		{
			"productId": "string",
			"productName": "string",
			"productPrice": int,
			"productUnit": "string" (satuan product,cth: Buah atau Kg)
		},
		{
			"productId": "string",
			"productName": "string",
			"productPrice": int,
			"productUnit": "string" (satuan product,cth: Buah atau Kg)
		}
	]
}
```

#### Product By Id

Request :

- Method : GET
- Endpoint : `/products/:id`
- Header :
  - Accept : application/json

Response :

- Status Code: 200 OK
- Body :

```json
{
	"error": "string",
	"data": {
		"productId": "string",
		"productName": "string",
		"productPrice": int,
		"productUnit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### Update Product

Request :

- Method : PUT
- Endpoint : `/products/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"productName": "string",
    "productPrice": int,
    "productUnit": "string" (satuan product,cth: Buah atau Kg)
}
```

Response :

- Status Code: 200 OK
- Body :

```json
{
	"error": "string",
	"data": {
		"productId": "string",
		"productName": "string",
		"productPrice": int,
		"productUnit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### Delete Product

Request :

- Method : DELETE
- Endpoint : `/products/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "errors": "string",
  "data": "OK"
}
```

### Transaction API

#### Create Transaction

Request :

- Method : POST
- Endpoint : `/transactions`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"BillDate": "string",
	"EntryDate": "string",
	"FinishDate": "string",
	"EmployeeId": "string",
	"CustomerId": "string",
	"BillDetails": [
		{
			"ProductId": "string",
			"ProductPrice": int,
			"Qty": int
		}
	]
}
```

Request :

- Status Code: 201 Created
- Body :

```json
{
	"error": "string",
	"data":  {
		"billId":  "string",
		"billDate":  "string",
		"entryDate":  "string",
		"finishDate":  "string",
		"employeeId":  "string",
		"customerId":  "string",
		"billDetails":  [
			{
				"billDetailsId":  "string",
				"productId":  "string",
				"productPrice": int,
				"qty": int
			}
		]
	}
}
```

#### Get Transaction

Request :

- Method : GET
- Endpoint : `/transactions/:id_bill`
- Header :
  - Accept : application/json
- Body :

Response :

- Status Code: 200 OK
- Body :

```json
{
	"error": "string",
	"data":  {
		"billId":  "string",
		"billDate":  "string",
		"entryDate":  "string",
		"finishDate":  "string",
		"employeeId":  "string",
		"customerId":  "string",
		"billDetails":  [
			{
				"billDetailsId":  "string",
				"productId":  "string",
				"productPrice": int,
				"qty": int
			}
		]
	}
}
```

#### List Transaction

Pattern string date : `dd-MM-yyyy`

Request :

- Method : GET
- Endpoint : `/transactions`
- Header :
  - Accept : application/json
- Query Param :
  - startDate : string `optional`
  - endDate : string `optional`
  - productName : string `optional`
- Body :

Response :

- Status Code: 200 OK
- Body :

```json
{
	"error": "string",
	"data":  [
		{
			"billId":  "string",
			"billDate":  "string",
			"entryDate":  "string",
			"finishDate":  "string",
			"employeeId":  "string",
			"customerId":  "string",
			"billDetails":  [
				{
					"billDetailsId":  "string",
					"productId":  "string",
					"productPrice": int,
					"qty": int
				}
			]
		}
	]
}
```
