# Location-Wise Pricing Implementation

## Overview
Implemented location-wise pricing for products, allowing different cost and selling prices per location instead of a single global price. Products can now be filtered by location in both single and list endpoints.

## Changes Made

### 1. Product Model Update ([shared/models/inventory.go](shared/models/inventory.go))

#### Added LocationPrice Struct
```go
type LocationPrice struct {
    LocationID   primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
    LocationName string             `bson:"location_name" json:"location_name"`
    CostPrice    float64            `bson:"cost_price" json:"cost_price"`
    SellingPrice float64            `bson:"selling_price" json:"selling_price"`
    MRP          float64            `bson:"mrp" json:"mrp"`
    Currency     string             `bson:"currency" json:"currency"`
    IsActive     bool               `bson:"is_active" json:"is_active"`
    CreatedAt    int64              `bson:"created_at" json:"created_at"`
    ModifiedAt   int64              `bson:"modified_at" json:"modified_at"`
}
```

#### Updated Product Model
- **Added**: `LocationPrices []LocationPrice` - Array to store location-wise pricing
- **Removed**: Old single-price fields (`cost_price`, `selling_price`, `mrp`, `currency`, `standard_cost`)

### 2. Service Layer Updates ([services/product-service/internal/service/product_service.go](services/product-service/internal/service/product_service.go))

#### Added LocationPriceRequest DTO
```go
type LocationPriceRequest struct {
    LocationID   string  `json:"location_id" binding:"required"`
    LocationName string  `json:"location_name"`
    CostPrice    float64 `json:"cost_price"`
    SellingPrice float64 `json:"selling_price"`
    MRP          float64 `json:"mrp"`
    Currency     string  `json:"currency"`
    IsActive     bool    `json:"is_active"`
}
```

#### Updated CreateProductRequest
- **Added**: `LocationPrices []LocationPriceRequest` - For creating products with location-wise prices
- **Removed**: Old single-price fields

#### Updated UpdateProductRequest
- **Added**: `LocationPrices []LocationPriceRequest` - For updating location-wise prices
- **Removed**: Old single-price fields

#### Updated ProductResponse
- **Added**: `LocationPrices []models.LocationPrice` - Returns location-wise prices in GET requests
- **Removed**: Old single-price fields

#### Enhanced Service Methods
- `createProductRequestToModel()`: Converts LocationPriceRequest array to LocationPrice models with timestamps and default currency
- `applyProductUpdates()`: Handles updates to location-wise prices, preserving CreatedAt timestamps for existing locations

### 3. Handler Layer Updates ([services/product-service/internal/handlers/product_handler.go](services/product-service/internal/handlers/product_handler.go))

#### GetProduct Endpoint
- **Added**: `location_id` query parameter to filter prices for a specific location
- Filters the `location_prices` array when `location_id` is provided
- Returns all location prices if no filter is specified

#### ListProducts Endpoint
- **Added**: `location_id` query parameter to filter products by location
- Only returns products that have prices defined for the specified location
- Works in combination with other filters (category, brand, search, etc.)

### 4. Repository Layer Updates ([services/product-service/internal/repository/product_repository.go](services/product-service/internal/repository/product_repository.go))

#### FindByOrganization Method
- **Added**: Location filtering using MongoDB array query `location_prices.location_id`
- Filters products that have at least one price entry for the specified location

## API Usage

### Creating a Product with Location-Wise Pricing

```json
POST /api/products
{
  "organization_id": "507f1f77bcf86cd799439011",
  "sku": "PROD-001",
  "name": "Sample Product",
  "location_prices": [
    {
      "location_id": "507f1f77bcf86cd799439012",
      "location_name": "Warehouse A",
      "cost_price": 100.00,
      "selling_price": 150.00,
      "mrp": 180.00,
      "currency": "USD",
      "is_active": true
    },
    {
      "location_id": "507f1f77bcf86cd799439013",
      "location_name": "Warehouse B",
      "cost_price": 105.00,
      "selling_price": 160.00,
      "mrp": 185.00,
      "currency": "USD",
      "is_active": true
    }
  ]
}
```

### Updating Product Location Prices

