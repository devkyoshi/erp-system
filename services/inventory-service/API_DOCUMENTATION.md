# Inventory Service API Documentation

## Overview

The Inventory Service manages all inventory-related operations including stock levels, stock movements, batch tracking, serial number management, stock adjustments, and physical inventory counts. All endpoints require JWT authentication.

**Base URL**: `/api/v1/inventory`

**Authentication**: All endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

---

## Table of Contents

1. [Stock Levels](#stock-levels)
2. [Stock Movements](#stock-movements)
3. [Batches](#batches)
4. [Stock Adjustments](#stock-adjustments)
5. [Inventory Counts](#inventory-counts)
6. [Serial Numbers](#serial-numbers)

---

## Stock Levels

Stock levels track the quantity of products available at different locations.

### 1. Get Stock Level

Get stock level information based on filters.

**Endpoint**: `GET /stock-levels`

**Query Parameters**:
- `product_id` (string, optional): Filter by product ID
- `location_id` (string, optional): Filter by location ID
- `organization_id` (string, optional): Filter by organization ID

**Response**:
```json
{
  "status": "success",
  "message": "Stock level retrieved successfully",
  "data": {
    "product_id": "507f1f77bcf86cd799439011",
    "location_id": "507f1f77bcf86cd799439012",
    "organization_id": "507f1f77bcf86cd799439013",
    "quantity_on_hand": 100,
    "quantity_available": 85,
    "quantity_allocated": 15,
    "quantity_committed": 10,
    "reorder_point": 20,
    "reorder_quantity": 50,
    "last_stock_count": "2026-01-01T10:00:00Z",
    "last_movement_date": "2026-01-01T12:00:00Z"
  }
}
```

### 2. Get Stock by Product

Get all stock levels for a specific product across all locations.

**Endpoint**: `GET /products/:product_id/stock`

**URL Parameters**:
- `product_id` (string, required): Product ID

**Response**:
```json
{
  "status": "success",
  "message": "Stock retrieved successfully",
  "data": [
    {
      "location_id": "507f1f77bcf86cd799439012",
      "location_name": "Main Warehouse",
      "quantity_on_hand": 100,
      "quantity_available": 85,
      "quantity_allocated": 15
    }
  ]
}
```

### 3. Get Stock by Location

Get all stock levels for a specific location.

**Endpoint**: `GET /locations/:location_id/stock`

**URL Parameters**:
- `location_id` (string, required): Location ID

**Response**:
```json
{
  "status": "success",
  "message": "Stock retrieved successfully",
  "data": [
    {
      "product_id": "507f1f77bcf86cd799439011",
      "product_name": "Widget A",
      "quantity_on_hand": 100,
      "quantity_available": 85
    }
  ]
}
```

### 4. Allocate Stock

Reserve stock for an order or customer.

**Endpoint**: `POST /stock/allocate`

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "location_id": "507f1f77bcf86cd799439012",
  "quantity": 10
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Stock allocated successfully",
  "data": null
}
```

### 5. Release Stock

Release previously allocated stock back to available inventory.

**Endpoint**: `POST /stock/release`

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "location_id": "507f1f77bcf86cd799439012",
  "quantity": 5
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Stock released successfully",
  "data": null
}
```

---

## Stock Movements

Stock movements track all inventory transactions (receipts, issues, transfers, adjustments).

### 1. Create Stock Movement

Record a new stock movement transaction.

**Endpoint**: `POST /organizations/:org_id/stock-movements`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "movement_type": "receipt",
  "from_location_id": "507f1f77bcf86cd799439012",
  "to_location_id": "507f1f77bcf86cd799439013",
  "quantity": 50,
  "unit_cost": 25.50,
  "reference_type": "purchase_order",
  "reference_no": "PO-2026-001",
  "reason": "Stock receipt from supplier",
  "notes": "Inspected and verified",
  "batch_number": "BATCH-001"
}
```

**Movement Types**:
- `receipt`: Receiving goods from supplier
- `issue`: Issuing goods to customer/production
- `transfer`: Moving stock between locations
- `adjustment`: Manual stock adjustment
- `return`: Goods returned from customer

**Response**:
```json
{
  "status": "success",
  "message": "Stock movement created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439014",
    "movement_no": "MOV-2026-001",
    "movement_type": "receipt",
    "quantity": 50,
    "created_at": "2026-01-01T10:00:00Z"
  }
}
```

### 2. List Stock Movements

Get all stock movements for an organization with filters.

**Endpoint**: `GET /organizations/:org_id/stock-movements`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Query Parameters**:
- `movement_type` (string, optional): Filter by movement type
- `location_id` (string, optional): Filter by location
- `product_id` (string, optional): Filter by product
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20)

**Response**:
```json
{
  "status": "success",
  "message": "Stock movements retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439014",
      "movement_no": "MOV-2026-001",
      "product_id": "507f1f77bcf86cd799439011",
      "movement_type": "receipt",
      "quantity": 50,
      "movement_date": "2026-01-01T10:00:00Z"
    }
  ]
}
```

### 3. Get Stock Movement

Get details of a specific stock movement.

**Endpoint**: `GET /stock-movements/:id`

**URL Parameters**:
- `id` (string, required): Stock movement ID

**Response**:
```json
{
  "status": "success",
  "message": "Stock movement retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439014",
    "movement_no": "MOV-2026-001",
    "product_id": "507f1f77bcf86cd799439011",
    "movement_type": "receipt",
    "from_location_id": "507f1f77bcf86cd799439012",
    "to_location_id": "507f1f77bcf86cd799439013",
    "quantity": 50,
    "unit_cost": 25.50,
    "total_cost": 1275.00,
    "reference_type": "purchase_order",
    "reference_no": "PO-2026-001",
    "reason": "Stock receipt from supplier",
    "movement_date": "2026-01-01T10:00:00Z"
  }
}
```

### 4. Get Movements by Product

Get all stock movements for a specific product.

**Endpoint**: `GET /products/:product_id/movements`

**URL Parameters**:
- `product_id` (string, required): Product ID

**Query Parameters**:
- `page` (integer, optional): Page number
- `limit` (integer, optional): Items per page

**Response**: Same as List Stock Movements

### 5. Get Movements by Location

Get all stock movements involving a specific location.

**Endpoint**: `GET /locations/:location_id/movements`

**URL Parameters**:
- `location_id` (string, required): Location ID

**Response**: Same as List Stock Movements

---

## Batches

Batch management for lot tracking and expiry date control.

### 1. Create Batch

Create a new batch/lot for a product.

**Endpoint**: `POST /organizations/:org_id/batches`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "location_id": "507f1f77bcf86cd799439012",
  "batch_number": "BATCH-2026-001",
  "initial_quantity": 100,
  "unit_cost": 25.50,
  "manufacture_date": "2025-12-01T00:00:00Z",
  "expiry_date": "2027-12-01T00:00:00Z",
  "supplier_id": "507f1f77bcf86cd799439015",
  "purchase_order_id": "507f1f77bcf86cd799439016",
  "qc_status": "passed",
  "qc_notes": "Quality check passed"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Batch created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439017",
    "batch_number": "BATCH-2026-001",
    "initial_quantity": 100,
    "current_quantity": 100,
    "manufacture_date": "2025-12-01T00:00:00Z",
    "expiry_date": "2027-12-01T00:00:00Z"
  }
}
```

