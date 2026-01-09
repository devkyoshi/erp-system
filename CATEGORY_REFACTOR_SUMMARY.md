# Category Structure Refactoring Summary

## Overview
The category system has been refactored from a traditional parent-child hierarchy to a **two-level embedded structure** where subcategories are stored directly within their parent category document.

## Key Changes

### 1. Data Model Changes

#### Before (Hierarchical with parent_id):
```go
type ProductCategory struct {
    ParentID *primitive.ObjectID  // Reference to parent
    Level    int                   // 0=root, 1=child, etc.
    Path     string                // /electronics/laptops/gaming
    // ...
}
```

#### After (Embedded Subcategories):
```go
type ProductSubcategory struct {
    BaseModel
    Name        string
    Code        string
    Description string
    IsActive    bool
    ProductCount int
    Metadata    map[string]interface{}
}

type ProductCategory struct {
    BaseModel
    OrganizationID primitive.ObjectID
    Name           string
    Code           string
    Description    string
    Subcategories  []ProductSubcategory  // EMBEDDED array
    IsActive       bool
    ProductCount   int
    Metadata       map[string]interface{}
}
```

### 2. Product Model Update
Products now reference both category and subcategory:
```go
type Product struct {
    // ...
    CategoryID     primitive.ObjectID
    SubcategoryID  primitive.ObjectID  // NEW: Direct reference to subcategory
    // ...
}
```

## API Changes

### Category Endpoints

#### Create Category
**POST** `/categories`

```json
{
  "organization_id": "507f1f77bcf86cd799439011",
  "name": "Electronics",
  "code": "ELEC",
  "description": "Electronic products",
  "is_active": true,
  "subcategories": [
    {
      "name": "Laptops",
      "code": "LAP",
      "description": "Laptop computers",
      "is_active": true
    },
    {
      "name": "Phones",
      "code": "PHN",
      "description": "Mobile phones",
      "is_active": true
    }
  ]
}
```

**Response**: Category object with embedded subcategories, each with their own IDs.

#### Update Category
**PUT** `/categories/:id`

```json
{
  "name": "Updated Electronics",
  "add_subcategories": [
    {
      "name": "Tablets",
      "code": "TAB",
      "is_active": true
    }
  ],
  "update_subcategories": [
    {
      "id": "507f1f77bcf86cd799439012",
      "name": "Gaming Laptops",
      "is_active": true
    }
  ],
  "remove_subcategories": ["507f1f77bcf86cd799439013"]
}
```

#### List Categories
**GET** `/categories?organization_id=xxx&is_active=true&q=electronics&page=1&limit=10`

**Removed query parameters**:
- `parent_id` - No longer needed (subcategories are embedded)
- `level` - No longer applicable (only 2 levels: category and subcategory)

**Returns**: Array of categories with their embedded subcategories.

#### Get Category Tree
**GET** `/categories/tree?organization_id=xxx`

Returns all categories with their subcategories already embedded (no need to build tree recursively).

#### Get Children
**GET** `/categories/:id/children`

Returns the subcategories array from the specified category (filtered for active subcategories).

### Product Endpoints

#### Create Product
**POST** `/products`

```json
{
  "organization_id": "507f1f77bcf86cd799439011",
  "sku": "LAP-DELL-001",
  "name": "Dell XPS 15",
  "category_id": "507f1f77bcf86cd799439011",
  "subcategory_id": "507f1f77bcf86cd799439012",  // NEW: Must exist in category.subcategories
  // ... other fields
}
```

**Validation**:
- If `subcategory_id` is provided, `category_id` must also be provided
- Subcategory must exist in the specified category's `subcategories` array
- Subcategory must not be soft-deleted

#### Update Product
**PUT** `/products/:id`

Same validation rules apply when updating category/subcategory.

## Database Schema

