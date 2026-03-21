# Pharmacy Shop ERP System

A production-ready Pharmacy Shop ERP system built with Golang (Gin) backend and Vue 3 frontend.

## Features

### Authentication
- Login/Logout functionality
- Role-based access (Admin, Staff)
- JWT token-based authentication

### Medicine Management
- Add, Edit, Delete medicines
- View all medicines with search functionality
- Track medicine details (name, manufacturer, batch number, expiry date, prices, quantity)

### Sales Module (POS)
- Create sales bills
- Add multiple medicines to a bill
- Auto-calculate totals
- Automatic stock reduction
- Bill printing functionality

### Purchase Module
- Add stock purchases
- Update medicine stock automatically
- Track purchase history

### Dashboard
- Total medicines count
- Low stock alerts
- Today's sales amount
- Total revenue
- Expiry date alerts

## Tech Stack

### Backend
- **Golang** 1.21+
- **Gin** - Web framework
- **GORM** - ORM
- **SQLite** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing

### Frontend
- **Vue 3** - Progressive JavaScript framework
- **Vite** - Build tool
- **TailwindCSS** - Utility-first CSS framework
- **Pinia** - State management
- **Vue Router** - Routing
- **Axios** - HTTP client

## Project Structure

```
easyinvoice/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── medicine.go
│   │   ├── sale.go
│   │   └── purchase.go
│   ├── routes/
│   │   └── routes.go
│   ├── controllers/
│   │   ├── auth_controller.go
│   │   ├── medicine_controller.go
│   │   ├── sale_controller.go
│   │   ├── purchase_controller.go
│   │   └── dashboard_controller.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── medicine_service.go
│   │   ├── sale_service.go
│   │   └── purchase_service.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── medicine_repository.go
│   │   ├── sale_repository.go
│   │   └── purchase_repository.go
│   └── middleware/
│       └── auth.go
│
└── frontend/
    ├── package.json
    ├── vite.config.js
    ├── tailwind.config.js
    ├── index.html
    └── src/
        ├── main.js
        ├── App.vue
        ├── style.css
        ├── components/
        │   └── Layout.vue
        ├── views/
        │   ├── Login.vue
        │   ├── Dashboard.vue
        │   ├── Medicines.vue
        │   ├── AddMedicine.vue
        │   ├── EditMedicine.vue
        │   ├── Sales.vue
        │   └── Purchases.vue
        ├── router/
        │   └── index.js
        ├── store/
        │   └── auth.js
        └── services/
            ├── api.js
            ├── authService.js
            ├── medicineService.js
            ├── saleService.js
            ├── purchaseService.js
            └── dashboardService.js
```

## Setup Instructions

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 16.x or higher
- **npm** or **yarn**

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Run the backend server:
```bash
go run main.go
```

The backend server will start on `http://localhost:8080`

The database file `pharmacy.db` will be created automatically in the backend directory.

**Default Users:**
- Admin: `admin` / `admin123`
- Staff: `staff` / `staff123`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

### Production Build

#### Backend
```bash
cd backend
go build -o pharmacy-erp
./pharmacy-erp
```

#### Frontend
```bash
cd frontend
npm run build
```

The built files will be in the `dist` directory. Serve them using any static file server.

## API Endpoints

### Authentication
- `POST /api/login` - Login user
- `GET /api/profile` - Get current user profile (Protected)

### Dashboard
- `GET /api/dashboard/stats` - Get dashboard statistics (Protected)

### Medicines
- `GET /api/medicines` - Get all medicines (with optional `?search=query`)
- `GET /api/medicines/:id` - Get medicine by ID
- `POST /api/medicines` - Create new medicine
- `PUT /api/medicines/:id` - Update medicine
- `DELETE /api/medicines/:id` - Delete medicine

### Sales
- `GET /api/sales` - Get all sales
- `GET /api/sales/:id` - Get sale by ID
- `POST /api/sales` - Create new sale

### Purchases
- `GET /api/purchases` - Get all purchases
- `GET /api/purchases/:id` - Get purchase by ID
- `POST /api/purchases` - Create new purchase

All endpoints except `/api/login` require authentication via Bearer token in the Authorization header.

## Database Schema

### Users
- id (Primary Key)
- username (Unique)
- password (bcrypt hashed)
- role (admin/staff)
- created_at, updated_at

### Medicines
- id (Primary Key)
- name
- manufacturer
- batch_number
- expiry_date
- purchase_price
- selling_price
- quantity
- created_at, updated_at

### Sales
- id (Primary Key)
- total_amount
- created_at, updated_at

### Sale Items
- id (Primary Key)
- sale_id (Foreign Key)
- medicine_id (Foreign Key)
- quantity
- price
- total
- created_at, updated_at

### Purchases
- id (Primary Key)
- medicine_id (Foreign Key)
- quantity
- purchase_price
- created_at, updated_at

## Bonus Features

✅ **Bill Print** - Print bills after completing a sale
✅ **Low Stock Alerts** - Dashboard shows medicines with low stock (≤10 units)
✅ **Expiry Alerts** - Medicines with expiry dates are tracked

## Environment Variables

You can set the following environment variables:

- `PORT` - Backend server port (default: 8080)
- `JWT_SECRET` - JWT secret key (default: "pharmacy-erp-secret-key-change-in-production")

## Development

### Running Both Servers

You can run both servers simultaneously:

**Terminal 1 (Backend):**
```bash
cd backend
go run main.go
```

**Terminal 2 (Frontend):**
```bash
cd frontend
npm run dev
```

## License

This project is built for educational/demonstration purposes.

## Support

For issues or questions, please check the code comments or create an issue in the repository.
