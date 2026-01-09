# Category & Product API Documentation

## Base URL
```
https://api.example.com/api/v1
```

All endpoints require authentication via JWT Bearer token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## Category Endpoints

### 1. Create Category

Create a new category with optional embedded subcategories.

**Endpoint:** `POST /categories`

**Request Body:**
```json
{
  "organization_id": "507f1f77bcf86cd799439011",
  "name": "Electronics",
  "code": "ELEC",
  "description": "Electronic products and accessories",
  "is_active": true,
  "metadata": {
    "display_order": 1,
    "icon": "electronics-icon.svg"
  },
  "subcategories": [
    {
      "name": "Laptops",
      "code": "LAP",
      "description": "Laptop computers and accessories",
      "is_active": true,
      "metadata": {
        "display_order": 1
      }
    },
    {
      "name": "Mobile Phones",
      "code": "PHN",
      "description": "Smartphones and feature phones",
      "is_active": true
    }
  ]
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Electronics",
    "code": "ELEC",
    "description": "Electronic products and accessories",
    "is_active": true,
    "product_count": 0,
    "subcategories": [
      {
        "id": "507f1f77bcf86cd799439021",
        "name": "Laptops",
        "code": "LAP",
        "description": "Laptop computers and accessories",
        "is_active": true,
        "product_count": 0,
        "metadata": {
          "display_order": 1
        },
        "created_at": "2026-01-09T10:30:00Z",
        "updated_at": "2026-01-09T10:30:00Z"
      },
      {
        "id": "507f1f77bcf86cd799439022",
        "name": "Mobile Phones",
        "code": "PHN",
        "description": "Smartphones and feature phones",
        "is_active": true,
        "product_count": 0,
        "created_at": "2026-01-09T10:30:00Z",
        "updated_at": "2026-01-09T10:30:00Z"
      }
    ],
    "metadata": {
      "display_order": 1,
      "icon": "electronics-icon.svg"
    },
    "created_at": "2026-01-09T10:30:00Z",
    "updated_at": "2026-01-09T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation error or duplicate category name
- `401 Unauthorized` - Invalid or missing authentication token
- `403 Forbidden` - User doesn't have permission for this organization
- `500 Internal Server Error` - Server error

---

### 2. Get Category by ID

Retrieve a single category with its embedded subcategories.

**Endpoint:** `GET /categories/:id`

**Path Parameters:**
- `id` (required) - Category ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Category retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Electronics",
    "code": "ELEC",
    "description": "Electronic products and accessories",
    "is_active": true,
    "product_count": 45,
    "subcategories": [
      {
        "id": "507f1f77bcf86cd799439021",
        "name": "Laptops",
        "code": "LAP",
        "description": "Laptop computers and accessories",
        "is_active": true,
        "product_count": 20,
        "created_at": "2026-01-09T10:30:00Z",
        "updated_at": "2026-01-09T10:30:00Z"
      },
      {
        "id": "507f1f77bcf86cd799439022",
        "name": "Mobile Phones",
        "code": "PHN",
        "description": "Smartphones and feature phones",
        "is_active": true,
        "product_count": 25,
        "created_at": "2026-01-09T10:30:00Z",
        "updated_at": "2026-01-09T10:30:00Z"
      }
    ],
    "created_at": "2026-01-09T10:30:00Z",
    "updated_at": "2026-01-09T10:30:00Z"
  }
}
```

---

### 3. List Categories

Retrieve all categories for an organization with optional filters and pagination.

**Endpoint:** `GET /categories`

**Query Parameters:**
- `organization_id` (required) - Organization ObjectID
- `is_active` (optional) - Filter by active status (true/false)
- `q` (optional) - Search query (searches name, code, description)
- `page` (optional, default: 1) - Page number
- `limit` (optional, default: 10, max: 100) - Items per page

**Example Request:**
```
GET /categories?organization_id=507f1f77bcf86cd799439011&is_active=true&q=electronics&page=1&limit=10
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "organization_id": "507f1f77bcf86cd799439011",
      "name": "Electronics",
      "code": "ELEC",
      "is_active": true,
      "product_count": 45,
      "subcategories": [
        {
          "id": "507f1f77bcf86cd799439021",
          "name": "Laptops",
          "code": "LAP",
          "is_active": true,
          "product_count": 20
        }
      ],
      "created_at": "2026-01-09T10:30:00Z",
      "updated_at": "2026-01-09T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 4. Get Category Tree

Retrieve all categories for an organization in a flat structure (subcategories are embedded).

**Endpoint:** `GET /categories/tree`

**Query Parameters:**
- `organization_id` (required) - Organization ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Category tree retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "name": "Electronics",
      "code": "ELEC",
      "is_active": true,
      "product_count": 45,
      "subcategories": [
        {
          "id": "507f1f77bcf86cd799439021",
          "name": "Laptops",
          "product_count": 20
        },
        {
          "id": "507f1f77bcf86cd799439022",
          "name": "Mobile Phones",
          "product_count": 25
        }
      ]
    },
    {
      "id": "507f1f77bcf86cd799439030",
      "name": "Furniture",
      "code": "FURN",
      "is_active": true,
      "product_count": 30,
      "subcategories": []
    }
  ]
}
```

