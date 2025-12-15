# Organization Service API Documentation

## Overview

The Organization Service manages organizations, companies, and locations in a multi-tenant ERP system. It provides a hierarchical structure where:
- **Organizations** are the top-level tenants
- **Companies** belong to organizations (e.g., subsidiaries, business units)
- **Locations** belong to companies (e.g., warehouses, stores, offices)

**Base URL:** `http://localhost:<PORT>/api/v1`

**Authentication:** All endpoints require JWT authentication via the `Authorization: Bearer <token>` header.

---

## Table of Contents

1. [Organizations](#organizations)
2. [Companies](#companies)
3. [Locations](#locations)
4. [Error Responses](#error-responses)
5. [Data Models](#data-models)

---

## Organizations

### 1. Create Organization

Creates a new organization with default settings.

**Endpoint:** `POST /organizations`

**Request Body:**
```json
{
  "name": "Acme Corporation",
  "legal_name": "Acme Corporation Ltd.",
  "domain": "acme.com",
  "email": "contact@acme.com",
  "phone": "+1-555-0100",
  "website": "https://www.acme.com",
  "tax_id": "12-3456789",
  "registration_number": "REG-2023-001",
  "industry": "Technology",
  "company_size": "enterprise",
  "billing_email": "billing@acme.com",
  "max_users": 500,
  "max_companies": 10,
  "max_locations": 50,
  "storage_limit_gb": 100.0
}
```

**Response:** `201 Created`
```json
{
  "status": "success",
  "code": 201,
  "message": "Organization created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corporation",
    "legal_name": "Acme Corporation Ltd.",
    "domain": "acme.com",
    "alternate_domains": [],
    "logo": "",
    "favicon": "",
    "email": "contact@acme.com",
    "phone": "+1-555-0100",
    "website": "https://www.acme.com",
    "tax_id": "12-3456789",
    "registration_number": "REG-2023-001",
    "industry": "Technology",
    "company_size": "enterprise",
    "status": "active",
    "is_active": true,
    "activated_at": "2025-12-15T10:00:00Z",
    "suspended_at": null,
    "suspension_reason": "",
    "subscription_id": null,
    "billing_email": "billing@acme.com",
    "billing_address": null,
    "max_users": 500,
    "max_companies": 10,
    "max_locations": 50,
    "current_users": 0,
    "current_companies": 0,
    "current_locations": 0,
    "storage_used_gb": 0,
    "storage_limit_gb": 100.0,
    "settings": {
      "timezone": "UTC",
      "date_format": "YYYY-MM-DD",
      "time_format": "HH:mm:ss",
      "currency": "USD",
      "language": "en",
      "enabled_modules": [],
      "allow_user_registration": false,
      "require_email_verification": true,
      "enable_mfa": false,
      "session_timeout": 30,
      "password_policy": {
        "min_length": 8,
        "require_uppercase": true,
        "require_lowercase": true,
        "require_numbers": true,
        "require_special_chars": false,
        "expiry_days": 90,
        "prevent_reuse_count": 5
      }
    },
    "brand_colors": null,
    "metadata": {},
    "tags": [],
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "created_by": "507f1f77bcf86cd799439012",
    "updated_by": null,
    "deleted_at": null,
    "deleted_by": null
  }
}
```

**Validation Rules:**
- `name`: Required
- `legal_name`: Required
- `domain`: Required, must be unique
- `email`: Required, must be valid email
- `billing_email`: Required, must be valid email

---

### 2. Get Organization

Retrieves a single organization by ID.

**Endpoint:** `GET /organizations/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the organization

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Organization retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corporation",
    "legal_name": "Acme Corporation Ltd.",
    // ... (same structure as create response)
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid organization ID format
- `404 Not Found`: Organization not found

---

### 3. List Organizations

Returns a paginated list of organizations with optional filtering by status.

**Endpoint:** `GET /organizations`

**Query Parameters:**
- `page` (integer, optional): Page number, default: 1
- `limit` (integer, optional): Items per page, default: 20, max: 100
- `status` (string, optional): Filter by status (`active`, `inactive`, `suspended`, `trial`, `cancelled`)

**Example Request:**
```
GET /organizations?page=1&limit=20&status=active
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Acme Corporation",
      // ... organization object
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "name": "Tech Innovations Inc.",
      // ... organization object
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

### 4. Update Organization

Updates an existing organization. All fields are optional - only provided fields will be updated.

**Endpoint:** `PUT /organizations/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the organization

**Request Body:**
```json
{
  "name": "Acme Corporation Updated",
  "legal_name": "Acme Corporation Ltd. Updated",
  "domain": "acme.com",
  "email": "newcontact@acme.com",
  "phone": "+1-555-0200",
  "website": "https://www.acme-new.com",
  "logo": "https://cdn.acme.com/logo.png",
  "is_active": true
}
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Organization updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corporation Updated",
    // ... updated organization object
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID or domain already exists
- `404 Not Found`: Organization not found

---

### 5. Delete Organization

Soft deletes an organization (sets `deleted_at` timestamp).

**Endpoint:** `DELETE /organizations/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the organization

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Organization deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID or delete operation failed
- `404 Not Found`: Organization not found

---

## Companies

### 1. Create Company

Creates a new company under an organization.

**Endpoint:** `POST /organizations/:id/companies`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the parent organization

**Request Body:**
```json
{
  "name": "Acme North America",
  "legal_name": "Acme North America Inc.",
  "code": "ACME-NA",
  "tax_id": "98-7654321",
  "registration_number": "REG-NA-2023-001",
  "vat_number": "VAT-123456",
  "email": "na@acme.com",
  "phone": "+1-555-0150",
  "address": {
    "street": "123 Main Street",
    "street2": "Suite 400",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "United States",
    "country_code": "US",
    "latitude": 40.7128,
    "longitude": -74.0060
  },
  "bank_accounts": [
    {
      "bank_name": "First National Bank",
      "account_number": "1234567890",
      "account_name": "Acme North America",
      "iban": "US12345678901234567890",
      "swift_code": "FNBUS33",
      "branch": "Manhattan Branch",
      "is_default": true
    }
  ],
  "is_default": true
}
```

**Response:** `201 Created`
```json
{
  "status": "success",
  "code": 201,
  "message": "Company created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Acme North America",
    "legal_name": "Acme North America Inc.",
    "code": "ACME-NA",
    "tax_id": "98-7654321",
    "registration_number": "REG-NA-2023-001",
    "vat_number": "VAT-123456",
    "email": "na@acme.com",
    "phone": "+1-555-0150",
    "fax": "",
    "website": "",
    "address": {
      "street": "123 Main Street",
      "street2": "Suite 400",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "United States",
      "country_code": "US",
      "latitude": 40.7128,
      "longitude": -74.0060
    },
    "bank_accounts": [
      {
        "bank_name": "First National Bank",
        "account_number": "1234567890",
        "account_name": "Acme North America",
        "iban": "US12345678901234567890",
        "swift_code": "FNBUS33",
        "branch": "Manhattan Branch",
        "is_default": true
      }
    ],
    "settings": {
      "fiscal_year_start": "01-01",
      "currency": "USD",
      "timezone": "UTC",
      "enable_multi_currency": false
    },
    "is_active": true,
    "is_default": true,
    "parent_company_id": null,
    "metadata": {},
    "tags": [],
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "created_by": "507f1f77bcf86cd799439012",
    "updated_by": null,
    "deleted_at": null,
    "deleted_by": null
  }
}
```

**Validation Rules:**
- `name`: Required
- `legal_name`: Required
- `code`: Required, must be unique within organization
- `email`: Required, must be valid email
- `address`: Required

**Error Responses:**
- `400 Bad Request`: Invalid organization ID, validation error, or company limit reached
- `404 Not Found`: Organization not found

---

### 2. Get Company

Retrieves a single company by ID.

**Endpoint:** `GET /companies/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the company

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Company retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Acme North America",
    // ... (same structure as create response)
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid company ID format
- `404 Not Found`: Company not found

---

### 3. List Companies

Returns a paginated list of companies under an organization.

**Endpoint:** `GET /organizations/:id/companies`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the parent organization

**Query Parameters:**
- `page` (integer, optional): Page number, default: 1
- `limit` (integer, optional): Items per page, default: 20, max: 100

**Example Request:**
```
GET /organizations/507f1f77bcf86cd799439011/companies?page=1&limit=20
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "organization_id": "507f1f77bcf86cd799439011",
      "name": "Acme North America",
      // ... company object
    },
    {
      "id": "507f1f77bcf86cd799439014",
      "organization_id": "507f1f77bcf86cd799439011",
      "name": "Acme Europe",
      // ... company object
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

---

### 4. Update Company

Updates an existing company. All fields are optional.

**Endpoint:** `PUT /companies/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the company

**Request Body:**
```json
{
  "name": "Acme North America Updated",
  "legal_name": "Acme North America Inc. Updated",
  "code": "ACME-NA-NEW",
  "email": "na-new@acme.com",
  "phone": "+1-555-0160",
  "is_active": true
}
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Company updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "name": "Acme North America Updated",
    // ... updated company object
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID, validation error, or code already exists
- `404 Not Found`: Company not found

---

### 5. Delete Company

Soft deletes a company.

**Endpoint:** `DELETE /companies/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the company

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Company deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID or delete operation failed
- `404 Not Found`: Company not found

---

## Locations

### 1. Create Location

Creates a new location under a company.

**Endpoint:** `POST /companies/:id/locations`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the parent company

**Request Body:**
```json
{
  "name": "Manhattan Warehouse",
  "code": "WH-NY-001",
  "type": "warehouse",
  "email": "warehouse-ny@acme.com",
  "phone": "+1-555-0170",
  "address": {
    "street": "456 Industrial Blvd",
    "street2": "",
    "city": "New York",
    "state": "NY",
    "postal_code": "10002",
    "country": "United States",
    "country_code": "US",
    "latitude": 40.7200,
    "longitude": -74.0100
  },
  "warehouse_info": {
    "total_area": 50000.0,
    "storage_capacity": 5000,
    "docking_bays": 12,
    "refrigerated_area": 10000.0,
    "has_cold_storage": true
  },
  "store_info": null,
  "is_default": true
}
```

**Response:** `201 Created`
```json
{
  "status": "success",
  "code": 201,
  "message": "Location created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "company_id": "507f1f77bcf86cd799439013",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Manhattan Warehouse",
    "code": "WH-NY-001",
    "type": "warehouse",
    "category": "",
    "email": "warehouse-ny@acme.com",
    "phone": "+1-555-0170",
    "fax": "",
    "address": {
      "street": "456 Industrial Blvd",
      "street2": "",
      "city": "New York",
      "state": "NY",
      "postal_code": "10002",
      "country": "United States",
      "country_code": "US",
      "latitude": 40.7200,
      "longitude": -74.0100
    },
    "manager_id": null,
    "applications": [],
    "warehouse_info": {
      "total_area": 50000.0,
      "storage_capacity": 5000,
      "docking_bays": 12,
      "refrigerated_area": 10000.0,
      "has_cold_storage": true
    },
    "store_info": null,
    "operating_hours": [],
    "settings": {
      "timezone": "UTC",
      "allow_backdated_transactions": false,
      "require_approval": false
    },
    "is_active": true,
    "is_default": true,
    "parent_location_id": null,
    "metadata": {},
    "tags": [],
    "created_at": "2025-12-15T10:00:00Z",
    "updated_at": "2025-12-15T10:00:00Z",
    "created_by": "507f1f77bcf86cd799439012",
    "updated_by": null,
    "deleted_at": null,
    "deleted_by": null
  }
}
```

**Location Types:**
- `head_office`: Headquarters
- `branch`: Branch office
- `warehouse`: Warehouse/distribution center
- `store`: Retail store
- `factory`: Manufacturing facility
- `distribution_center`: Distribution center

**Validation Rules:**
- `name`: Required
- `code`: Required, must be unique within company
- `type`: Required, must be valid LocationType
- `email`: Required, must be valid email
- `address`: Required

**Type-Specific Information:**
- When `type` is `warehouse`, include `warehouse_info`
- When `type` is `store`, include `store_info`

**Error Responses:**
- `400 Bad Request`: Invalid company ID, validation error, or location limit reached
- `404 Not Found`: Company not found

---

### 2. Get Location

Retrieves a single location by ID.

**Endpoint:** `GET /locations/:id`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the location

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Location retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "company_id": "507f1f77bcf86cd799439013",
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Manhattan Warehouse",
    // ... (same structure as create response)
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid location ID format
- `404 Not Found`: Location not found

---

### 3. List Locations

Returns a paginated list of locations under a company.

**Endpoint:** `GET /companies/:id/locations`

**Path Parameters:**
- `id` (string): MongoDB ObjectID of the parent company

**Query Parameters:**
- `page` (integer, optional): Page number, default: 1
- `limit` (integer, optional): Items per page, default: 20, max: 100

**Example Request:**
```
GET /companies/507f1f77bcf86cd799439013/locations?page=1&limit=20
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "company_id": "507f1f77bcf86cd799439013",
      "name": "Manhattan Warehouse",
      "type": "warehouse",
      // ... location object
    },
    {
      "id": "507f1f77bcf86cd799439016",
      "company_id": "507f1f77bcf86cd799439013",
      "name": "Brooklyn Store",
      "type": "store",
      // ... location object
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

## Error Responses

All error responses follow this structure:

```json
{
  "status": "error",
  "code": 400,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "must be a valid email address"
    }
  }
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | VALIDATION_ERROR | Request validation failed |
| 400 | INVALID_ID | Invalid MongoDB ObjectID format |
| 400 | CREATE_FAILED | Resource creation failed |
| 400 | UPDATE_FAILED | Resource update failed |
| 400 | DELETE_FAILED | Resource deletion failed |
| 400 | LIST_FAILED | Resource listing failed |
| 401 | UNAUTHORIZED | Missing or invalid authentication token |
| 404 | NOT_FOUND | Resource not found |
| 500 | INTERNAL_ERROR | Internal server error |

---

## Data Models

### Organization Status

```typescript
type OrganizationStatus = 
  | "active"      // Active and operational
  | "inactive"    // Inactive
  | "suspended"   // Temporarily suspended
  | "trial"       // Trial period
  | "cancelled"   // Cancelled subscription
```

### Company Size

```typescript
type CompanySize = 
  | "small"       // 1-50 employees
  | "medium"      // 51-250 employees
  | "large"       // 251-1000 employees
  | "enterprise"  // 1000+ employees
```

### Location Type

```typescript
type LocationType = 
  | "head_office"          // Headquarters
  | "branch"               // Branch office
  | "warehouse"            // Warehouse
  | "store"                // Retail store
  | "factory"              // Manufacturing facility
  | "distribution_center"  // Distribution center
```

### Address Object

```json
{
  "street": "123 Main Street",
  "street2": "Suite 400",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "country": "United States",
  "country_code": "US",
  "latitude": 40.7128,
  "longitude": -74.0060
}
```

### Bank Account Object

```json
{
  "bank_name": "First National Bank",
  "account_number": "1234567890",
  "account_name": "Acme Corporation",
  "iban": "US12345678901234567890",
  "swift_code": "FNBUS33",
  "branch": "Manhattan Branch",
  "is_default": true
}
```

### Warehouse Info Object

```json
{
  "total_area": 50000.0,
  "storage_capacity": 5000,
  "docking_bays": 12,
  "refrigerated_area": 10000.0,
  "has_cold_storage": true
}
```

### Store Info Object

```json
{
  "floor_area": 5000.0,
  "pos_count": 5,
  "parking_spaces": 50,
  "has_online_ordering": true
}
```

### Operating Hours Object

```json
{
  "day_of_week": 1,
  "open_time": "09:00",
  "close_time": "18:00",
  "is_closed": false
}
```

**Day of Week Values:**
- 0: Sunday
- 1: Monday
- 2: Tuesday
- 3: Wednesday
- 4: Thursday
- 5: Friday
- 6: Saturday

---

## Authentication

All API endpoints require JWT authentication. Include the JWT token in the request header:

```
Authorization: Bearer <your-jwt-token>
```

The JWT token should be obtained from the Authentication Service. The token contains the user ID which is used to track who created, updated, or deleted resources.

---

## Rate Limiting

The Organization Service implements rate limiting:
- **Default:** 100 requests per minute per IP address
- Rate limit information is returned in response headers:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time when the rate limit resets

---

## Pagination

List endpoints support pagination with the following query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

Pagination information is included in the response:

```json
{
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

## Health Check

**Endpoint:** `GET /health`

**Authentication:** Not required

**Response:** `200 OK`
```json
{
  "status": "success",
  "code": 200,
  "message": "Service is healthy",
  "data": {
    "status": "healthy"
  }
}
```

---

## Best Practices

1. **Error Handling**: Always check the `status` field in responses. On error, examine the `error.code` and `error.message` for details.

2. **Pagination**: Use pagination for list endpoints to avoid loading too much data at once.

3. **Filtering**: Use query parameters like `status` to filter results and reduce response size.

4. **Validation**: Validate data on the client side before sending requests to reduce unnecessary API calls.

5. **Idempotency**: Update and delete operations are idempotent - calling them multiple times with the same data produces the same result.

6. **Soft Deletes**: All delete operations are soft deletes (set `deleted_at` timestamp). Data is not physically removed from the database.

7. **Hierarchical Queries**: 
   - Companies belong to organizations
   - Locations belong to companies
   - Always ensure proper hierarchy when creating resources

8. **Unique Constraints**:
   - Organization domain must be globally unique
   - Company code must be unique within an organization
   - Location code must be unique within a company

---

## Examples

### Complete Workflow Example

#### 1. Create an Organization
```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Solutions Inc.",
    "legal_name": "Tech Solutions Incorporated",
    "domain": "techsolutions.com",
    "email": "info@techsolutions.com",
    "phone": "+1-555-1000",
    "billing_email": "billing@techsolutions.com",
    "max_users": 100,
    "max_companies": 5,
    "max_locations": 20,
    "storage_limit_gb": 50.0
  }'
```

#### 2. Create a Company under the Organization
```bash
curl -X POST http://localhost:8080/api/v1/organizations/507f1f77bcf86cd799439011/companies \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Solutions USA",
    "legal_name": "Tech Solutions USA LLC",
    "code": "TS-USA",
    "email": "usa@techsolutions.com",
    "phone": "+1-555-1001",
    "address": {
      "street": "789 Tech Park",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94102",
      "country": "United States",
      "country_code": "US"
    },
    "is_default": true
  }'
```

#### 3. Create a Warehouse Location
```bash
curl -X POST http://localhost:8080/api/v1/companies/507f1f77bcf86cd799439013/locations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "San Francisco Warehouse",
    "code": "WH-SF-001",
    "type": "warehouse",
    "email": "warehouse-sf@techsolutions.com",
    "phone": "+1-555-1002",
    "address": {
      "street": "100 Warehouse District",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94103",
      "country": "United States",
      "country_code": "US"
    },
    "warehouse_info": {
      "total_area": 25000.0,
      "storage_capacity": 2500,
      "docking_bays": 8,
      "has_cold_storage": false
    },
    "is_default": true
  }'
```

---

## Support

For issues, questions, or feature requests related to the Organization Service, please contact the development team or create an issue in the project repository.

**Version:** 1.0.0  
**Last Updated:** December 15, 2025
