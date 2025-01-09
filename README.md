# Golang API Get Started

## Description
A simple API project built with **Go** and **Gin** for order management. This API can be used to manage order, user and authentication data using JWT.

## How to run the project

### 1. Clone Repository
Clone the repository to your local computer using the following command:

```bash
git clone https://github.com/Ladybert/slash-api-recruitment.git
cd slash-api-recruitment
```

### 2. Dependency Instalation

```bash
go mod tidy

```

### 3. Run the project

```bash
go run cmd/server/main.go

```

# Slash API Documentation

This is the documentation for the Slash API. It provides detailed descriptions of each available endpoint.

---

## Authentication Endpoints (`/auth`)

### Register
**POST** `/auth/register`
- **Description:** Endpoint for user registration.
- **Request Body Example:**

```bash
{
  "username": "exampleUser",
  "password": "examplePass",
  "email": "example@user.com",
}
```

### Login
**POST** `/auth/login`
- **Description:** Endpoint for user login.
- **Request Body Example:**

```bash
{
  "username": "exampleUser",
  "password": "examplePass"
}
```

---

## Product Endpoints (`/products`)

### Get All Products
**GET** `/products`
- **Description:** Retrieves a list of all available products.

### Get Product by ID
**GET** `/products/:id`
- **Description:** Retrieves a specific product by its ID.

---

## Order Endpoints (`/order`)

### Create Order
**POST** `/order`
- **Description:** Creates a new order.
- **Request Body Example:**

```bash
{
  "user_id": "<16 digit of id>"
  "customer_name": "John Doe",
  "customer_phone": "081234567890",
  "customer_address": "123 Example Street",
  "items": [
    {
      "product_id": "123",
      "quantity": 2
    },
    {
      "product_id": "456",
      "quantity": 1
    }
  ]
}
```

### Get Order Details
**GET** `/order/:id`
- **Description:** Retrieves the details of a specific order by its ID.

### Pay Now
**POST** `/order/pay-now`
- **Description:** Completes payment for the current order before it expires.
- **Request Body Example:**

```bash
{
  "order_id": "12345"
}
```

### Update Order
**PUT** `/order/update/:id`
- **Description:** Updates items in an existing order, adjusting stock and total amounts accordingly.
- **Request Body Example:**

```bash
{
  "productID": "789",
  "quantity": 3
}
```

### Delete Order
**DELETE** `/order/order-delete/:id`
- **Description:** Deletes an order by ID, only if it was created by the currently logged-in user.

---

## Notes
- Replace `:id` in the endpoints with the appropriate resource ID.
- Use proper headers (e.g., `Authorization: Bearer <JWT Token>`) where authentication is required.