### 2. Get Batch

Get details of a specific batch.

**Endpoint**: `GET /batches/:id`

**URL Parameters**:
- `id` (string, required): Batch ID

**Response**:
```json
{
  "status": "success",
  "message": "Batch retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439017",
    "batch_number": "BATCH-2026-001",
    "product_id": "507f1f77bcf86cd799439011",
    "location_id": "507f1f77bcf86cd799439012",
    "initial_quantity": 100,
    "current_quantity": 85,
    "unit_cost": 25.50,
    "manufacture_date": "2025-12-01T00:00:00Z",
    "expiry_date": "2027-12-01T00:00:00Z",
    "is_active": true
  }
}
```

### 3. Update Batch Quantity

Adjust the quantity of a batch (increase or decrease).

**Endpoint**: `PUT /batches/:id/quantity`

**URL Parameters**:
- `id` (string, required): Batch ID

**Request Body**:
```json
{
  "delta": -15
}
```

**Note**: Use positive values to increase, negative to decrease.

**Response**:
```json
{
  "status": "success",
  "message": "Batch quantity updated successfully",
  "data": null
}
```

### 4. Get Batches by Product

Get all batches for a specific product.

**Endpoint**: `GET /products/:product_id/batches`

**URL Parameters**:
- `product_id` (string, required): Product ID