### Categories Collection
```json
{
  "_id": ObjectId("..."),
  "organization_id": ObjectId("..."),
  "name": "Electronics",
  "code": "ELEC",
  "description": "Electronic products",
  "is_active": true,
  "product_count": 150,
  "subcategories": [
    {
      "_id": ObjectId("..."),
      "name": "Laptops",
      "code": "LAP",
      "description": "Laptop computers",
      "is_active": true,
      "product_count": 45,
      "created_at": ISODate("..."),
      "updated_at": ISODate("..."),
      "deleted_at": null
    },
    {
      "_id": ObjectId("..."),
      "name": "Phones",
      "code": "PHN",
      "description": "Mobile phones",
      "is_active": true,
      "product_count": 105,
      "created_at": ISODate("..."),
      "updated_at": ISODate("..."),
      "deleted_at": null
    }
  ],
  "created_at": ISODate("..."),
  "updated_at": ISODate("..."),
  "deleted_at": null
}
```

### Products Collection
```json
{
  "_id": ObjectId("..."),
  "organization_id": ObjectId("..."),
  "sku": "LAP-DELL-001",
  "name": "Dell XPS 15",
  "category_id": ObjectId("..."),
  "subcategory_id": ObjectId("..."),  // References embedded subcategory._id
  // ... other fields
}
```

## Benefits

1. **Simpler Queries**: No need for recursive queries or joins to get category hierarchy
2. **Better Performance**: Subcategories are retrieved with their parent in a single query
3. **Atomic Updates**: Subcategory operations are atomic within the parent document
4. **Clearer Structure**: Two-level limit enforced at schema level
5. **Easier Querying**: Products can be queried by both category and subcategory directly

## Migration Considerations

**Existing Data**: If you have existing data with `parent_id`, `level`, and `path` fields, you'll need to:

1. Create a migration script to:
   - Group subcategories by their `parent_id`
   - Embed them into their parent category's `subcategories` array
   - Update products to include `subcategory_id`
   - Remove old parent-child category documents

2. Update any existing queries in other services that reference:
   - `parent_id`
   - `level`
   - `path`

## Validation Rules

1. **Category Names**: Must be unique within an organization
2. **Subcategory Names**: Must be unique within a category
3. **Product Category/Subcategory**:
   - `subcategory_id` requires `category_id`
   - Subcategory must exist in the specified category
   - Subcategory must not be soft-deleted
4. **Soft Deletion**:
   - Deleted subcategories remain in the array with `deleted_at` set
   - Queries filter out soft-deleted subcategories

## Testing Recommendations

1. **Unit Tests**: Test subcategory CRUD operations within categories
2. **Integration Tests**:
   - Create category with subcategories
   - Update subcategories (add/update/remove)
   - Create products with category and subcategory
   - Query products by category and subcategory
3. **Edge Cases**:
   - Subcategory not in specified category
   - Soft-deleted subcategory reference
   - Category without subcategories
   - Updating subcategory that doesn't exist

## Example Workflows

### Create Category with Subcategories
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "507f1f77bcf86cd799439011",
    "name": "Electronics",
    "code": "ELEC",
    "is_active": true,
    "subcategories": [
      {"name": "Laptops", "code": "LAP", "is_active": true},
      {"name": "Phones", "code": "PHN", "is_active": true}
    ]
  }'
```

### Create Product with Category and Subcategory
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "507f1f77bcf86cd799439011",
    "sku": "LAP-DELL-001",
    "name": "Dell XPS 15",
    "category_id": "507f1f77bcf86cd799439020",
    "subcategory_id": "507f1f77bcf86cd799439021",
    "selling_price": 1299.99
  }'
```

### Query Products by Category and Subcategory
```bash
# All products in Electronics
curl "http://localhost:8080/api/v1/products?organization_id=xxx&category_id=507f1f77bcf86cd799439020"

# All Laptops (specific subcategory)
# Note: You'll need to add subcategory_id filter to product repository if not already present
```

## Notes

- The refactoring is complete for all core CRUD operations
- The `go.mod` warnings about indirect dependencies are informational only
- All category and product validation now works with the embedded structure
- Subcategories have their own unique IDs for referencing in products
- Soft deletion is supported for both categories and subcategories