---

### 5. Get Category Children (Subcategories)

Retrieve the subcategories for a specific category.

**Endpoint:** `GET /categories/:id/children`

**Path Parameters:**
- `id` (required) - Category ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Children categories retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439021",
      "name": "Laptops",
      "code": "LAP",
      "description": "Laptop computers and accessories",
      "is_active": true,
      "product_count": 20,
      "created_at": "2026-01-09T10:30:00Z",
      "updated_at": "2026-01-09T10:30:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439022",
      "name": "Mobile Phones",
      "code": "PHN",
      "description": "Smartphones and feature phones",
      "is_active": true,
      "product_count": 25,
      "created_at": "2026-01-09T10:30:00Z",
      "updated_at": "2026-01-09T10:30:00Z"
    }
  ]
}
```

---

### 6. Update Category

Update a category and manage its subcategories (add, update, or remove).

**Endpoint:** `PUT /categories/:id`

**Path Parameters:**
- `id` (required) - Category ObjectID

**Request Body:**
```json
{
  "name": "Consumer Electronics",
  "description": "Updated description",
  "is_active": true,
  "add_subcategories": [
    {
      "name": "Tablets",
      "code": "TAB",
      "description": "Tablet computers",
      "is_active": true
    }
  ],
  "update_subcategories": [
    {
      "id": "507f1f77bcf86cd799439021",
      "name": "Gaming Laptops",
      "description": "High-performance gaming laptops"
    }
  ],
  "remove_subcategories": ["507f1f77bcf86cd799439022"]
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Consumer Electronics",
    "description": "Updated description",
    "is_active": true,
    "subcategories": [
      {
        "id": "507f1f77bcf86cd799439021",
        "name": "Gaming Laptops",
        "description": "High-performance gaming laptops",
        "is_active": true,
        "updated_at": "2026-01-09T11:00:00Z"
      },
      {
        "id": "507f1f77bcf86cd799439023",
        "name": "Tablets",
        "code": "TAB",
        "description": "Tablet computers",
        "is_active": true,
        "created_at": "2026-01-09T11:00:00Z",
        "updated_at": "2026-01-09T11:00:00Z"
      }
    ],
    "updated_at": "2026-01-09T11:00:00Z"
  }
}
```

**Notes:**
- Removing a subcategory soft-deletes it (sets `deleted_at`)
- You can add, update, and remove subcategories in the same request
- All operations are performed atomically

---

### 7. Delete Category

Soft delete a category (sets `deleted_at` timestamp).

**Endpoint:** `DELETE /categories/:id`

**Path Parameters:**
- `id` (required) - Category ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Category deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request` - Category has active subcategories or products
- `404 Not Found` - Category not found

**Validation:**
- Cannot delete category with active (non-deleted) subcategories
- Cannot delete category with associated products

---

## Product Endpoints

### 1. Create Product

Create a new product with category and optional subcategory.

**Endpoint:** `POST /products`

**Request Body:**
```json
{
  "organization_id": "507f1f77bcf86cd799439011",
  "sku": "LAP-DELL-XPS15-001",
  "barcode": "123456789012",
  "name": "Dell XPS 15 Laptop",
  "description": "15.6-inch laptop with Intel Core i7",
  "type": "finished",
  "status": "active",
  "category_id": "507f1f77bcf86cd799439020",
  "subcategory_id": "507f1f77bcf86cd799439021",
  "brand_id": "507f1f77bcf86cd799439040",
  "track_inventory": true,
  "track_batches": false,
  "track_serial_numbers": true,
  "valuation_method": "fifo",
  "cost_price": 999.99,
  "selling_price": 1299.99,
  "mrp": 1399.99,
  "currency": "USD",
  "reorder_level": 5,
  "reorder_quantity": 10,
  "min_stock_level": 3,
  "max_stock_level": 50,
  "images": [
    "https://cdn.example.com/products/dell-xps15-1.jpg",
    "https://cdn.example.com/products/dell-xps15-2.jpg"
  ],
  "specifications": {
    "processor": "Intel Core i7-12700H",
    "ram": "16GB DDR5",
    "storage": "512GB NVMe SSD",
    "display": "15.6\" FHD"
  },
  "metadata": {
    "warranty_months": 24,
    "supplier_part_no": "XPS-15-9520"
  }
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439050",
    "organization_id": "507f1f77bcf86cd799439011",
    "sku": "LAP-DELL-XPS15-001",
    "barcode": "123456789012",
    "name": "Dell XPS 15 Laptop",
    "description": "15.6-inch laptop with Intel Core i7",
    "type": "finished",
    "status": "active",
    "category_id": "507f1f77bcf86cd799439020",
    "subcategory_id": "507f1f77bcf86cd799439021",
    "brand_id": "507f1f77bcf86cd799439040",
    "track_inventory": true,
    "total_stock": 0,
    "available_stock": 0,
    "cost_price": 999.99,
    "selling_price": 1299.99,
    "mrp": 1399.99,
    "currency": "USD",
    "created_at": "2026-01-09T10:30:00Z",
    "updated_at": "2026-01-09T10:30:00Z"
  }
}
```