**Query Parameters**:
- `location_id` (string, optional): Filter by location
- `active_only` (boolean, optional): Show only active batches

**Response**:
```json
{
  "status": "success",
  "message": "Batches retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439017",
      "batch_number": "BATCH-2026-001",
      "current_quantity": 85,
      "expiry_date": "2027-12-01T00:00:00Z",
      "is_active": true
    }
  ]
}
```

---

## Stock Adjustments

Stock adjustments handle manual inventory corrections with approval workflow.

### 1. Create Stock Adjustment

Create a new stock adjustment request.

**Endpoint**: `POST /organizations/:org_id/stock-adjustments`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Request Body**:
```json
{
  "location_id": "507f1f77bcf86cd799439012",
  "adjustment_no": "ADJ-2026-001",
  "adjustment_date": "2026-01-01T10:00:00Z",
  "reason": "physical_count",
  "reason_details": "Variance found during physical count",
  "notes": "Investigated discrepancy",
  "items": [
    {
      "product_id": "507f1f77bcf86cd799439011",
      "expected_qty": 100,
      "actual_qty": 95,
      "uom": "EA",
      "unit_cost": 25.50,
      "batch_id": "507f1f77bcf86cd799439017",
      "reason": "Missing items"
    }
  ]
}
```

**Adjustment Reasons**:
- `physical_count`: Discrepancy from physical count
- `damage`: Damaged goods
- `theft`: Stolen items
- `obsolete`: Obsolete inventory
- `data_entry_error`: Correction of data entry error
- `other`: Other reasons

**Response**:
```json
{
  "status": "success",
  "message": "Stock adjustment created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439018",
    "adjustment_no": "ADJ-2026-001",
    "status": "draft",
    "total_items": 1,
    "total_variance": -5,
    "created_at": "2026-01-01T10:00:00Z"
  }
}
```

### 2. List Stock Adjustments

Get all stock adjustments for an organization.

**Endpoint**: `GET /organizations/:org_id/stock-adjustments`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Query Parameters**:
- `status` (string, optional): Filter by status (draft, pending, approved, rejected)
- `location_id` (string, optional): Filter by location
- `page` (integer, optional): Page number
- `limit` (integer, optional): Items per page

**Response**:
```json
{
  "status": "success",
  "message": "Adjustments retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439018",
      "adjustment_no": "ADJ-2026-001",
      "location_id": "507f1f77bcf86cd799439012",
      "status": "draft",
      "adjustment_date": "2026-01-01T10:00:00Z",
      "total_items": 1
    }
  ]
}
```

### 3. Get Stock Adjustment

Get details of a specific stock adjustment.

**Endpoint**: `GET /stock-adjustments/:id`

**URL Parameters**:
- `id` (string, required): Adjustment ID

**Response**:
```json
{
  "status": "success",
  "message": "Adjustment retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439018",
    "adjustment_no": "ADJ-2026-001",
    "location_id": "507f1f77bcf86cd799439012",
    "adjustment_date": "2026-01-01T10:00:00Z",
    "reason": "physical_count",
    "reason_details": "Variance found during physical count",
    "status": "draft",
    "items": [
      {
        "product_id": "507f1f77bcf86cd799439011",
        "expected_qty": 100,
        "actual_qty": 95,
        "difference_qty": -5,
        "unit_cost": 25.50,
        "total_cost": -127.50
      }
    ]
  }
}
```

