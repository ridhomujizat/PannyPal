# Cash Flow API Documentation

API lengkap untuk monitoring pengeluaran dan pemasukan cash flow.

## Base URL
```
http://localhost:9001/api
```

## Authentication
API ini menggunakan phone_number sebagai identifier user. Setiap request memerlukan phone_number baik di body (POST/PUT) atau query parameter (GET/DELETE).

---

## Transaction APIs

### 1. Create Transaction
**POST** `/transactions`

Create new income or expense transaction.

**Request Body:**
```json
{
  "phone_number": "6285811588248",
  "amount": 100000.00,
  "category_id": 1,
  "type": "INCOME", // or "EXPENSE"
  "description": "Gaji bulanan" // optional
}
```

**Response:**
```json
{
  "code": 201,
  "message": "Transaction created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Salary"
    },
    "amount": 100000.00,
    "description": "Gaji bulanan",
    "transaction_date": "2025-12-08T10:30:00Z",
    "type": "INCOME",
    "created_at": "2025-12-08T10:30:00Z",
    "updated_at": "2025-12-08T10:30:00Z"
  }
}
```

### 2. Get Transactions
**GET** `/transactions?phone_number=6285811588248&page=1&limit=10`

List all transactions with filters and pagination.

**Query Parameters:**
- `phone_number` (required): User's phone number
- `type` (optional): "INCOME" or "EXPENSE"
- `category_id` (optional): Filter by category ID
- `start_date` (optional): Filter transactions from this date (ISO format)
- `end_date` (optional): Filter transactions until this date (ISO format)
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page (max: 100)

**Response:**
```json
{
  "code": 200,
  "message": "Transactions retrieved successfully",
  "data": {
    "transactions": [
      {
        "id": 1,
        "user_id": 1,
        "category_id": 1,
        "category": {
          "id": 1,
          "name": "Salary"
        },
        "amount": 100000.00,
        "description": "Gaji bulanan",
        "transaction_date": "2025-12-08T10:30:00Z",
        "type": "INCOME",
        "created_at": "2025-12-08T10:30:00Z",
        "updated_at": "2025-12-08T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

### 3. Get Transaction by ID
**GET** `/transactions/{id}?phone_number=6285811588248`

Get specific transaction by ID.

**Response:** Same as single transaction object above.

### 4. Update Transaction
**PUT** `/transactions/{id}?phone_number=6285811588248`

Update existing transaction. All fields are optional.

**Request Body:**
```json
{
  "amount": 120000.00,
  "category_id": 2,
  "type": "EXPENSE",
  "description": "Updated description"
}
```

### 5. Delete Transaction
**DELETE** `/transactions/{id}?phone_number=6285811588248`

Delete transaction.

**Response:**
```json
{
  "code": 200,
  "message": "Transaction deleted successfully",
  "data": null
}
```

### 6. Get Transaction Summary
**GET** `/transactions/summary?phone_number=6285811588248`

Get summary of transactions with totals and category breakdown.

**Query Parameters:**
- `phone_number` (required): User's phone number
- `start_date` (optional): Summary from this date
- `end_date` (optional): Summary until this date
- `month` (optional): Filter by specific month (1-12)
- `year` (optional): Filter by specific year

**Response:**
```json
{
  "code": 200,
  "message": "Transaction summary retrieved successfully",
  "data": {
    "total_income": 500000.00,
    "total_expense": 300000.00,
    "balance": 200000.00,
    "transaction_count": 15,
    "income_count": 5,
    "expense_count": 10,
    "category_summary": [
      {
        "category_id": 1,
        "category_name": "Food",
        "type": "EXPENSE",
        "total_amount": 150000.00,
        "count": 8,
        "percentage": 50.0
      }
    ],
    "period": {
      "start_date": "2025-12-01T00:00:00Z",
      "end_date": "2025-12-31T23:59:59Z",
      "month": 12,
      "year": 2025
    }
  }
}
```

---

## Category APIs

### 1. Create Category
**POST** `/categories`

**Request Body:**
```json
{
  "name": "Food & Drinks"
}
```

### 2. Get All Categories
**GET** `/categories`

**Response:**
```json
{
  "code": 200,
  "message": "Categories retrieved successfully",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Food & Drinks",
        "created_at": "2025-12-08T10:30:00Z",
        "updated_at": "2025-12-08T10:30:00Z"
      }
    ]
  }
}
```

### 3. Get Category by ID
**GET** `/categories/{id}`

### 4. Update Category
**PUT** `/categories/{id}`

**Request Body:**
```json
{
  "name": "Updated Category Name"
}
```

### 5. Delete Category
**DELETE** `/categories/{id}`

---

## Budget APIs

### 1. Create Budget
**POST** `/budgets`

Set monthly budget for a category.

**Request Body:**
```json
{
  "phone_number": "6285811588248",
  "category_id": 1,
  "amount": 500000.00,
  "month": 12,
  "year": 2025
}
```

### 2. Get Budgets
**GET** `/budgets?phone_number=6285811588248`

**Query Parameters:**
- `phone_number` (required): User's phone number
- `month` (optional): Filter by month
- `year` (optional): Filter by year
- `category_id` (optional): Filter by category

### 3. Get Budget by ID
**GET** `/budgets/{id}?phone_number=6285811588248`

### 4. Update Budget
**PUT** `/budgets/{id}?phone_number=6285811588248`

### 5. Delete Budget
**DELETE** `/budgets/{id}?phone_number=6285811588248`

### 6. Get Budget Status
**GET** `/budgets/status?phone_number=6285811588248&month=12&year=2025`

Check budget usage status.

**Response:**
```json
{
  "code": 200,
  "message": "Budget status retrieved successfully",
  "data": {
    "budget_statuses": [
      {
        "budget_id": 1,
        "category_id": 1,
        "category_name": "Food",
        "budget_amount": 500000.00,
        "spent_amount": 350000.00,
        "remaining_amount": 150000.00,
        "percentage_used": 70.0,
        "is_over_budget": false,
        "month": 12,
        "year": 2025
      }
    ],
    "total_budget": 500000.00,
    "total_spent": 350000.00,
    "total_remaining": 150000.00,
    "month": 12,
    "year": 2025
  }
}
```

---

## Analytics APIs

### 1. Monthly Analytics
**GET** `/analytics/monthly?phone_number=6285811588248&year=2025`

Get monthly breakdown for a year.

**Response:**
```json
{
  "code": 200,
  "message": "Monthly analytics retrieved successfully",
  "data": {
    "data": [
      {
        "month": 12,
        "month_name": "December",
        "year": 2025,
        "total_income": 500000.00,
        "total_expense": 300000.00,
        "balance": 200000.00,
        "transaction_count": 15
      }
    ],
    "year": 2025,
    "total_income": 6000000.00,
    "total_expense": 3600000.00,
    "total_balance": 2400000.00,
    "best_month": {
      "month": 8,
      "month_name": "August",
      "balance": 400000.00
    },
    "worst_month": {
      "month": 3,
      "month_name": "March",
      "balance": 50000.00
    }
  }
}
```

### 2. Yearly Analytics
**GET** `/analytics/yearly?phone_number=6285811588248&start_year=2020&end_year=2025`

Get yearly comparison.

### 3. Category Analytics
**GET** `/analytics/categories?phone_number=6285811588248&type=EXPENSE`

Get spending/income by category.

**Query Parameters:**
- `phone_number` (required): User's phone number
- `start_date` (optional): Analysis from this date
- `end_date` (optional): Analysis until this date
- `type` (optional): "INCOME" or "EXPENSE"

**Response:**
```json
{
  "code": 200,
  "message": "Category analytics retrieved successfully",
  "data": {
    "data": [
      {
        "category_id": 1,
        "category_name": "Food",
        "type": "EXPENSE",
        "total_amount": 150000.00,
        "count": 8,
        "percentage": 50.0,
        "average_amount": 18750.00
      }
    ],
    "total_amount": 300000.00,
    "transaction_count": 15,
    "period": {
      "start_date": "2025-12-01T00:00:00Z",
      "end_date": "2025-12-31T23:59:59Z"
    },
    "top_category": {
      "category_id": 1,
      "category_name": "Food",
      "total_amount": 150000.00
    }
  }
}
```

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "code": 400,
  "message": "Validation error",
  "data": "phone_number is required",
  "error": "validation failed"
}
```

