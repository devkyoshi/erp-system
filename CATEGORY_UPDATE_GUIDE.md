# Category Update with Subcategory Management

## Overview

The category update endpoint now supports comprehensive subcategory management, allowing you to add, update, and remove subcategories in a single update operation. When removing subcategories, the system automatically validates that they don't have any associated products.

## Features

1. **Add New Subcategories**: Create new subcategories under the category being updated
2. **Update Existing Subcategories**: Modify subcategory details (name, code, description, etc.)
3. **Remove Subcategories**: Delete subcategories with automatic product validation

## API Endpoint

```
PUT /categories/{id}
```

## Request Body Structure

```json
{
  "name": "Updated Category Name",
  "code": "UPDATED_CODE",
  "description": "Updated description",
  "is_active": true,
  "metadata": {
    "custom_field": "value"
  },
  "add_subcategories": [
    {
      "name": "New Subcategory 1",
      "code": "SUB1",
      "description": "First new subcategory",
      "is_active": true,
      "metadata": {}
    },
    {
      "name": "New Subcategory 2",
      "code": "SUB2",
      "description": "Second new subcategory with nested subcategories",
      "is_active": true,
      "subcategories": [
        {
          "name": "Nested Subcategory",
          "code": "NESTED1",
          "description": "This is nested",
          "is_active": true
        }
      ]
    }
  ],
  "update_subcategories": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Updated Subcategory Name",
      "description": "Updated description",
      "is_active": false
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "code": "UPDATED_SUB_CODE"
    }
  ],
  "remove_subcategories": [
    "507f1f77bcf86cd799439013",
    "507f1f77bcf86cd799439014"
  ]
}
```

## Field Descriptions

### Main Category Fields (Optional)
- `name`: New name for the category
- `code`: New code for the category
- `description`: New description
- `parent_id`: Move category to a different parent (use empty string to make it a root category)
- `is_active`: Enable or disable the category
- `metadata`: Custom metadata as key-value pairs

### Subcategory Management

#### `add_subcategories` (array)
Array of new subcategories to create. Each subcategory has:
- `name` (required): Name of the subcategory
- `code`: Code for the subcategory
- `description`: Description
- `is_active`: Whether subcategory is active (default: true)
- `metadata`: Custom metadata
- `subcategories`: Nested subcategories (supports unlimited nesting)

#### `update_subcategories` (array)
Array of existing subcategories to update. Each update has:
- `id` (required): The subcategory ID to update
- `name`: New name (optional)
- `code`: New code (optional)
- `description`: New description (optional)
- `is_active`: New active status (optional)
- `metadata`: New metadata (optional)

#### `remove_subcategories` (array of strings)
Array of subcategory IDs to remove. The system will:
1. Verify each subcategory is a direct child of the category
2. Check if the subcategory has any child categories (prevents removal if yes)
3. Check if the subcategory has any products (prevents removal if yes)
4. Only remove if all validations pass

## Validation Rules

### For Removing Subcategories
1. **Direct Child Validation**: Can only remove direct children of the category
2. **Has Children Check**: Cannot remove subcategories that have child categories
3. **Product Check**: Cannot remove subcategories that have products assigned to them
4. **Authorization**: Must belong to the same organization

### For Adding Subcategories
1. **Name Uniqueness**: Name must be unique at the same level (among siblings)
2. **Organization Inheritance**: Automatically inherits organization from parent
3. **Authorization**: Must belong to the same organization

### For Updating Subcategories
1. **Existence Check**: Subcategory must exist
2. **Direct Child Validation**: Must be a direct child of the category
3. **Name Uniqueness**: If name is changed, must be unique among siblings
4. **Authorization**: Must belong to the same organization

## Example Use Cases

### Example 1: Add New Subcategories Only

```json
{
  "add_subcategories": [
    {
      "name": "Smart Watches",
      "code": "SMART_WATCHES",
      "is_active": true
    },
    {
      "name": "Fitness Trackers",
      "code": "FITNESS",
      "is_active": true
    }
  ]
}
```