### 4. Update Stock Adjustment

Update a draft stock adjustment.

**Endpoint**: `PUT /stock-adjustments/:id`

**URL Parameters**:
- `id` (string, required): Adjustment ID

**Request Body**: Same as Create Stock Adjustment

**Note**: Can only update adjustments with status "draft"

**Response**:
```json
{
  "status": "success",
  "message": "Adjustment updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439018",
    "adjustment_no": "ADJ-2026-001",
    "status": "draft",
    "updated_at": "2026-01-01T11:00:00Z"
  }
}
```

### 5. Approve Stock Adjustment

Approve an adjustment and apply stock changes.

**Endpoint**: `POST /stock-adjustments/:id/approve`

**URL Parameters**:
- `id` (string, required): Adjustment ID

**Request Body**:
```json
{
  "notes": "Approved after verification"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Adjustment approved successfully",
  "data": null
}
```

### 6. Reject Stock Adjustment

Reject a stock adjustment request.

**Endpoint**: `POST /stock-adjustments/:id/reject`

**URL Parameters**:
- `id` (string, required): Adjustment ID

**Request Body**:
```json
{
  "reason": "Insufficient justification",
  "notes": "Please provide more details"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Adjustment rejected successfully",
  "data": null
}
```

### 7. Delete Stock Adjustment

Delete a draft stock adjustment.

**Endpoint**: `DELETE /stock-adjustments/:id`

**URL Parameters**:
- `id` (string, required): Adjustment ID

**Note**: Can only delete adjustments with status "draft"

**Response**:
```json
{
  "status": "success",
  "message": "Adjustment deleted successfully",
  "data": null
}
```

---

## Inventory Counts

Physical inventory count operations with variance tracking.

### 1. Create Inventory Count

Start a new physical inventory count.

**Endpoint**: `POST /organizations/:org_id/inventory-counts`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Request Body**:
```json
{
  "location_id": "507f1f77bcf86cd799439012",
  "count_no": "COUNT-2026-001",
  "count_date": "2026-01-01T10:00:00Z",
  "count_type": "full",
  "notes": "Year-end inventory count",
  "items": [
    {
      "product_id": "507f1f77bcf86cd799439011",
      "system_qty": 100,
      "uom": "EA",
      "unit_cost": 25.50,
      "batch_id": "507f1f77bcf86cd799439017"
    }
  ]
}
```

**Count Types**:
- `full`: Complete inventory count
- `cycle`: Cycle count (partial inventory)
- `spot`: Spot check specific items

**Response**:
```json
{
  "status": "success",
  "message": "Inventory count created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439019",
    "count_no": "COUNT-2026-001",
    "status": "draft",
    "count_type": "full",
    "total_items": 1,
    "created_at": "2026-01-01T10:00:00Z"
  }
}
```

### 2. List Inventory Counts

Get all inventory counts for an organization.

**Endpoint**: `GET /organizations/:org_id/inventory-counts`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Query Parameters**:
- `status` (string, optional): Filter by status
- `location_id` (string, optional): Filter by location
- `count_type` (string, optional): Filter by count type
- `page` (integer, optional): Page number
- `limit` (integer, optional): Items per page

**Response**:
```json
{
  "status": "success",
  "message": "Inventory counts retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439019",
      "count_no": "COUNT-2026-001",
      "location_id": "507f1f77bcf86cd799439012",
      "count_type": "full",
      "status": "draft",
      "count_date": "2026-01-01T10:00:00Z"
    }
  ]
}
```

### 3. Get Inventory Count

Get details of a specific inventory count.

**Endpoint**: `GET /inventory-counts/:id`

**URL Parameters**:
- `id` (string, required): Inventory count ID

**Response**:
```json
{
  "status": "success",
  "message": "Inventory count retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439019",
    "count_no": "COUNT-2026-001",
    "location_id": "507f1f77bcf86cd799439012",
    "count_date": "2026-01-01T10:00:00Z",
    "count_type": "full",
    "status": "draft",
    "items": [
      {
        "product_id": "507f1f77bcf86cd799439011",
        "system_qty": 100,
        "counted_qty": 0,
        "variance_qty": 0,
        "uom": "EA",
        "unit_cost": 25.50
      }
    ]
  }
}
```