**Common Error Codes:**
- `400` - Bad Request (validation errors)
- `404` - Not Found (user/resource not found)
- `403` - Forbidden (access denied)
- `500` - Internal Server Error

---

## Usage Examples

### 1. Complete Cash Flow Setup

1. **Create Categories:**
```bash
curl -X POST http://localhost:9001/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Salary"}'

curl -X POST http://localhost:9001/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Food"}'
```

2. **Set Monthly Budget:**
```bash
curl -X POST http://localhost:9001/api/budgets \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "6285811588248",
    "category_id": 2,
    "amount": 500000.00,
    "month": 12,
    "year": 2025
  }'
```

3. **Record Transactions:**
```bash
# Income
curl -X POST http://localhost:9001/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "6285811588248",
    "amount": 3000000.00,
    "category_id": 1,
    "type": "INCOME",
    "description": "Monthly salary"
  }'

# Expense
curl -X POST http://localhost:9001/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "6285811588248",
    "amount": 150000.00,
    "category_id": 2,
    "type": "EXPENSE",
    "description": "Lunch"
  }'
```

4. **Monitor Budget Status:**
```bash
curl "http://localhost:9001/api/budgets/status?phone_number=6285811588248&month=12&year=2025"
```

5. **Get Analytics:**
```bash
curl "http://localhost:9001/api/analytics/monthly?phone_number=6285811588248&year=2025"
```

### 2. Daily Workflow

**Morning - Check yesterday's transactions:**
```bash
curl "http://localhost:9001/api/transactions?phone_number=6285811588248&start_date=2025-12-07&end_date=2025-12-07"
```

**Throughout the day - Record expenses:**
```bash
curl -X POST http://localhost:9001/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "6285811588248",
    "amount": 25000.00,
    "category_id": 2,
    "type": "EXPENSE",
    "description": "Coffee"
  }'
```

**Evening - Check budget status:**
```bash
curl "http://localhost:9001/api/budgets/status?phone_number=6285811588248"
```

**End of month - Generate summary:**
```bash
curl "http://localhost:9001/api/transactions/summary?phone_number=6285811588248&month=12&year=2025"
```