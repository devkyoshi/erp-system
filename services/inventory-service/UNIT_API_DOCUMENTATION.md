# Unit and Unit Chart API Endpoints

This document describes the newly created API endpoints for managing Units and Unit Charts (conversion rules) in the inventory service.

## Authentication

All endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Units API

Base URL: `/api/v1/units`

### 1. Create Unit
Create a new unit of measure.

**Endpoint:** `POST /api/v1/units`

**Request Body:**
```json
{
  "code": "PCS",
  "name": "Pieces",
  "symbol": "pcs",
  "description": "Individual pieces",
  "unit_type": "quantity",
  "is_base_unit": true,
  "metadata": {}
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "code": "PCS",
    "name": "Pieces",
    "symbol": "pcs",
    "description": "Individual pieces",
    "unit_type": "quantity",
    "is_base_unit": true,
    "is_active": true,
    "metadata": {},
    "created_at": "2026-01-11T10:00:00Z",
    "updated_at": "2026-01-11T10:00:00Z",
    "version": 1
  },
  "message": "Unit created successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 2. Get Unit by ID
Retrieve a specific unit by its ID.

**Endpoint:** `GET /api/v1/units/:id`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "code": "PCS",
    "name": "Pieces",
    ...
  },
  "message": "Unit retrieved successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 3. Get All Units
Retrieve all units with optional filters.

**Endpoint:** `GET /api/v1/units`

**Query Parameters:**
- `unit_type` (optional): Filter by unit type (e.g., "quantity", "weight", "volume", "length")
- `active_only` (optional): If "true", only return active units

**Examples:**
- `GET /api/v1/units` - Get all units
- `GET /api/v1/units?unit_type=quantity` - Get all quantity units
- `GET /api/v1/units?active_only=true` - Get only active units

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "code": "PCS",
      "name": "Pieces",
      ...
    },
    {
      "_id": "507f1f77bcf86cd799439012",
      "code": "BOX",
      "name": "Box",
      ...
    }
  ],
  "message": "Units retrieved successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 4. Update Unit
Update an existing unit.

**Endpoint:** `PUT /api/v1/units/:id`

**Request Body:**
```json
{
  "name": "Updated Pieces",
  "symbol": "pc",
  "description": "Updated description",
  "is_active": true,
  "metadata": {}
}
```

**Note:** All fields are optional. Only provide fields you want to update.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "code": "PCS",
    "name": "Updated Pieces",
    ...
  },
  "message": "Unit updated successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 5. Delete Unit
Soft delete a unit (marks as inactive).

**Endpoint:** `DELETE /api/v1/units/:id`

**Note:** Base units and units with active conversions cannot be deleted.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": null,
  "message": "Unit deleted successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

---

## Unit Charts API (Conversion Rules)

Base URL: `/api/v1/unit-charts`

### 1. Create Unit Chart
Create a new unit conversion rule.

**Endpoint:** `POST /api/v1/unit-charts`

**Request Body:**
```json
{
  "from_unit_id": "507f1f77bcf86cd799439011",
  "to_unit_id": "507f1f77bcf86cd799439012",
  "conversion_rate": 12.0,
  "metadata": {}
}
```

**Example:** If 1 BOX = 12 PCS, then:
- `from_unit_id`: BOX unit ID
- `to_unit_id`: PCS unit ID
- `conversion_rate`: 12.0

**Validation Rules:**
- Both units must exist
- Both units must be of the same type (e.g., both quantity units)
- Cannot create conversion from a unit to itself
- Cannot create duplicate conversions
- Circular conversion paths are not allowed

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439013",
    "from_unit_id": "507f1f77bcf86cd799439011",
    "to_unit_id": "507f1f77bcf86cd799439012",
    "conversion_rate": 12.0,
    "is_active": true,
    "metadata": {},
    "created_at": "2026-01-11T10:00:00Z",
    "updated_at": "2026-01-11T10:00:00Z",
    "version": 1
  },
  "message": "Unit chart created successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 2. Get Unit Chart by ID
Retrieve a specific unit chart by its ID.

**Endpoint:** `GET /api/v1/unit-charts/:id`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439013",
    "from_unit_id": "507f1f77bcf86cd799439011",
    "to_unit_id": "507f1f77bcf86cd799439012",
    "conversion_rate": 12.0,
    ...
  },
  "message": "Unit chart retrieved successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 3. Get All Unit Charts
Retrieve all unit charts with optional filters.

**Endpoint:** `GET /api/v1/unit-charts`

**Query Parameters:**
- `active_only` (optional): If "true", only return active unit charts

**Examples:**
- `GET /api/v1/unit-charts` - Get all unit charts
- `GET /api/v1/unit-charts?active_only=true` - Get only active unit charts

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439013",
      "from_unit_id": "507f1f77bcf86cd799439011",
      "to_unit_id": "507f1f77bcf86cd799439012",
      "conversion_rate": 12.0,
      ...
    }
  ],
  "message": "Unit charts retrieved successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 4. Update Unit Chart
Update an existing unit chart.

**Endpoint:** `PUT /api/v1/unit-charts/:id`

**Request Body:**
```json
{
  "conversion_rate": 15.0,
  "is_active": true,
  "metadata": {}
}
```

**Note:** All fields are optional. Only provide fields you want to update. You cannot change `from_unit_id` or `to_unit_id`.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439013",
    "from_unit_id": "507f1f77bcf86cd799439011",
    "to_unit_id": "507f1f77bcf86cd799439012",
    "conversion_rate": 15.0,
    ...
  },
  "message": "Unit chart updated successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 5. Delete Unit Chart
Soft delete a unit chart (marks as inactive).

**Endpoint:** `DELETE /api/v1/unit-charts/:id`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": null,
  "message": "Unit chart deleted successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### 6. Get Conversion Rate
Get the conversion rate between two units.

**Endpoint:** `GET /api/v1/unit-charts/conversion-rate`

**Query Parameters:**
- `from_unit_id` (required): Source unit ID
- `to_unit_id` (required): Target unit ID

**Example:**
```
GET /api/v1/unit-charts/conversion-rate?from_unit_id=507f1f77bcf86cd799439011&to_unit_id=507f1f77bcf86cd799439012
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "from_unit_id": "507f1f77bcf86cd799439011",
    "to_unit_id": "507f1f77bcf86cd799439012",
    "conversion_rate": 12.0
  },
  "message": "Conversion rate retrieved successfully",
  "timestamp": "2026-01-11T10:00:00Z"
}
```

**Note:** If the units are the same, the conversion rate is 1.0.

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": []
  },
  "timestamp": "2026-01-11T10:00:00Z"
}
```

