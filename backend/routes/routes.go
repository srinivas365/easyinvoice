package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"pharmacy-erp/controllers"
	"pharmacy-erp/middleware"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Initialize controllers
	authController := controllers.NewAuthController()
	medicineController := controllers.NewMedicineController()
	saleController := controllers.NewSaleController()
	purchaseController := controllers.NewPurchaseController()
	purchaseInvoiceController := controllers.NewPurchaseInvoiceController()
	productController := controllers.NewProductController()
	dashboardController := controllers.NewDashboardController()
	settingsController := controllers.NewSettingsController()
	distributorController := controllers.NewDistributorController()

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/login", authController.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", authController.GetProfile)

		// Dashboard
		protected.GET("/dashboard/stats", dashboardController.GetStats)

		// Medicines
		protected.GET("/medicines", medicineController.GetAll)
		protected.GET("/medicines/:id", medicineController.GetByID)
		protected.POST("/medicines", medicineController.Create)
		protected.PUT("/medicines/:id", medicineController.Update)
		protected.DELETE("/medicines/:id", medicineController.Delete)

		// Sales (Invoices)
		protected.GET("/sales", saleController.GetAll)
		protected.GET("/sales/archived", saleController.GetArchived)
		protected.GET("/sales/:id", saleController.GetByID)
		protected.POST("/sales", saleController.Create)
		protected.PUT("/sales/:id", saleController.Update)
		protected.DELETE("/sales/:id", saleController.Delete)
		protected.POST("/sales/:id/restore", saleController.Restore)

		// Purchases
		protected.GET("/purchases", purchaseController.GetAll)
		protected.GET("/purchases/:id", purchaseController.GetByID)
		protected.POST("/purchases", purchaseController.Create)

		// Purchase Invoices (distributor + date + line items)
		protected.GET("/purchase-invoices", purchaseInvoiceController.GetAll)
		protected.GET("/purchase-invoices/:id", purchaseInvoiceController.GetByID)
		protected.POST("/purchase-invoices", purchaseInvoiceController.Create)
		protected.PUT("/purchase-invoices/:id", purchaseInvoiceController.Update)
		protected.DELETE("/purchase-invoices/:id", purchaseInvoiceController.Delete)

		// Products (medicines catalog; purchase invoice items reference these)
		protected.GET("/products", productController.GetAll)
		protected.GET("/products/:id", productController.GetByID)
		protected.POST("/products", productController.Create)
		protected.PUT("/products/:id", productController.Update)
		protected.DELETE("/products/:id", productController.Delete)

		// Settings
		protected.GET("/settings", settingsController.Get)
		protected.PUT("/settings", settingsController.Update)

		// Distributors
		protected.GET("/distributors", distributorController.GetAll)
		protected.GET("/distributors/:id", distributorController.GetByID)
		protected.POST("/distributors", distributorController.Create)
		protected.PUT("/distributors/:id", distributorController.Update)
		protected.DELETE("/distributors/:id", distributorController.Delete)
	}

	return r
}