### Example 2: Update Existing Subcategories

```json
{
  "update_subcategories": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Smartwatches",
      "description": "Modern wearable smartwatches",
      "is_active": true
    }
  ]
}
```

### Example 3: Remove Subcategories

```json
{
  "remove_subcategories": [
    "507f1f77bcf86cd799439013",
    "507f1f77bcf86cd799439014"
  ]
}
```

**Response if subcategory has products:**
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "failed to remove subcategories: cannot remove subcategories with products: [Obsolete Category, Deprecated Category]"
  }
}
```

### Example 4: Complete Update (All Operations)

```json
{
  "name": "Electronics & Gadgets",
  "description": "Updated comprehensive electronics category",
  "is_active": true,
  "add_subcategories": [
    {
      "name": "Gaming Accessories",
      "code": "GAMING_ACC",
      "is_active": true
    }
  ],
  "update_subcategories": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Smartphones & Tablets",
      "is_active": true
    }
  ],
  "remove_subcategories": [
    "507f1f77bcf86cd799439015"
  ]
}
```

## Response Format

### Success Response (200 OK)

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439010",
    "organization_id": "507f191e810c19729de860ea",
    "name": "Electronics & Gadgets",
    "code": "ELECTRONICS",
    "description": "Updated comprehensive electronics category",
    "parent_id": null,
    "level": 0,
    "path": "/electronics & gadgets",
    "is_active": true,
    "product_count": 42,
    "created_at": "2026-01-05T10:30:00Z",
    "updated_at": "2026-01-08T14:25:30Z"
  },
  "message": "Category updated successfully"
}
```

### Error Responses

#### 400 Bad Request - Validation Error
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid subcategory ID '123': encoding/hex: odd length hex string"
  }
}
```

#### 400 Bad Request - Subcategory Has Products
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "failed to remove subcategories: cannot remove subcategories with products: [Laptops, Smartphones]"
  }
}
```

#### 400 Bad Request - Subcategory Has Children
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "failed to remove subcategories: cannot remove subcategory 'Mobile Devices' because it has child categories"
  }
}
```

#### 400 Bad Request - Not a Direct Child
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "failed to update subcategories: category 'Laptops' is not a direct subcategory of the parent"
  }
}
```

#### 404 Not Found - Category Not Found
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "category not found"
  }
}
```

#### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UPDATE_FAILED",
    "message": "unauthorized: category belongs to different organization"
  }
}
```

## Implementation Details

### Repository Methods Added

1. **`FindByIDs(ctx, ids)`**: Fetch multiple categories by their IDs
2. **`HasProductsInCategories(ctx, categoryIDs)`**: Check multiple categories for products in bulk
3. **`DeleteMultiple(ctx, ids)`**: Soft delete multiple categories at once

### Service Methods Added

1. **`handleRemoveSubcategories()`**: Validates and removes subcategories
   - Checks if subcategories exist
   - Validates they're direct children
   - Checks for child categories
   - Checks for products
   - Performs batch deletion

2. **`handleUpdateSubcategories()`**: Updates existing subcategories
   - Validates existence and hierarchy
   - Applies partial updates
   - Checks name uniqueness

3. **`handleAddSubcategories()`**: Creates new subcategories
   - Inherits organization from parent
   - Supports nested subcategory creation
   - Validates name uniqueness

## Testing Recommendations

1. **Test removing subcategories with products**: Should fail with appropriate error
2. **Test removing subcategories with children**: Should fail with appropriate error
3. **Test removing empty subcategories**: Should succeed
4. **Test updating subcategory names**: Should validate uniqueness
5. **Test adding nested subcategories**: Should create full hierarchy
6. **Test combined operations**: Add, update, and remove in single request
7. **Test authorization**: Try operations across organizations
8. **Test invalid subcategory IDs**: Should return validation errors

## Notes

- All operations are performed within the same request transaction
- Subcategory removal is processed first, then updates, then additions
- Product validation prevents accidental data loss
- The system supports unlimited nesting levels for subcategories
- All operations respect organization boundaries for security
