package main

import (
	"log"
	"os"

	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"pharmacy-erp/routes"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize database
	config.InitDB()

	// Run migrations
	migrate()

	// Create default admin user if not exists
	createDefaultAdmin()

	// Setup routes
	r := routes.SetupRoutes()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func migrate() {
	// First run AutoMigrate for models that don't need special handling
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Medicine{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Purchase{},
		&models.Distributor{},
		&models.Product{},
		&models.PurchaseInvoice{},
		&models.PurchaseInvoiceItem{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Handle Settings table separately to add new columns if they don't exist
	var settingsTableExists int
	config.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='settings'").Scan(&settingsTableExists)

	if settingsTableExists > 0 {
		// Check if address column exists
		var addressExists int
		config.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('settings') WHERE name='address'").Scan(&addressExists)
		if addressExists == 0 {
			config.DB.Exec("ALTER TABLE settings ADD COLUMN address TEXT")
			log.Println("Added address column to settings table")
		}

		// Check if phone_number column exists
		var phoneExists int
		config.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('settings') WHERE name='phone_number'").Scan(&phoneExists)
		if phoneExists == 0 {
			config.DB.Exec("ALTER TABLE settings ADD COLUMN phone_number TEXT")
			log.Println("Added phone_number column to settings table")
		}
	}

	// Now run AutoMigrate for Settings to ensure all columns are correct
	err = config.DB.AutoMigrate(&models.Settings{})
	if err != nil {
		log.Fatal("Migration failed for Settings:", err)
	}

	// Handle sale_items table separately to avoid GORM conflicts
	var tableExists int
	config.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sale_items'").Scan(&tableExists)

	if tableExists > 0 {
		// Table exists - check if it needs migration
		var medicineIdNotNull int
		config.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('sale_items') WHERE name='medicine_id' AND \"notnull\"=1").Scan(&medicineIdNotNull)

		var colExists int
		config.DB.Raw("SELECT COUNT(*) FROM pragma_table_info('sale_items') WHERE name='discount_type'").Scan(&colExists)

		// If medicine_id is NOT NULL or discount_type doesn't exist, we need to migrate
		if medicineIdNotNull > 0 || colExists == 0 {
			log.Println("Migrating sale_items table...")

			config.DB.Exec("PRAGMA foreign_keys=off")
			config.DB.Exec("BEGIN TRANSACTION")

			// Create backup
			config.DB.Exec("CREATE TABLE sale_items_backup AS SELECT * FROM sale_items")

			// Drop old table
			config.DB.Exec("DROP TABLE sale_items")

			config.DB.Exec("COMMIT")
			config.DB.Exec("PRAGMA foreign_keys=on")

			// Let GORM create the table with correct schema
			err = config.DB.AutoMigrate(&models.SaleItem{})
			if err != nil {
				log.Fatal("Migration failed for SaleItem:", err)
			}

			// Copy data back
			config.DB.Exec("PRAGMA foreign_keys=off")
			if colExists > 0 {
				config.DB.Exec(`
					INSERT INTO sale_items 
					SELECT id, sale_id, medicine_id, medicine_name, quantity, price, discount, 
					       COALESCE(discount_type, 'amount'), total, created_at, updated_at, deleted_at 
					FROM sale_items_backup
				`)
			} else {
				config.DB.Exec(`
					INSERT INTO sale_items (id, sale_id, medicine_id, medicine_name, quantity, price, discount, discount_type, total, created_at, updated_at, deleted_at)
					SELECT id, sale_id, medicine_id, medicine_name, quantity, price, discount, 
					       'amount', total, created_at, updated_at, deleted_at 
					FROM sale_items_backup
				`)
			}
			config.DB.Exec("DROP TABLE sale_items_backup")
			config.DB.Exec("PRAGMA foreign_keys=on")

			log.Println("sale_items table migrated successfully")
		} else {
			// Table is already correct, just verify with GORM (but don't let it modify)
			// We skip AutoMigrate to avoid conflicts
			log.Println("sale_items table schema is up to date")
		}
	} else {
		// Table doesn't exist - let GORM create it
		err = config.DB.AutoMigrate(&models.SaleItem{})
		if err != nil {
			log.Fatal("Migration failed for SaleItem:", err)
		}
	}

	log.Println("Migrations completed successfully")
}

func createDefaultAdmin() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     "admin",
		}
		config.DB.Create(&admin)
		log.Println("Default admin user created (username: admin, password: admin123)")

		// Create a staff user as well
		hashedPasswordStaff, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
		staff := models.User{
			Username: "staff",
			Password: string(hashedPasswordStaff),
			Role:     "staff",
		}
		config.DB.Create(&staff)
		log.Println("Default staff user created (username: staff, password: staff123)")
	}
}
