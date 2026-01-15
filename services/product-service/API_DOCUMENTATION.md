# Product Service API Documentation

## Overview
This service manages products, brands, categories, and units.

## Base URL
`/api/v1`

---

## 1. Units

### Get All Units
Retrieves a list of all active units of measure along with their conversion rules.

- **Endpoint**: `GET /units`
- **Auth Required**: Yes

#### Response
```json
{
  "success": true,
  "status": 200,
  "data": [
    {
      "id": "60d5ec49f1b2c62c8428a1e1",
      "name": "Box",
      "code": "BOX",
      "unit_type": "quantity",
      "is_base_unit": false,
      "conversions": [
        {
          "to_unit_id": "60d5ec49f1b2c62c8428a1e0",
          "to_unit_name": "Pieces",
          "to_unit_code": "PCS",
          "conversion_rate": 12
        }
      ]
    },
    {
      "id": "60d5ec49f1b2c62c8428a1e0",
      "name": "Pieces",
      "code": "PCS",
      "unit_type": "quantity",
      "is_base_unit": true,
      "conversions": []
    }
  ],
  "message": "Units retrieved successfully"
}
```

---

## 2. Products

### Product Object Structure (Key Fields)
```json
{
  "id": "60d5ec49f1b2c62c8428a1e5",
  "name": "Example Product",
  "sku": "PROD-001",
  "location_prices": [
    {
      "location_id": "60d5ec49f1b2c62c8428a1e2",
      "location_name": "Main Warehouse",
      "unit_id": "60d5ec49f1b2c62c8428a1e1", // Added Field
      "unit": {                          // Populated Field (Response Only)
        "id": "60d5ec49f1b2c62c8428a1e1",
        "name": "Box",
        "code": "BOX",
        "symbol": "box"
      },
      "cost_price": 100,
      "selling_price": 150,
      "mrp": 160,
      "currency": "USD"
    }
  ]
  // ... other fields
}
```

### Create Product
Creates a new product.

- **Endpoint**: `POST /products`
- **Auth Required**: Yes

#### Request Body
```json
{
  "organization_id": "60d5ec49f1b2c62c8428a1e9",
  "name": "Premium Widget",
  "sku": "WIDGET-001",
  "location_prices": [
    {
      "location_id": "60d5ec49f1b2c62c8428a1e2",
      "location_name": "Main Store",
      "unit_id": "60d5ec49f1b2c62c8428a1e1", // Optional: Specify unit for this price
      "cost_price": 50.00,
      "selling_price": 85.00,
      "mrp": 100.00,
      "currency": "USD",
      "is_active": true
    }
  ]
  // ... other standard product fields (category_id, brand_id, etc.)
}
```

#### Response
```json
{
  "success": true,
  "status": 201,
  "data": {
    "id": "60d5ec49f1b2c62c8428a1e5",
    "name": "Premium Widget",
    "sku": "WIDGET-001",
    // ... created product data
  },
  "message": "Product created successfully"
}
```

### Get Product
Retrieves a single product by ID. The `location_prices` will include full `unit` details.

- **Endpoint**: `GET /products/{id}`
- **Auth Required**: Yes

#### Response
```json
{
  "success": true,
  "status": 200,
  "data": {
    "id": "60d5ec49f1b2c62c8428a1e5",
    "sku": "WIDGET-001",
    "location_prices": [
      {
        "location_id": "60d5ec49f1b2c62c8428a1e2",
        "location_name": "Main Store",
        "unit_id": "60d5ec49f1b2c62c8428a1e1",
        "unit": {
           "id": "60d5ec49f1b2c62c8428a1e1",
           "name": "Box",
           "code": "BOX",
           "unit_type": "quantity"
        },
        "selling_price": 85.00
        // ...
      }
    ]
    // ...
  },
  "message": "Product retrieved successfully"
}
```

### List Products
Retrieves a list of products. Location prices in the list will also include unit details.

- **Endpoint**: `GET /products`
- **Query Params**:
  - `organization_id` (Required)
  - `page`, `limit`
  - `search`, `category_id`, `brand_id`, etc.
  - `location_id` (Filter by presence of price in location)

#### Response
```json
{
  "success": true,
  "status": 200,
  "data": [
    {
      "id": "60d5ec49f1b2c62c8428a1e5",
      "sku": "WIDGET-001",
      "location_prices": [
        {
          "location_id": "...",
          "unit": { ... },
          "selling_price": 85.00
        }
      ]
    }
    // ...
  ],
  "pagination": { ... },
  "message": "Products retrieved successfully"
}
```

### Update Product
Updates an existing product. You can update the `unit_id` in location prices.

- **Endpoint**: `PUT /products/{id}`

#### Request Body
```json
{
  "name": "Updated Name",
  "location_prices": [
    {
      "location_id": "60d5ec49f1b2c62c8428a1e2",
      "unit_id": "60d5ec49f1b2c62c8428a1e0", // Changing unit to 'Pieces'
      "selling_price": 8.00
    }
  ]
}
```
