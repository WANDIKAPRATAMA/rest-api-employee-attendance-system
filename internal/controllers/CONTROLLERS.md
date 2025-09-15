# Controllers Documentation

Controllers handle HTTP requests, validate inputs, and coordinate with use cases to process business logic. All responses follow the `RefreshTokenResponse` format.

## ProductController

- **Purpose**: Manages CRUD operations for products, categories, stocks, warehouse locations, and dashboard.
- **Methods**:
  - `CreateProduct`: Creates a new product with validation.
  - `GetProductByID`: Retrieves a product by ID.
  - `UpdateProduct`: Updates a product (admin/super_admin).
  - `DeleteProduct`: Soft deletes a product (super_admin).
  - `GetProductsList`: Lists products with pagination, filter, and search.
  - `CreateProductCategory`: Creates a new category.
  - `GetProductCategoryByID`: Retrieves a category.
  - `UpdateProductCategory`: Updates a category.
  - `DeleteProductCategory`: Deletes a category.
  - `GetProductCategoriesList`: Lists categories.
  - `CreateProductStock`: Creates stock with movement tracking.
  - `GetProductStockByID`: Retrieves stock.
  - `UpdateProductStock`: Updates stock with movement delta.
  - `DeleteProductStock`: Deletes stock.
  - `GetProductStocksList`: Lists stocks.
  - `CreateWarehouseLocation`: Creates a warehouse.
  - `GetWarehouseLocationByID`: Retrieves a warehouse.
  - `UpdateWarehouseLocation`: Updates a warehouse.
  - `DeleteWarehouseLocation`: Deletes a warehouse.
  - `GetWarehouseLocationsList`: Lists warehouses.
  - `GetDashboardSummary`: Provides detailed dashboard data (total stock, low/out-of-stock items, recent additions).

## AuthController

- **Purpose**: Handles user authentication and authorization.
- **Methods**:
  - `Signup`: Registers a new user with email and password.
  - `Signin`: Authenticates user and issues tokens.
  - `ChangePassword`: Updates user password.
  - `RefreshToken`: Refreshes access token.
  - `ChangeRole`: Updates user role (admin/super_admin).
  - `Signout`: Revokes refresh token.
