# User Access Endpoint Documentation

## Overview
This endpoint allows logged-in users to view their organizations, companies, and the locations they can access within the ERP system.

## Endpoint

### Get User Access Information
**GET** `/api/v1/users/me/access`

Returns the organization, companies, and locations that the authenticated user has access to.

#### Authentication
- Requires Bearer token in the Authorization header
- The user ID and organization ID are extracted from the JWT token

#### Response Format
```json
{
  "success": true,
  "message": "User access data retrieved successfully",
  "data": {
    "organization": {
      "id": "507f1f77bcf86cd799439011",
      "name": "Acme Corporation",
      "domain": "acme.com",
      "email": "info@acme.com",
      "phone": "+1-234-567-8900",
      "status": "active",
      "is_active": true,
      ...
    },
    "companies": [
      {
        "company": {
          "id": "507f1f77bcf86cd799439012",
          "organization_id": "507f1f77bcf86cd799439011",
          "name": "Acme Manufacturing",
          "code": "ACM-MFG",
          "email": "mfg@acme.com",
          "phone": "+1-234-567-8901",
          "is_active": true,
          ...
        },
        "locations": [
          {
            "id": "507f1f77bcf86cd799439013",
            "company_id": "507f1f77bcf86cd799439012",
            "organization_id": "507f1f77bcf86cd799439011",
            "name": "Main Warehouse",
            "code": "WH-001",
            "type": "warehouse",
            "email": "warehouse@acme.com",
            "phone": "+1-234-567-8902",
            "address": {
              "street": "123 Industrial Blvd",
              "city": "Springfield",
              "state": "IL",
              "postal_code": "62701",
              "country": "USA"
            },
            "is_active": true,
            ...
          },
          {
            "id": "507f1f77bcf86cd799439014",
            "company_id": "507f1f77bcf86cd799439012",
            "organization_id": "507f1f77bcf86cd799439011",
            "name": "Retail Store Downtown",
            "code": "ST-001",
            "type": "store",
            "email": "downtown@acme.com",
            "phone": "+1-234-567-8903",
            "address": {
              "street": "456 Main Street",
              "city": "Springfield",
              "state": "IL",
              "postal_code": "62701",
              "country": "USA"
            },
            "is_active": true,
            ...
          }
        ]
      }
    ]
  }
}
```

#### Response Fields

##### Organization Object
- `id`: Organization unique identifier
- `name`: Organization display name
- `legal_name`: Legal registered name
- `domain`: Primary domain
- `email`: Organization contact email
- `phone`: Organization contact phone
- `status`: Organization status (active, inactive, suspended, trial, cancelled)
- `is_active`: Boolean indicating if organization is active
- Other organization-specific fields...

##### Company Object
- `id`: Company unique identifier
- `organization_id`: Parent organization ID
- `name`: Company display name
- `legal_name`: Legal registered name
- `code`: Unique company code within organization
- `email`: Company contact email
- `phone`: Company contact phone
- `address`: Company physical address
- `is_active`: Boolean indicating if company is active
- Other company-specific fields...

##### Location Object
- `id`: Location unique identifier
- `company_id`: Parent company ID
- `organization_id`: Organization ID (denormalized)
- `name`: Location display name
- `code`: Unique location code within company
- `type`: Location type (head_office, branch, warehouse, store, factory, distribution_center)
- `email`: Location contact email
- `phone`: Location contact phone
- `address`: Location physical address
- `is_active`: Boolean indicating if location is active
- `warehouse_info`: Additional warehouse-specific information (if type is warehouse)
- `store_info`: Additional store-specific information (if type is store)
- Other location-specific fields...

#### Error Responses

##### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "User not authenticated"
  }
}
```

##### 400 Bad Request
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Organization ID not found"
  }
}
```

##### 500 Internal Server Error
```json
{
  "success": false,
  "error": {
    "code": "FETCH_FAILED",
    "message": "Error message details"
  }
}
```

## Usage Example

### cURL
```bash
curl -X GET "http://localhost:8082/api/v1/users/me/access" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### JavaScript (Fetch API)
```javascript
fetch('http://localhost:8082/api/v1/users/me/access', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => {
  console.log('Organization:', data.data.organization);
  console.log('Companies:', data.data.companies);
})
.catch(error => console.error('Error:', error));
```

### Python (requests)
```python
import requests

headers = {
    'Authorization': f'Bearer {access_token}',
    'Content-Type': 'application/json'
}

response = requests.get(
    'http://localhost:8082/api/v1/users/me/access',
    headers=headers
)

data = response.json()
print('Organization:', data['data']['organization'])
print('Companies:', data['data']['companies'])
```

## Data Model Relationships

```
Organization (1)
    │
    └──> Companies (N)
            │
            └──> Locations (N)
                    │
                    └──> LocationUsers (N) ──> Users (N)
```

- An **Organization** can have multiple **Companies**
- A **Company** belongs to one **Organization** and can have multiple **Locations**
- A **Location** belongs to one **Company** and one **Organization**
- Users access locations through **LocationUser** records (many-to-many relationship)
- The endpoint filters locations based on the user's LocationUser assignments

## Implementation Notes

### Access Control
- Users only see locations they have been explicitly granted access to via the `location_users` collection
- Each `LocationUser` record has:
  - `user_id`: The user's ID
  - `location_id`: The location they can access
  - `is_active`: Whether the access is currently active
  - `access_level`: The level of access (full, read_only, restricted)
  - `role_id`: Location-specific role assignment

### Performance Considerations
- The endpoint performs efficient queries using MongoDB indexes
- Location IDs are retrieved first, then batch queries are used to fetch company and location details
- Results are organized in memory to avoid multiple database round trips
- Indexes exist on:
  - `location_users.user_id`
  - `location_users.location_id`
  - `location_users.is_active`
  - `locations._id`
  - `companies._id`

### Future Enhancements
Potential improvements for this endpoint:
1. Add pagination support for organizations with many companies/locations
2. Add filtering options (e.g., by location type, company, active status)
3. Include user's role and permissions for each location
4. Add search and sorting capabilities
5. Include additional aggregated data (e.g., user counts per location)