```json
PUT /api/products/:id
{
  "location_prices": [
    {
      "location_id": "507f1f77bcf86cd799439012",
      "location_name": "Warehouse A",
      "cost_price": 110.00,
      "selling_price": 160.00,
      "mrp": 190.00,
      "currency": "USD",
      "is_active": true
    },
    {
      "location_id": "507f1f77bcf86cd799439014",
      "location_name": "Warehouse C",
      "cost_price": 108.00,
      "selling_price": 155.00,
      "mrp": 185.00,
      "currency": "USD",
      "is_active": true
    }
  ]
}
```

### Getting Product with Location Filter

```json
GET /api/products/:id?location_id=507f1f77bcf86cd799439012
```

Returns product with prices filtered to only the specified location.

### Listing Products by Location

```json
GET /api/products?organization_id=507f1f77bcf86cd799439011&location_id=507f1f77bcf86cd799439012
```

Returns only products that have pricing defined for the specified location.

### Response Format (GET /api/products/:id or GET /api/products)

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "organization_id": "507f1f77bcf86cd799439011",
    "sku": "PROD-001",
    "name": "Sample Product",
    "location_prices": [
      {
        "location_id": "507f1f77bcf86cd799439012",
        "location_name": "Warehouse A",
        "cost_price": 110.00,
        "selling_price": 160.00,
        "mrp": 190.00,
        "currency": "USD",
        "is_active": true,
        "created_at": 1736697600000,
        "modified_at": 1736697600000
      },
      {
        "location_id": "507f1f77bcf86cd799439014",
        "location_name": "Warehouse C",
        "cost_price": 108.00,
        "selling_price": 155.00,
        "mrp": 185.00,
        "currency": "USD",
        "is_active": true,
        "created_at": 1736697600000,
        "modified_at": 1736697600000
      }
    ]
  }
}
```

## Database Schema

The product document in MongoDB now includes:

```javascript
{
  _id: ObjectId("..."),
  organization_id: ObjectId("..."),
  sku: "PROD-001",
  name: "Sample Product",
  // ... other fields ...
  
  // Location-wise pricing
  location_prices: [
    {
      location_id: ObjectId("..."),
      location_name: "Warehouse A",
      cost_price: 110.00,
      selling_price: 160.00,
      mrp: 190.00,
      currency: "USD",
      is_active: true,
      created_at: 1736697600000,
      modified_at: 1736697600000
    },
    // ... more locations
  ]
}
```

## Important Notes

1. **Breaking Change**: Old single-price fields have been completely removed
2. **Migration Required**: Existing products need data migration (see Migration Guide below)
3. **Location Filtering**: Products can be filtered by location in list and get endpoints
4. **Flexible Pricing**: Different prices for different locations/warehouses
5. **Location Tracking**: Track which locations have which prices
6. **Active/Inactive Control**: Enable/disable pricing for specific locations
7. **Audit Trail**: Created and modified timestamps for each location price
8. **Currency Support**: Different currencies per location if needed

## Full API Documentation

For complete API documentation including all endpoints, request/response examples, filtering options, and integration guides, see:
**[Product API Documentation](services/product-service/PRODUCT_API_DOCUMENTATION.md)**

## Next Steps

Consider implementing:
1. **Data Migration Script**: Convert existing products with old price fields to location-wise pricing
2. Location validation to ensure location_id exists in the system
3. Bulk location price updates endpoint
4. Price history tracking for each location
5. Location-based stock and pricing reports
6. Validation to require at least one active location price per product

## Migration Guide

If you have existing products with old pricing fields, you'll need to migrate them:

```javascript
// Example migration script for MongoDB
db.products.find({ cost_price: { $exists: true } }).forEach(function(product) {
  // Create default location price from old fields
  const defaultLocation = {
    location_id: ObjectId("YOUR_DEFAULT_LOCATION_ID"),
    location_name: "Main Warehouse",
    cost_price: product.cost_price || 0,
    selling_price: product.selling_price || 0,
    mrp: product.mrp || 0,
    currency: product.currency || "USD",
    is_active: true,
    created_at: new Date().getTime(),
    modified_at: new Date().getTime()
  };
  
  // Update product with location prices and remove old fields
  db.products.updateOne(
    { _id: product._id },
    {
      $set: { location_prices: [defaultLocation] },
      $unset: { 
        cost_price: "", 
        selling_price: "", 
        mrp: "", 
        currency: "",
        standard_cost: ""
      }
    }
  );
});
```

## Notes
