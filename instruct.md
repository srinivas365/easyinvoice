You are an expert full-stack engineer. Build a production-ready Pharmacy Shop ERP system using:

Backend:
- Golang
- Gin framework
- SQLite database
- GORM ORM

Frontend:
- Vue 3
- Vite
- TailwindCSS
- Pinia for state management
- Vue Router

Architecture:
- Clean architecture
- Modular structure
- REST API

The ERP must support the following modules:

-----------------------------------
AUTH MODULE
-----------------------------------

Features:
- Login
- Logout
- Simple role support:
  - Admin
  - Staff

Fields:
- id
- username
- password (bcrypt hashed)
- role

-----------------------------------
MEDICINE MANAGEMENT
-----------------------------------

Features:
- Add medicine
- Edit medicine
- Delete medicine
- View medicines
- Search medicines

Fields:
- id
- name
- manufacturer
- batch_number
- expiry_date
- purchase_price
- selling_price
- quantity
- created_at

-----------------------------------
SALES MODULE (POS)
-----------------------------------

Features:
- Create bill
- Add multiple medicines to bill
- Auto calculate total
- Reduce stock automatically

Fields:

Sales table:
- id
- total_amount
- created_at

Sale Items table:
- id
- sale_id
- medicine_id
- quantity
- price
- total

-----------------------------------
PURCHASE MODULE
-----------------------------------

Features:
- Add stock
- Update stock

Fields:
- id
- medicine_id
- quantity
- purchase_price
- created_at

-----------------------------------
DASHBOARD
-----------------------------------

Show:

- Total medicines
- Low stock medicines
- Today sales amount
- Total revenue

-----------------------------------
BACKEND REQUIREMENTS
-----------------------------------

Use proper folder structure:

backend/
    main.go
    config/
    models/
    routes/
    controllers/
    services/
    repository/
    middleware/

Include:

- Database initialization
- Migration
- REST APIs for all modules
- Error handling
- Logging

-----------------------------------
FRONTEND REQUIREMENTS
-----------------------------------

Use proper folder structure:

frontend/
    src/
        components/
        views/
        router/
        store/
        services/

Pages:

- Login Page
- Dashboard
- Medicine List
- Add Medicine
- Sales Page (POS)
- Purchase Page

-----------------------------------
IMPORTANT FEATURES

Must include:

- Proper API integration
- Clean UI
- Responsive design
- Loading indicators
- Error handling

-----------------------------------
DATABASE

Use SQLite file:

pharmacy.db

-----------------------------------
DELIVERABLE

Generate complete working project with:

- Backend code
- Frontend code
- Database models
- API integration
- Instructions to run

-----------------------------------
BONUS FEATURES (if possible)

- Bill print
- Expiry alerts
- Low stock alerts

-----------------------------------

Start building step by step.
First generate backend, then frontend.