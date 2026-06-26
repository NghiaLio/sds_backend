# API Documentation - SDS Mobile Backend

This document details the REST API endpoints provided by the backend to manage products, categories, and authentication.

### Base Configuration
*   **Base URL**: `http://localhost:8080`
*   **Headers**: `Content-Type: application/json`

---

## 🔑 Access Token Expiration
The `accessToken` returned upon successful login is a **JSON Web Token (JWT)**.
*   **Expiration Duration**: **10 minutes** from the time of generation.
*   **Payload Claims**: Contains `user_id`, `username`, and `exp` (unix timestamp of expiration).

---

## 📬 Common Response envelopes

### Success Envelope
All successful requests return a `200 OK` or `201 Created` status with:
```json
{
  "success": true,
  "data": { ... } // Can be an object, array, or message string
}
```

### Error Envelope
All failed requests return appropriate HTTP status codes (400, 401, 404, 500) with:
```json
{
  "success": false,
  "error": "Detailed error message describing the failure"
}
```

---

## ❌ Common Error Codes & Troubleshooting

Below are the typical error responses you will encounter when interacting with the API:

### 1. `400 Bad Request`
Triggered by invalid payloads, validation constraints, or format errors.

*   **JSON parsing error**: Sending invalid JSON structure.
    ```json
    {
      "success": false,
      "error": "Invalid request body: invalid character..."
    }
    ```
*   **Validation error (Register)**: Username < 3 chars or Password < 6 chars.
    ```json
    {
      "success": false,
      "error": "username must be at least 3 characters long"
    }
    ```
*   **Validation error (Category/Product)**: Missing required fields or negative numbers.
    ```json
    {
      "success": false,
      "error": "validation error: product name is required"
    }
    ```
*   **Conflict error**: Creating a product with a `code` that already exists.
    ```json
    {
      "success": false,
      "error": "product code 'SP01' already exists"
    }
    ```
*   **Reference error**: Creating a product with a `category_id` that does not exist.
    ```json
    {
      "success": false,
      "error": "category_id 999 does not exist"
    }
    ```

### 2. `401 Unauthorized`
Triggered by authentication failures.

*   **Missing Header**: Calling protected endpoints without the `Authorization` header.
    ```json
    {
      "success": false,
      "error": "Missing Authorization header"
    }
    ```
*   **Invalid Header Format**: Header is not in `Bearer <token>` format.
    ```json
    {
      "success": false,
      "error": "Invalid Authorization header format. Must be 'Bearer <token>'"
    }
    ```
*   **Invalid/Expired Token**: Token is forged, corrupted, or has expired (exceeded 10 minutes).
    ```json
    {
      "success": false,
      "error": "Invalid or expired access token"
    }
    ```
*   **Invalid Credentials**: Wrong username or password during login.
    ```json
    {
      "success": false,
      "error": "invalid username or password"
    }
    ```

### 3. `404 Not Found`
Triggered when requesting a resource (Category/Product) that does not exist.

*   **Resource Not Found**:
    ```json
    {
      "success": false,
      "error": "product not found"
    }
    ```

---

## 🚦 Endpoints Reference

### 🔐 Part A: Public Endpoints

These endpoints do NOT require an authorization token.

#### A.1. Ping
Check the server status.

*   **URL**: `/ping`
*   **Method**: `GET`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "message": "pong"
      }
    }
    ```

#### A.2. User Registration
Create a new user account.

*   **URL**: `/register`
*   **Method**: `POST`
*   *Body example*: `{"username": "cuongpc10", "password": "password123"}`
*   **Success Response (201 Created)**:
    ```json
    {
      "success": true,
      "data": {
        "message": "User registered successfully"
      }
    }
    ```

#### A.3. User Login
Authenticate user credentials and generate a JWT access token.

*   **URL**: `/login`
*   **Method**: `POST`
*   *Body example*: `{"username": "cuongpc10", "password": "123456"}`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "username": "cuongpc10"
      }
    }
    ```

---

