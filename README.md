# Voucher Redemption Service

## Prerequisites

- Go version `1.23.4` or higher is installed.
- MySQL database server is installed and running.
- `migrate` tool for database migration management (optional, if you want to run migrations manually).

## How to Run the Project

This project uses `Makefile` to simplify the build, run, and database management processes. Here are the available commands:

**1. Environment Variables Configuration**

The application uses environment variables for configuration. You need to create a `.env` file in the root directory of the project. You can use `.env.example` as a template.

Here is an example of the environment variables you need to configure:

```env
PORT=8000
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_ADDRESS=127.0.0.1:3306
DB_NAME=your_db_name
```

**2. Build Application**

To compile the application, use the command:

```bash
make build
```

**3. Run Application**

To compile the application, use the command:

```bash
make run
```

**4. Run Database Migration**

To compile the application, use the command:

```bash
make migrate-up
```

**5. Run Database Seeder**

To compile the application, use the command:

```bash
make seed
```

## API Endpoints and Example Payloads

Below is a list of available API endpoints along with example payloads to interact with the application. The base URL for all endpoints is http://localhost:<PORT>/api/v1.

**Brand Endpoints**

- POST /api/v1/brand - Create New Brand
  Payload (Example):
  `  {
      "name": "Brand ABC",
      "description": "Description for Brand ABC"
  }`

**Voucher Endpoints**

- POST /api/v1/voucher - Create New Brand
  Payload (Example):
  `  {
      "name": "10% Discount Voucher",
      "brand_id": 1,
      "cost_in_points": 100,
      "description": "10% discount voucher for all products"
  }`
- GET /api/v1/voucher/brand?id={brand_id} - Get Voucher List by Brand ID
  Replace {brand_id} with the desired Brand ID.

- GET /api/v1/voucher?id={voucher_id} - Get Voucher Details by Voucher ID
  Replace {voucher_id} with the desired Voucher ID.

**Transaction (Redemption) Endpoints**

- POST /api/v1/transaction/redemption - Create Voucher Redemption Transaction
  Payload (Example):

  ```
  {
  "customer_id": 123,
  "vouchers": [
  {
  "voucher_id": 1,
  "quantity": 2
  },
  {
  "voucher_id": 2,
  "quantity": 1
  }
  ]
  }

  ```

  ```

  ```

- GET /api/v1/transaction/redemption?transactionId={transactionId} - Get Redemption Transaction Details by Transaction ID
  Replace {transactionId} with the desired Redemption Transaction ID.

**_Make sure to run go mod tidy and go mod vendor to manage project dependencies if needed._**