### Common Error Codes:
- `VALIDATION_ERROR` - Request validation failed
- `UNAUTHORIZED` - Authentication required or failed
- `INVALID_ID` - Invalid ObjectID format
- `INVALID_USER_ID` - Invalid user ID
- `NOT_FOUND` - Resource not found
- `CREATE_FAILED` - Failed to create resource
- `UPDATE_FAILED` - Failed to update resource
- `DELETE_FAILED` - Failed to delete resource
- `FETCH_FAILED` - Failed to fetch resources
- `MISSING_PARAMS` - Required parameters missing

### Common HTTP Status Codes:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Database Indexes

The following indexes have been created for optimal performance:

### Units Collection:
- Unique index on `code`
- Index on `unit_type` + `is_active`
- Index on `is_base_unit` + `unit_type`

### Unit Charts Collection:
- Unique index on `from_unit_id` + `to_unit_id`
- Index on `from_unit_id` + `is_active`
- Index on `to_unit_id` + `is_active`

---

## Usage Examples

### Example 1: Create a Quantity System
```bash
# 1. Create base unit (Pieces)
curl -X POST http://localhost:8082/api/v1/units \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PCS",
    "name": "Pieces",
    "symbol": "pcs",
    "unit_type": "quantity",
    "is_base_unit": true
  }'

# 2. Create derived unit (Box)
curl -X POST http://localhost:8082/api/v1/units \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BOX",
    "name": "Box",
    "symbol": "box",
    "unit_type": "quantity",
    "is_base_unit": false
  }'

# 3. Create conversion (1 BOX = 12 PCS)
curl -X POST http://localhost:8082/api/v1/unit-charts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "from_unit_id": "<box_unit_id>",
    "to_unit_id": "<pcs_unit_id>",
    "conversion_rate": 12.0
  }'
```

### Example 2: Get Conversion Rate
```bash
curl -X GET "http://localhost:8082/api/v1/unit-charts/conversion-rate?from_unit_id=<box_id>&to_unit_id=<pcs_id>" \
  -H "Authorization: Bearer <token>"
```

### Example 3: List All Quantity Units
```bash
curl -X GET "http://localhost:8082/api/v1/units?unit_type=quantity&active_only=true" \
  -H "Authorization: Bearer <token>"
```
