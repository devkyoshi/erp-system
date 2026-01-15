# Full Documentation: Products and Unit Charts

This document provides a comprehensive guide to the Product Service's data models and API endpoints, focusing on Products and Unit Conversions (Unit Charts).

## 1. Data Models

### 1.1 Product Model
The `Product` entity is the core of the inventory system. It represents items that are bought, sold, or manufactured.

**Key Fields:**
- `_id`: Unique Identifier (MongoDB ObjectID).
- `sku` (String): Stock Keeping Unit, unique per organization.
- `name` (String): Product name.
- `type` (Enum): `raw_material`, `finished_goods`, `service`, etc.
- `status` (Enum): `active`, `inactive`, `discontinued`.
- `organization_id`: Reference to the tenant organization.

**Inventory & Units:**
- `base_unit_id`: The primary unit for stock keeping (e.g., "Pieces").
- `allowed_unit_ids`: List of other units this product can be transacted in (e.g., "Box", "Carton").
- `track_inventory` (Boolean): Whether to track stock levels.

**Location-wise Pricing (`location_prices`):**
Products can have different prices in different locations (stores/warehouses).
- `location_id`: Reference to the Location.
- `unit_id`: **(New)** The specific unit for this price (e.g., Price per "Box" at "Warehouse A").
- `cost_price`: Cost to the company.
- `selling_price`: Selling price.
- `currency`: Currency code (e.g., "USD").

### 1.2 Unit Model
Represents a standard Unit of Measure (UOM).
- `name`: Full name (e.g., "Kilogram").
- `code`: Abbreviation (e.g., "kg").
- `unit_type`: Classification (e.g., "weight", "quantity").
- `is_base_unit`: Flag indicating if this is a reference unit for its type.

### 1.3 Unit Chart (Conversion) Model
Defines the conversion rules between units. This allows the system to convert stock and prices between different UOMs.
- `from_unit_id`: Source Unit.
- `to_unit_id`: Target Unit.
- `conversion_rate`: Multiplier.
  - *Logic*: `1 FromUnit = ConversionRate * ToUnit`
  - *Example*: 1 Box = 12 Pieces. `From`="Box", `To`="Pieces", `Rate`=12.

---

## 2. API Reference

### 2.1 Units API
**Endpoint**: `GET /api/v1/units`

Returns a hierarchical view of units and their available conversions.

**Response Structure:**
```json
[
  {
    "id": "...",
    "name": "Box",
    "code": "BOX",
    "unit_type": "quantity",
    "conversions": [  // Derived from Unit Charts
      {
        "to_unit_id": "...",     // ID of "Pieces"
        "to_unit_name": "Pieces",
        "to_unit_code": "PCS",
        "conversion_rate": 12.0  // 1 Box = 12 Pieces
      }
    ]
  }
]
```

### 2.2 Products API

#### Create Product `POST /api/v1/products`
Create a product with location-specific pricing and units.

**Payload:**
```json
{
  "name": "Soda Can",
  "sku": "SODA-001",
  "base_unit_id": "...", // ID for "Can"
  "location_prices": [
    {
      "location_id": "...",
      "unit_id": "...",    // ID for "Six-Pack"
      "selling_price": 5.99
    }
  ]
}
```

#### Get Product `GET /api/v1/products/{id}`
Retrieves product details. The response auto-populates the full `unit` object within `location_prices` for display purposes.

**Response Snippet:**
```json
{
  "sku": "SODA-001",
  "location_prices": [
    {
      "location_id": "...",
      "price": 5.99,
      "unit_id": "...",
      "unit": {          // Populated details
        "name": "Six-Pack",
        "code": "6PK"
      }
    }
  ]
}
```