**Validation Rules:**
- `sku` must be unique within the organization
- If `subcategory_id` is provided, `category_id` must also be provided
- `subcategory_id` must exist in the specified category's subcategories array
- Subcategory must not be soft-deleted

**Error Responses:**
- `400 Bad Request` - Validation error (duplicate SKU, invalid IDs, etc.)
- `404 Not Found` - Category or subcategory not found
- `422 Unprocessable Entity` - Subcategory doesn't belong to specified category

---

### 2. Get Product by ID

Retrieve a single product with its category details.

**Endpoint:** `GET /products/:id`

**Path Parameters:**
- `id` (required) - Product ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439050",
    "sku": "LAP-DELL-XPS15-001",
    "name": "Dell XPS 15 Laptop",
    "description": "15.6-inch laptop with Intel Core i7",
    "category_id": "507f1f77bcf86cd799439020",
    "subcategory_id": "507f1f77bcf86cd799439021",
    "selling_price": 1299.99,
    "total_stock": 15,
    "available_stock": 12,
    "category": {
      "id": "507f1f77bcf86cd799439020",
      "name": "Electronics",
      "code": "ELEC"
    },
    "created_at": "2026-01-09T10:30:00Z",
    "updated_at": "2026-01-09T10:30:00Z"
  }
}
```

---

### 3. Get Product by SKU

Retrieve a product by its SKU.

**Endpoint:** `GET /products/sku/:sku`

**Path Parameters:**
- `sku` (required) - Product SKU

**Query Parameters:**
- `organization_id` (required) - Organization ObjectID

**Example:**
```
GET /products/sku/LAP-DELL-XPS15-001?organization_id=507f1f77bcf86cd799439011
```

**Response:** `200 OK` (same structure as Get Product by ID)

---

### 4. List Products

Retrieve products with filtering and pagination.

**Endpoint:** `GET /products`

**Query Parameters:**
- `organization_id` (required) - Organization ObjectID
- `category_id` (optional) - Filter by category
- `subcategory_id` (optional) - Filter by subcategory
- `brand_id` (optional) - Filter by brand
- `status` (optional) - Filter by status (active, inactive, discontinued)
- `type` (optional) - Filter by type (raw, component, finished, service)
- `track_inventory` (optional) - Filter by inventory tracking (true/false)
- `search` (optional) - Search in name, SKU, or description
- `page` (optional, default: 1) - Page number
- `limit` (optional, default: 10, max: 100) - Items per page

**Example Request:**
```
GET /products?organization_id=507f1f77bcf86cd799439011&category_id=507f1f77bcf86cd799439020&subcategory_id=507f1f77bcf86cd799439021&search=dell&page=1&limit=20
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439050",
      "sku": "LAP-DELL-XPS15-001",
      "name": "Dell XPS 15 Laptop",
      "category_id": "507f1f77bcf86cd799439020",
      "subcategory_id": "507f1f77bcf86cd799439021",
      "selling_price": 1299.99,
      "total_stock": 15,
      "available_stock": 12,
      "status": "active"
    },
    {
      "id": "507f1f77bcf86cd799439051",
      "sku": "LAP-DELL-XPS13-001",
      "name": "Dell XPS 13 Laptop",
      "category_id": "507f1f77bcf86cd799439020",
      "subcategory_id": "507f1f77bcf86cd799439021",
      "selling_price": 999.99,
      "total_stock": 20,
      "available_stock": 18,
      "status": "active"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

**Use Cases:**
- Get all products: `?organization_id=xxx`
- Get products in a category: `?organization_id=xxx&category_id=yyy`
- Get products in a subcategory: `?organization_id=xxx&category_id=yyy&subcategory_id=zzz`
- Search products: `?organization_id=xxx&search=laptop`

---

### 5. Update Product

Update an existing product.

**Endpoint:** `PUT /products/:id`

**Path Parameters:**
- `id` (required) - Product ObjectID

**Request Body (all fields optional):**
```json
{
  "name": "Dell XPS 15 9520 Laptop",
  "description": "Updated description",
  "category_id": "507f1f77bcf86cd799439020",
  "subcategory_id": "507f1f77bcf86cd799439021",
  "selling_price": 1199.99,
  "status": "active",
  "specifications": {
    "processor": "Intel Core i7-12700H",
    "ram": "32GB DDR5",
    "storage": "1TB NVMe SSD",
    "display": "15.6\" FHD"
  }
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439050",
    "name": "Dell XPS 15 9520 Laptop",
    "description": "Updated description",
    "selling_price": 1199.99,
    "updated_at": "2026-01-09T12:00:00Z"
  }
}
```

**Validation:**
- Same validation rules as create (SKU uniqueness, category/subcategory relationship)
- Cannot change SKU to one that's already in use

---

### 6. Delete Product

Soft delete a product (sets `deleted_at` timestamp).

**Endpoint:** `DELETE /products/:id`

**Path Parameters:**
- `id` (required) - Product ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Product deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request` - Product has existing stock (cannot delete with stock)
- `404 Not Found` - Product not found

---

### 7. Get Low Stock Products

Retrieve products below their reorder level.

**Endpoint:** `GET /products/low-stock`

**Query Parameters:**
- `organization_id` (required) - Organization ObjectID

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Low stock products retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439050",
      "sku": "LAP-DELL-XPS15-001",
      "name": "Dell XPS 15 Laptop",
      "available_stock": 3,
      "reorder_level": 5,
      "reorder_quantity": 10,
      "category_id": "507f1f77bcf86cd799439020",
      "subcategory_id": "507f1f77bcf86cd799439021"
    }
  ]
}
```

---

## Common Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "field": "sku",
      "issue": "SKU already exists"
    }
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or missing authentication token"
  }
}
```

### 403 Forbidden
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found"
  }
}
```

