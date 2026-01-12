# Product Service API Documentation

## Overview
The Product Service manages product information with location-wise pricing. Each product can have different cost and selling prices for different locations/warehouses.

**Base URL**: `/api/v1/products`

**Authentication**: All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

---

## Endpoints

### 1. Create Product

Creates a new product with location-wise pricing.

**Endpoint**: `POST /api/v1/products`

**Request Body**:
```json
{
  "organization_id": "507f1f77bcf86cd799439011",
  "sku": "PROD-001",
  "barcode": "1234567890123",
  "name": "Sample Product",
  "description": "Product description",
  "type": "finished_goods",
  "status": "active",
  "category_id": "507f1f77bcf86cd799439012",
  "subcategory_id": "507f1f77bcf86cd799439013",
  "brand_id": "507f1f77bcf86cd799439014",
  "manufacturer_id": "507f1f77bcf86cd799439015",
  "track_inventory": true,
  "track_batches": false,
  "track_serial_numbers": false,
  "valuation_method": "fifo",
  "base_unit_id": "507f1f77bcf86cd799439016",
  "allowed_unit_ids": ["507f1f77bcf86cd799439016", "507f1f77bcf86cd799439017"],
  "weight": 1.5,
  "weight_unit": "kg",
  "length": 10.0,
  "width": 5.0,
  "height": 3.0,
  "dimension_unit": "cm",
  "volume": 150.0,
  "volume_unit": "cm3",
  "location_prices": [
    {
      "location_id": "507f1f77bcf86cd799439020",
      "location_name": "Main Warehouse",
      "cost_price": 100.00,
      "selling_price": 150.00,
      "mrp": 180.00,
      "currency": "USD",
      "is_active": true
    },
    {
      "location_id": "507f1f77bcf86cd799439021",
      "location_name": "Regional Store",
      "cost_price": 105.00,
      "selling_price": 160.00,
      "mrp": 185.00,
      "currency": "USD",
      "is_active": true
    }
  ],
  "tax_category_id": "507f1f77bcf86cd799439022",
  "hsn_code": "1234",
  "sac_code": "5678",
  "reorder_level": 10,
  "reorder_quantity": 50,
  "min_stock_level": 5,
  "max_stock_level": 100,
  "safety_stock": 15,
  "default_supplier_id": "507f1f77bcf86cd799439023",
  "supplier_ids": ["507f1f77bcf86cd799439023", "507f1f77bcf86cd799439024"],
  "lead_time_days": 7,
  "shelf_life_days": 365,
  "requires_qc": true,
  "perishable": false,
  "hazardous": false,
  "images": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
  "thumbnail": "https://example.com/thumbnail.jpg",
  "specifications": {
    "color": "Blue",
    "material": "Plastic"
  },
  "metadata": {
    "custom_field": "value"
  }
}
```

**Response**: `201 Created`
```json
{
  "status": "success",
  "message": "Product created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439025",
    "organization_id": "507f1f77bcf86cd799439011",
    "sku": "PROD-001",
    "name": "Sample Product",
    "location_prices": [
      {
        "location_id": "507f1f77bcf86cd799439020",
        "location_name": "Main Warehouse",
        "cost_price": 100.00,
        "selling_price": 150.00,
        "mrp": 180.00,
        "currency": "USD",
        "is_active": true,
        "created_at": 1736697600000,
        "modified_at": 1736697600000
      }
    ],
    "created_at": "2026-01-12T10:00:00Z",
    "updated_at": "2026-01-12T10:00:00Z"
  }
}
```

**Validation Rules**:
- `organization_id`, `sku`, `name` are required
- `location_id` in `location_prices` must be a valid ObjectID
- SKU must be unique within the organization
- If `subcategory_id` is provided, `category_id` must also be provided
- Product types: `raw_material`, `finished_goods`, `semi_finished`, `consumable`, `service`
- Product status: `active`, `inactive`, `discontinued`
- Valuation methods: `fifo`, `lifo`, `weighted_average`, `standard`

---

### 2. Get Product by ID

Retrieves a single product by its ID. Optionally filter location prices by location.

**Endpoint**: `GET /api/v1/products/:id`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `location_id` | string | No | Filter to show only prices for this location |

**Example Request**:
```
GET /api/v1/products/507f1f77bcf86cd799439025
GET /api/v1/products/507f1f77bcf86cd799439025?location_id=507f1f77bcf86cd799439020
```

**Response**: `200 OK`
```json
{
  "status": "success",
  "message": "Product retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439025",
    "organization_id": "507f1f77bcf86cd799439011",
    "sku": "PROD-001",
    "barcode": "1234567890123",
    "name": "Sample Product",
    "description": "Product description",
    "type": "finished_goods",
    "status": "active",
    "category_id": "507f1f77bcf86cd799439012",
    "location_prices": [
      {
        "location_id": "507f1f77bcf86cd799439020",
        "location_name": "Main Warehouse",
        "cost_price": 100.00,
        "selling_price": 150.00,
        "mrp": 180.00,
        "currency": "USD",
        "is_active": true,
        "created_at": 1736697600000,
        "modified_at": 1736697600000
      }
    ],
    "category": {
      "id": "507f1f77bcf86cd799439012",
      "name": "Electronics",
      "code": "ELEC"
    },
    "total_stock": 45.0,
    "available_stock": 40.0,
    "created_at": "2026-01-12T10:00:00Z",
    "updated_at": "2026-01-12T10:00:00Z"
  }
}
```