### 4. Update Count Item

Update the counted quantity for a specific item.

**Endpoint**: `POST /inventory-counts/:id/items`

**URL Parameters**:
- `id` (string, required): Inventory count ID

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "counted_qty": 95,
  "notes": "Physical count confirmed",
  "batch_id": "507f1f77bcf86cd799439017"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Count item updated successfully",
  "data": null
}
```

### 5. Complete Inventory Count

Finalize the inventory count.

**Endpoint**: `POST /inventory-counts/:id/complete`

**URL Parameters**:
- `id` (string, required): Inventory count ID

**Request Body**:
```json
{
  "create_adjustment": true
}
```

**Note**: If `create_adjustment` is true, system will automatically create a stock adjustment for any variances found.

**Response**:
```json
{
  "status": "success",
  "message": "Count completed successfully",
  "data": null
}
```

### 6. Cancel Inventory Count

Cancel an in-progress inventory count.

**Endpoint**: `POST /inventory-counts/:id/cancel`

**URL Parameters**:
- `id` (string, required): Inventory count ID

**Response**:
```json
{
  "status": "success",
  "message": "Count cancelled successfully",
  "data": null
}
```

### 7. Delete Inventory Count

Delete a draft or cancelled inventory count.

**Endpoint**: `DELETE /inventory-counts/:id`

**URL Parameters**:
- `id` (string, required): Inventory count ID

**Note**: Can only delete counts with status "draft" or "cancelled"

**Response**:
```json
{
  "status": "success",
  "message": "Count deleted successfully",
  "data": null
}
```

---

## Serial Numbers

Serial number tracking for individual item management.

### 1. Create Serial Number

Register a new serial number.

**Endpoint**: `POST /organizations/:org_id/serial-numbers`

**URL Parameters**:
- `org_id` (string, required): Organization ID

**Request Body**:
```json
{
  "product_id": "507f1f77bcf86cd799439011",
  "location_id": "507f1f77bcf86cd799439012",
  "serial_no": "SN-2026-001",
  "manufacture_date": "2025-12-01T00:00:00Z",
  "warranty_expiry": "2027-12-01T00:00:00Z",
  "batch_id": "507f1f77bcf86cd799439017",
  "unit_cost": 25.50
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Serial number created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "serial_no": "SN-2026-001",
    "product_id": "507f1f77bcf86cd799439011",
    "status": "available",
    "is_available": true,
    "created_at": "2026-01-01T10:00:00Z"
  }
}
```

### 2. Get Serial Number

Get details of a specific serial number by ID.

**Endpoint**: `GET /serial-numbers/:id`

**URL Parameters**:
- `id` (string, required): Serial number ID

**Response**:
```json
{
  "status": "success",
  "message": "Serial number retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "serial_no": "SN-2026-001",
    "product_id": "507f1f77bcf86cd799439011",
    "location_id": "507f1f77bcf86cd799439012",
    "status": "available",
    "is_available": true,
    "manufacture_date": "2025-12-01T00:00:00Z",
    "warranty_expiry": "2027-12-01T00:00:00Z",
    "unit_cost": 25.50
  }
}
```

### 3. Get Serial Number by Serial String

Look up a serial number by its serial string.

**Endpoint**: `GET /organizations/:org_id/serial-numbers/:serial_no`

**URL Parameters**:
- `org_id` (string, required): Organization ID
- `serial_no` (string, required): Serial number string

**Response**: Same as Get Serial Number

### 4. List Serial Numbers by Product

Get all serial numbers for a specific product.

**Endpoint**: `GET /products/:product_id/serial-numbers`

**URL Parameters**:
- `product_id` (string, required): Product ID

**Query Parameters**:
- `location_id` (string, optional): Filter by location
- `available` (boolean, optional): Filter by availability (true/false)

**Response**:
```json
{
  "status": "success",
  "message": "Serial numbers retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "serial_no": "SN-2026-001",
      "status": "available",
      "location_id": "507f1f77bcf86cd799439012",
      "manufacture_date": "2025-12-01T00:00:00Z"
    }
  ]
}
```

### 5. Update Serial Number

Update serial number information.

**Endpoint**: `PUT /serial-numbers/:id`

**URL Parameters**:
- `id` (string, required): Serial number ID

**Request Body**:
```json
{
  "location_id": "507f1f77bcf86cd799439013",
  "warranty_expiry": "2028-12-01T00:00:00Z",
  "notes": "Warranty extended"
}
```

**Note**: Cannot update sold serial numbers.

**Response**:
```json
{
  "status": "success",
  "message": "Serial number updated successfully",
  "data": null
}
```

### 6. Allocate Serial Number

Allocate a serial number to a customer/sales order.

**Endpoint**: `POST /serial-numbers/:id/allocate`

**URL Parameters**:
- `id` (string, required): Serial number ID

**Request Body**:
```json
{
  "customer_id": "507f1f77bcf86cd799439021",
  "sales_order_id": "507f1f77bcf86cd799439022"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Serial number allocated successfully",
  "data": null
}
```

### 7. Mark Serial Number as Sold

Mark a serial number as sold.

**Endpoint**: `POST /serial-numbers/:id/sold`

**URL Parameters**:
- `id` (string, required): Serial number ID

**Response**:
```json
{
  "status": "success",
  "message": "Serial number marked as sold successfully",
  "data": null
}
```

### 8. Delete Serial Number

Soft-delete a serial number.

**Endpoint**: `DELETE /serial-numbers/:id`

**URL Parameters**:
- `id` (string, required): Serial number ID

**Note**: Cannot delete sold serial numbers.

**Response**:
```json
{
  "status": "success",
  "message": "Serial number deleted successfully",
  "data": null
}
```

---

## Error Responses

All endpoints return consistent error responses:

```json
{
  "status": "error",
  "error_code": "ERROR_CODE",
  "message": "Human readable error message",
  "details": {}
}
```

### Common Error Codes:

- `INVALID_ID`: Invalid ID format
- `NOT_FOUND`: Resource not found
- `VALIDATION_ERROR`: Request validation failed
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Insufficient permissions
- `CREATE_FAILED`: Failed to create resource
- `UPDATE_FAILED`: Failed to update resource
- `DELETE_FAILED`: Failed to delete resource
- `INSUFFICIENT_STOCK`: Not enough stock available
- `ALREADY_EXISTS`: Resource already exists
- `INVALID_STATUS`: Operation not allowed for current status

### HTTP Status Codes:

- `200 OK`: Successful GET/PUT/DELETE request
- `201 Created`: Successful POST request
- `400 Bad Request`: Validation error or bad input
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## Pagination

List endpoints support pagination with the following query parameters:

- `page` (integer): Page number, starting from 1 (default: 1)
- `limit` (integer): Number of items per page (default: 20, max: 100)

Paginated responses include metadata:

```json
{
  "status": "success",
  "message": "Items retrieved successfully",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

## Best Practices

1. **Stock Movements**: Always create stock movements for all inventory transactions to maintain an accurate audit trail.

2. **Batch Tracking**: Use batches for products with expiry dates or that require lot traceability.

3. **Serial Numbers**: Implement serial number tracking for high-value items or items requiring warranty management.

4. **Stock Adjustments**: Always provide detailed reasons for stock adjustments and use the approval workflow.

5. **Inventory Counts**: Perform regular cycle counts to maintain inventory accuracy. Set `create_adjustment: true` to automatically adjust variances.

6. **Allocations**: Always allocate stock when creating sales orders to prevent overselling.

7. **Error Handling**: Implement proper error handling for all API calls and check stock availability before transactions.

8. **Idempotency**: For critical operations, implement retry logic with idempotency keys to prevent duplicate transactions.

---

## Support

For questions or issues, please contact the development team or refer to the main ERP system documentation.
