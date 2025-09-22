# Go-Shop Project Structure

## 📁 Root Directory
- **go-shop** - Compiled executable file (binary)
- **go.mod** - Go module dependencies
- **go.sum** - Dependency checksums
- **env.example** - Environment variables template

## 📁 config/
- **config.go** - Application configuration loader
  - Loads environment variables from .env file
  - Database connection settings
  - JWT secret configuration
  - SMTP email settings

## 📁 database/
- **database.go** - PostgreSQL database connection and migrations
  - GORM database initialization
  - Auto-migration for all models
  - Connection pooling configuration
- **redis.go** - Redis connection for caching and temporary data
  - OTP storage
  - Pending user data storage
  - Session management

## 📁 models/
- **user.go** - User data structure and request/response models
  - User registration, login, profile management
- **role.go** - Role-based access control (RBAC) models
  - Role definitions (super_admin, seller, user)
  - User-role relationships
- **product.go** - Product catalog models
  - Product creation, updates, categories
  - Order count tracking for popularity
  - Search request models and search logging
- **category.go** - Product category management
- **order.go** - Order lifecycle models
  - Order statuses: pending, paid, confirmed, shipped, delivered, cancelled
  - Order items and payment tracking
- **favorite.go** - User favorites system
  - Support for products and categories
  - Item type validation

## 📁 handlers/
- **auth.go** - Authentication endpoints
  - User registration with OTP
  - Login/logout
  - Password reset
- **user.go** - User profile management
  - Get/update user profile
  - User information retrieval
- **role.go** - Role management (Admin only)
  - Assign/remove roles
  - Get users by role
  - Role creation and management
- **product.go** - Product catalog endpoints
  - Browse products
  - Product details
  - Category filtering
  - Advanced search with filters and sorting
- **category.go** - Category management
  - List categories
  - Category details
- **order.go** - User order operations
  - Create orders
  - View order history
  - Pay orders (users only)
  - Cancel orders (users only)
- **favorite.go** - Favorites management
  - Add/remove favorites
  - View user favorites
- **admin.go** - Admin operations
  - Product management (CRUD)
  - Category management (CRUD)
  - Order management (confirm, ship, deliver, cancel)
  - User management

## 📁 services/
- **auth.go** - Authentication business logic
  - User registration with email verification
  - OTP generation and validation
  - JWT token management
  - Password hashing and verification
- **user.go** - User management logic
  - Profile updates
  - User data validation
- **role.go** - Role management logic
  - Role assignment/removal with transactions
  - Permission checking
  - User role validation
- **product.go** - Product catalog logic
  - Product CRUD operations
  - Stock management
  - Category relationships
  - Advanced search with filters and sorting
  - Search query logging
- **category.go** - Category management logic
- **order.go** - Order processing logic
  - Order creation with stock validation
  - Status transitions with business rules
  - Payment processing (user payment)
  - Payment confirmation (admin confirmation)
  - Shipping and delivery tracking
  - Product popularity tracking (order_count)
- **favorite.go** - Favorites logic
  - Add/remove items from favorites
  - Duplicate prevention
- **email.go** - Email service
  - OTP emails
  - Password reset emails
  - Welcome emails
  - SMTP configuration

## 📁 middleware/
- **auth.go** - Authentication middleware
  - JWT token validation
  - User context injection
- **role.go** - Role-based access control
  - Super admin middleware
  - Seller middleware
  - Role validation
- **admin.go** - Admin-specific middleware
  - Admin access control
  - Sensitive operation logging

## 📁 routes/
- **routes.go** - API route definitions
  - Public routes (auth, products, categories, search)
  - Protected routes (user operations)
  - Admin routes (super admin only)
  - Seller routes (seller + super admin)
  - Middleware application

## 📁 utils/
- **jwt.go** - JWT token utilities
  - Token generation
  - Token validation
  - Claims extraction
- **password.go** - Password utilities
  - Hashing with bcrypt
  - Password validation
- **otp.go** - OTP generation
  - 6-digit numeric codes
  - Expiration handling

## 🔄 Order Lifecycle Flow
1. **User** creates order (pending)
2. **User** pays order (paid) - via POST /orders/{id}/pay
3. **Super Admin** confirms order (confirmed)
4. **Admin/Seller** ships order (shipped)
5. **Super Admin** delivers order (delivered)
6. **User/Admin** can cancel at any stage (cancelled)

## 🔐 Role Permissions
- **User**: Create orders, pay orders, cancel own orders, manage favorites
- **Seller**: Manage products, ship orders, view all orders
- **Super Admin**: Full access to all operations, user management, role assignment, confirm/deliver orders

## 📧 Email System
- OTP verification for registration
- Password reset functionality
- SMTP configuration required
- Error handling for email failures

## 🗄️ Database Features
- PostgreSQL with GORM ORM
- Soft deletes for data integrity
- Foreign key relationships
- Automatic migrations
- Redis for temporary data storage

## 🔍 Search & Analytics Features
- **Product Search API** (`/api/v1/products/search`)
  - Text search by title (ILIKE)
  - Filter by category, price range
  - Sort by price, popularity, date
  - Pagination support
- **Product Popularity Tracking**
  - `order_count` field in products table
  - Auto-increment on order confirmation
  - Popularity-based sorting
- **Search Analytics**
  - Search query logging
  - User behavior tracking
  - Results analytics
- **Database Optimization**
  - Indexes for search performance
  - Full-text search support
  - Composite indexes for complex queries

## 💳 Payment & Order Management Features
- **User Payment API** (`POST /api/v1/orders/{id}/pay`)
  - Users can mark their orders as paid
  - Only pending orders can be paid
  - Automatic status transition: pending → paid
- **Order Status Management**
  - Complete lifecycle: pending → paid → confirmed → shipped → delivered
  - Role-based status transitions
  - Cancellation at any stage (except delivered)
- **Admin Order Operations**
  - Super Admin: Confirm paid orders (with stock update)
  - Admin/Seller: Ship confirmed orders
  - Super Admin: Deliver shipped orders
  - Admin: Cancel any order