### 422 Unprocessable Entity
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Subcategory not found in the specified category"
  }
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred"
  }
}
```

---

## Data Models

### Category Model
```json
{
  "id": "ObjectId",
  "organization_id": "ObjectId",
  "name": "string",
  "code": "string",
  "description": "string",
  "is_active": "boolean",
  "product_count": "integer",
  "subcategories": [
    {
      "id": "ObjectId",
      "name": "string",
      "code": "string",
      "description": "string",
      "is_active": "boolean",
      "product_count": "integer",
      "metadata": {},
      "created_at": "datetime",
      "updated_at": "datetime",
      "deleted_at": "datetime|null"
    }
  ],
  "metadata": {},
  "created_at": "datetime",
  "updated_at": "datetime",
  "deleted_at": "datetime|null"
}
```

### Product Model
```json
{
  "id": "ObjectId",
  "organization_id": "ObjectId",
  "sku": "string (unique)",
  "barcode": "string",
  "name": "string",
  "description": "string",
  "type": "enum (raw|component|finished|service)",
  "status": "enum (active|inactive|discontinued)",
  "category_id": "ObjectId",
  "subcategory_id": "ObjectId",
  "brand_id": "ObjectId",
  "track_inventory": "boolean",
  "track_batches": "boolean",
  "track_serial_numbers": "boolean",
  "valuation_method": "enum (fifo|lifo|weighted_average)",
  "cost_price": "decimal",
  "selling_price": "decimal",
  "mrp": "decimal",
  "currency": "string",
  "total_stock": "decimal",
  "available_stock": "decimal",
  "allocated_stock": "decimal",
  "reorder_level": "integer",
  "reorder_quantity": "integer",
  "images": ["string"],
  "specifications": {},
  "metadata": {},
  "created_at": "datetime",
  "updated_at": "datetime",
  "deleted_at": "datetime|null"
}
```

---

## Best Practices

1. **Category & Subcategory Structure**
   - Keep categories broad and general
   - Use subcategories for specific product types
   - Maximum depth is 2 levels (category → subcategory)

2. **Product Organization**
   - Always specify both category and subcategory when applicable
   - Use consistent naming conventions for SKUs
   - Keep product descriptions clear and searchable

3. **Filtering**
   - Use category_id for broad filtering
   - Use subcategory_id for specific filtering
   - Combine with search for more precise results

4. **Performance**
   - Use pagination for large result sets
   - Consider caching frequently accessed categories
   - Index on category_id and subcategory_id for faster queries

5. **Data Integrity**
   - Validate category/subcategory relationships before creating products
   - Handle soft-deleted subcategories in queries
   - Update product counts when products are created/deleted