### 📝 Part B: Protected Product Categories Endpoints

These endpoints **require** a valid JWT token in headers: `Authorization: Bearer <accessToken>`.

#### B.1. Get All Categories
Retrieve a list of all product categories.

*   **URL**: `/categories`
*   **Method**: `GET`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": [
        {
          "id": 1,
          "name": "Chai lọ"
        },
        {
          "id": 2,
          "name": "Hộp nhựa"
        }
      ]
    }
    ```

#### B.2. Create Category
Add a new product category.

*   **URL**: `/categories`
*   **Method**: `POST`
*   *Body example*: `{"name": "Bao bì giấy"}`
*   **Success Response (201 Created)**:
    ```json
    {
      "success": true,
      "data": {
        "id": 4,
        "name": "Bao bì giấy"
      }
    }
    ```

#### B.3. Update Category
Modify an existing category.

*   **URL**: `/categories/{id}`
*   **Method**: `PUT`
*   *Body example*: `{"name": "Bao bì tái chế"}`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "id": 4,
        "name": "Bao bì tái chế"
      }
    }
    ```

#### B.4. Delete Category
Delete a category from the database.

*   **URL**: `/categories/{id}`
*   **Method**: `DELETE`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "message": "Category deleted successfully"
      }
    }
    ```

---

### 🛍️ Part C: Protected Products Endpoints

These endpoints **require** a valid JWT token in headers: `Authorization: Bearer <accessToken>`.

#### C.1. Get All Products (Filtered & Paginated)
Retrieve a list of products. Supports pagination, keyword searches, and category filtering.

*   **URL**: `/products`
*   **Method**: `GET`
*   **Query Parameters**:
    *   `page` (optional): Default is `1`.
    *   `limit` (optional): Default is `10`.
    *   `category_id` (optional): Filter products by Category ID.
    *   `keyword` (optional): Search query matched against product name, code, or description.
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": [
        {
          "id": 1,
          "name": "Chai thủy tinh 500ml",
          "code": "SP01",
          "price": 10.0,
          "stock": 100,
          "category_id": 1,
          "description": "Chai thủy tinh cao cấp đựng nước hoa quả",
          "image": "https://example.com/images/chai-tt-500ml.png",
          "created_at": "2026-06-26T12:00:58+07:00",
          "updated_at": "2026-06-26T12:00:58+07:00"
        }
      ]
    }
    ```

#### C.2. Create Product
Add a new product to the catalog.

*   **URL**: `/products`
*   **Method**: `POST`
*   *Body example*:
    ```json
    {
      "name": "Example Product",
      "code": "SP22",
      "price": 12.5,
      "stock": 100,
      "category_id": 1,
      "description": "An example product with optional description.",
      "image": "https://example.com/image.png"
    }
    ```
*   **Success Response (201 Created)**:
    ```json
    {
      "success": true,
      "data": {
        "id": 8,
        "name": "Example Product",
        "code": "SP22",
        "price": 12.5,
        "stock": 100,
        "category_id": 1,
        "description": "An example product with optional description.",
        "image": "https://example.com/image.png",
        "created_at": "2026-06-26T12:00:58+07:00",
        "updated_at": "2026-06-26T12:00:58+07:00"
      }
    }
    ```

#### C.3. Update Product
Modify details of an existing product.

*   **URL**: `/products/{id}`
*   **Method**: `PUT`
*   *Body example*: Same schema as Product Creation.
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "id": 8,
        "name": "Updated Product Name",
        "code": "SP22",
        "price": 15.0,
        "stock": 95,
        "category_id": 1,
        "description": "Updated description text.",
        "image": "https://example.com/image.png",
        "created_at": "2026-06-26T12:00:58+07:00",
        "updated_at": "2026-06-26T12:15:30+07:00"
      }
    }
    ```

#### C.4. Delete Product
Remove a product from the database.

*   **URL**: `/products/{id}`
*   **Method**: `DELETE`
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "message": "Product deleted successfully"
      }
    }
    ```
