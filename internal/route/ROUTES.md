# API Routes Documentation

All endpoints are prefixed with `/api` and require authentication via `AuthMiddleware.Authenticate` unless specified. Authorization is handled by `ProductMiddleware.Authorize` based on user roles (user, admin, super_admin).

## Authentication Routes

- **Base Path**: `/api/auth`
- **Controller**: `AuthController`
  - `POST /signup`: Register a new user.
  - `POST /signin`: Login and get access/refresh tokens.
  - `POST /change-password`: Update user password (authenticated).
  - `POST /refresh-token`: Refresh access token.
  - `POST /change-role`: Change user role (authenticated, admin/super_admin only).
  - `POST /signout`: Revoke refresh token (authenticated).

## Product Routes

- **Base Path**: `/api/products`
- **Controller**: `ProductController`
  - `POST /`: Create a product (admin/super_admin).
  - `GET /:id`: Get product by ID (all roles).
  - `PUT /:id`: Update product (admin/super_admin).
  - `DELETE /:id`: Delete product (super_admin).
  - `GET /`: List products with pagination/filter (all roles).

## Product Category Routes

- **Base Path**: `/api/product-categories`
- **Controller**: `ProductController`
  - `POST /`: Create product category (admin/super_admin).
  - `GET /:id`: Get category by ID (all roles).
  - `PUT /:id`: Update category (admin/super_admin).
  - `DELETE /:id`: Delete category (super_admin).
  - `GET /`: List categories with pagination/filter (all roles).

## Product Stock Routes

- **Base Path**: `/api/product-stocks`
- **Controller**: `ProductController`
  - `POST /`: Create product stock (admin/super_admin).
  - `GET /:id`: Get stock by ID (all roles).
  - `PUT /:id`: Update stock (admin/super_admin).
  - `DELETE /:id`: Delete stock (super_admin).
  - `GET /`: List stocks with pagination/filter (all roles).

## Warehouse Location Routes

- **Base Path**: `/api/warehouse-locations`
- **Controller**: `ProductController`
  - `POST /`: Create warehouse location (admin/super_admin).
  - `GET /:id`: Get location by ID (all roles).
  - `PUT /:id`: Update location (admin/super_admin).
  - `DELETE /:id`: Delete location (super_admin).
  - `GET /`: List locations with pagination/filter (all roles).

## Dashboard Routes

- **Base Path**: `/api/dashboard`
- **Controller**: `ProductController`
  - `GET /`: Get dashboard summary with detailed low-stock, out-of-stock, and recent additions (all roles).
