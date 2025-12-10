# Database Migration Guide

This project supports two migration approaches:

## 🔄 **Option 1: GORM AutoMigrate (Existing System)**

This is your current system that uses GORM's AutoMigrate functionality.

### Usage:
```bash
# Run GORM AutoMigrate
make migrate-gorm
```

### Features:
- ✅ Automatically creates tables based on Go struct definitions
- ✅ Handles foreign keys and relationships
- ✅ Creates custom indexes and triggers
- ✅ Creates PostgreSQL extensions and custom types
- ✅ Maintains dependency order for table creation

---

## 🗃️ **Option 2: SQL Migrations (New System)**

This uses `golang-migrate/migrate` tool with SQL files for more precise control.

### Setup:
```bash
# Install golang-migrate tool
make install-migrate
```

### Usage:
```bash
# Apply all migrations
make migrate-up
# OR
make migrate-sql-up

# Rollback all migrations
make migrate-down
# OR  
make migrate-sql-down

# Apply/rollback one step
make migrate-up-1
make migrate-down-1

# Check current version
make migrate-version
# OR
make migrate-sql-version

# Force to specific version (use carefully!)
make migrate-force version=5

# Create new migration
make migrate name=add_new_table
```

### Available SQL Migration Files:
1. `20250728202902_create_extensions_and_enums` - Extensions and custom types
2. `20250728202929_create_customers_table` - Customer table
3. `20250728203010_create_credit_applications_table` - Credit applications
4. `20250728203038_create_data_room_table` - Data room table
5. `20250728203237_create_bank_statement_tables` - Bank statement tables
6. `20250728203318_create_bank_statement_indicator_tables` - Analysis indicators
7. `20250728203523_create_bank_statement_analysis_table` - Analysis results
8. `20250728203556_create_collateral_docs_table` - Collateral documents
9. `20250728203835_create_slik_ojk_tables` - SLIK OJK tables
10. `20250728204009_create_related_parties_tables` - Related parties
11. `20250728204051_create_junction_tables` - Junction/mapping tables

---

## 🤔 **Which Approach to Use?**

### **Use GORM AutoMigrate if:**
- You want to keep your existing workflow
- You prefer Go-centric migration management
- You want automatic schema sync with your models
- You're doing rapid development

### **Use SQL Migrations if:**
- You need precise control over database schema
- You want to review all schema changes
- You're working in a team environment
- You need to maintain complex migration history
- You want to ensure consistent deployments

---

## 🚨 **Important Notes:**

1. **Don't mix both approaches** - Choose one and stick with it
2. **Backup your database** before running any migrations
3. **Test migrations** in development environment first
4. **Review migration files** before applying in production

---

## 🔧 **Environment Variables Required:**

Make sure your `.env` file contains:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASS=your_password
DB_NAME=your_database
```