---

### 3. Get Product by SKU

Retrieves a product by its SKU code.

**Endpoint**: `GET /api/v1/products/sku/:sku`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `organization_id` | string | Yes | Organization ID |

**Example Request**:
```
GET /api/v1/products/sku/PROD-001?organization_id=507f1f77bcf86cd799439011
```

**Response**: `200 OK` - Same structure as Get Product by ID

---

### 4. List Products

Retrieves a paginated list of products with optional filters.

**Endpoint**: `GET /api/v1/products`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `organization_id` | string | Yes | Organization ID |
| `category_id` | string | No | Filter by category |
| `subcategory_id` | string | No | Filter by subcategory |
| `brand_id` | string | No | Filter by brand |
| `status` | string | No | Filter by status (`active`, `inactive`, `discontinued`) |
| `type` | string | No | Filter by product type |
| `track_inventory` | boolean | No | Filter by inventory tracking |
| `location_id` | string | No | Filter products that have prices for this location |
| `search` | string | No | Search by name, SKU, or description |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10, max: 100) |

**Example Requests**:
```
GET /api/v1/products?organization_id=507f1f77bcf86cd799439011&page=1&limit=10
GET /api/v1/products?organization_id=507f1f77bcf86cd799439011&location_id=507f1f77bcf86cd799439020
GET /api/v1/products?organization_id=507f1f77bcf86cd799439011&search=laptop&status=active
GET /api/v1/products?organization_id=507f1f77bcf86cd799439011&category_id=507f1f77bcf86cd799439012
```

**Response**: `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439025",
      "organization_id": "507f1f77bcf86cd799439011",
      "sku": "PROD-001",
      "name": "Sample Product",
      "type": "finished_goods",
      "status": "active",
      "location_prices": [
        {
          "location_id": "507f1f77bcf86cd799439020",
          "location_name": "Main Warehouse",
          "cost_price": 100.00,
          "selling_price": 150.00,
          "mrp": 180.00,
          "currency": "USD",
          "is_active": true,
          "created_at": 1736697600000,
          "modified_at": 1736697600000
        }
      ],
      "total_stock": 45.0,
      "available_stock": 40.0,
      "created_at": "2026-01-12T10:00:00Z",
      "updated_at": "2026-01-12T10:00:00Z"
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

**Use Cases**:
- **Get all products for a location**: Use `location_id` filter to get products available at a specific warehouse
- **Search products**: Use `search` parameter to find products by name, SKU, or description
- **Filter by category**: Combine `category_id` and `status=active` to get active products in a category
- **Multi-location filtering**: Products will only appear if they have prices defined for the specified location

---

### 5. Update Product

Updates an existing product. All fields are optional except the ID.

**Endpoint**: `PUT /api/v1/products/:id`

**Request Body**:
```json
{
  "name": "Updated Product Name",
  "description": "Updated description",
  "status": "active",
  "location_prices": [
    {
      "location_id": "507f1f77bcf86cd799439020",
      "location_name": "Main Warehouse",
      "cost_price": 110.00,
      "selling_price": 165.00,
      "mrp": 190.00,
      "currency": "USD",
      "is_active": true
    },
    {
      "location_id": "507f1f77bcf86cd799439026",
      "location_name": "New Branch",
      "cost_price": 115.00,
      "selling_price": 170.00,
      "mrp": 195.00,
      "currency": "USD",
      "is_active": true
    }
  ],
  "images": ["https://example.com/new-image.jpg"]
}
```

**Response**: `200 OK`
```json
{
  "status": "success",
  "message": "Product updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439025",
    "name": "Updated Product Name",
    "location_prices": [
      {
        "location_id": "507f1f77bcf86cd799439020",
        "location_name": "Main Warehouse",
        "cost_price": 110.00,
        "selling_price": 165.00,
        "mrp": 190.00,
        "currency": "USD",
        "is_active": true,
        "created_at": 1736697600000,
        "modified_at": 1736701200000
      }
    ],
    "updated_at": "2026-01-12T11:00:00Z"
  }
}
```

**Important Notes**:
- Updating `location_prices` **replaces all existing location prices** with the new array
- To preserve existing location prices, fetch the product first, modify the array, and send the complete array
- `created_at` for each location price is preserved if the location already existed
- SKU can be updated but must remain unique within the organization
- Cannot update `organization_id`

---

### 6. Delete Product

Soft deletes a product. Products with existing stock cannot be deleted.

**Endpoint**: `DELETE /api/v1/products/:id`

**Example Request**:
```
DELETE /api/v1/products/507f1f77bcf86cd799439025
```

**Response**: `200 OK`
```json
{
  "status": "success",
  "message": "Product deleted successfully",
  "data": null
}
```

**Validation**:
- Cannot delete products with `total_stock > 0`
- Product must belong to the user's organization

---

### 7. Get Low Stock Products

Retrieves products that are below their reorder level.

**Endpoint**: `GET /api/v1/products/low-stock`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `organization_id` | string | Yes | Organization ID |

**Example Request**:
```
GET /api/v1/products/low-stock?organization_id=507f1f77bcf86cd799439011
```

**Response**: `200 OK`
```json
{
  "status": "success",
  "message": "Low stock products retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439025",
      "sku": "PROD-001",
      "name": "Sample Product",
      "total_stock": 3.0,
      "reorder_level": 10,
      "reorder_quantity": 50,
      "location_prices": [...]
    }
  ]
}
```
