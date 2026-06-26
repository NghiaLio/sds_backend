# API Documentation - SDS Mobile Backend

This document details the REST API endpoints provided by the backend to manage products, categories, and authentication.

### Base Configuration
*   **Base URL**: `http://localhost:8080`
*   **Headers**: `Content-Type: application/json`

---

## 📬 Common Response Envelope

All API endpoints return responses encapsulated in a standard JSON envelope:

### Success Response
```json
{
  "success": true,
  "data": { ... } // Can be an object, array, or message
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error description details here"
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
*   **Request Body**:
    *   `username` (string, required): Minimum 3 characters. Must be unique.
    *   `password` (string, required): Minimum 6 characters. Will be securely hashed with bcrypt.
*   **Example Request Body**:
    ```json
    {
      "username": "cuongpc10",
      "password": "password123"
    }
    ```
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
*   **Request Body**:
    *   `username` (string, required)
    *   `password` (string, required)
*   **Example Request Body**:
    ```json
    {
      "username": "cuongpc10",
      "password": "123456"
    }
    ```
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI1Mz...",
        "username": "cuongpc10"
      }
    }
    ```

---

### 📝 Part B: Protected Product Categories Endpoints

These endpoints **require** a valid JSON Web Token (JWT) in the request headers:
*   **Header**: `Authorization: Bearer <accessToken>`

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
*   **Request Body**:
    *   `name` (string, required): Cannot be empty.
*   **Example Request Body**:
    ```json
    {
      "name": "Bao bì giấy"
    }
    ```
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
*   **Path Parameters**:
    *   `id` (integer, required)
*   **Example Request Body**:
    ```json
    {
      "name": "Bao bì tái chế"
    }
    ```
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
*   **Path Parameters**:
    *   `id` (integer, required)
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

These endpoints **require** a valid JSON Web Token (JWT) in the request headers:
*   **Header**: `Authorization: Bearer <accessToken>`

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
*   **Request Body**:
    *   `name` (string, required)
    *   `code` (string, required): Must be unique.
    *   `price` (float, required): Cannot be negative.
    *   `stock` (int, required): Cannot be negative.
    *   `category_id` (int, required): Must reference an existing Category.
    *   `description` (string, optional)
    *   `image` (string, optional)
*   **Example Request Body**:
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
*   **Path Parameters**:
    *   `id` (integer, required)
*   **Request Body**: Same schema as Product Creation.
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
*   **Path Parameters**:
    *   `id` (integer, required)
*   **Success Response (200 OK)**:
    ```json
    {
      "success": true,
      "data": {
        "message": "Product deleted successfully"
      }
    }
    ```
