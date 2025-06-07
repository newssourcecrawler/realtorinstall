package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/handlers"
)

func Register(router *gin.Engine, h *handlers.AllHandlers) {
	// 5. Build Gin router with CORS + JWT middleware
	router := gin.Default()
	router.Use(cors.Default())

	// 6. Authentication routes
	router.POST("/login", authH.Login)
	router.Use(AuthMiddleware(authSvc, userRepo))
	router.POST("/register",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "register_user"),
		authH.Register,
	)

	// 7. User CRUD routes
	router.GET("/users", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_user"),
		userH.List,
	)
	router.POST("/users", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_user"),
		userH.Create,
	)
	router.PUT("/users/:id", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_user"),
		userH.Update,
	)
	router.DELETE("/users/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_user"),
		userH.Delete,
	)

	// 8. Property routes
	router.GET("/properties", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_property"),
		propH.List,
	)
	router.POST("/properties",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_property"),
		propH.Create,
	)
	router.PUT("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_property"),
		propH.Update,
	)
	router.DELETE("/properties/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_property"),
		propH.Delete,
	)

	// 9. Buyer routes
	router.GET("/buyers",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_buyer"),
		buyerH.List,
	)
	router.POST("/buyers",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_buyer"),
		buyerH.Create,
	)
	router.PUT("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_buyer"),
		buyerH.Update,
	)
	router.DELETE("/buyers/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_buyer"),
		buyerH.Delete,
	)

	// 10. Pricing routes
	router.GET("/pricing",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_pricing"),
		priceH.List,
	)
	router.POST("/pricing",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_pricing"),
		priceH.Create,
	)
	router.PUT("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_pricing"),
		priceH.Update,
	)
	router.DELETE("/pricing/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_pricing"),
		priceH.Delete,
	)

	// 11. Sales routes
	router.GET("/sales", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_sale"),
		salesH.List,
	)
	router.POST("/sales",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		salesH.Create,
	)
	router.PUT("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_sale"),
		salesH.Update,
	)
	router.DELETE("/sales/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_sale"),
		salesH.Delete,
	)

	// 12. Introduction routes
	router.GET("/introductions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_introduction"),
		introH.List,
	)
	router.POST("/introductions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_introduction"),
		introH.Create,
	)
	router.PUT("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_introduction"),
		introH.Update,
	)
	router.DELETE("/introductions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_introduction"),
		introH.Delete,
	)

	// 13. Lettings routes
	router.GET("/lettings",
		RequirePermission(userRepo, "view_lettings"),
		lettingsH.List,
	)
	router.POST("/lettings",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Create,
	)
	router.PUT("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Update,
	)
	router.DELETE("/lettings/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		lettingsH.Delete,
	)

	// 14. Plan routes
	router.GET("/plans",
		RequirePermission(userRepo, "view_plans"),
		planH.List,
	)
	router.POST("/plans",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Create,
	)
	router.PUT("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Update,
	)
	router.DELETE("/plans/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_sale"),
		planH.Delete,
	)

	// 15. Installment routes
	router.GET("/installments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments"),
		instH.List,
	)
	router.GET("/installments/plan/:planId",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments_byplan"),
		instH.ListByPlan,
	)
	router.POST("/installments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_installments"),
		instH.Create,
	)
	router.PUT("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_installments"),
		instH.Update,
	)
	router.DELETE("/installments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_installments"),
		instH.Delete,
	)

	// 16. Payment routes
	router.GET("/payments", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_payments"),
		payH.List,
	)
	router.POST("/payments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_payments"),
		payH.Create,
	)
	router.PUT("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_payments"),
		payH.Update,
	)
	router.DELETE("/payments/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_payments"),
		payH.Delete,
	)

	// 17. Commission routes
	router.GET("/commissions", AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_commission"),
		commissionH.List,
	)
	router.POST("/commissions",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "create_commission"),
		commissionH.Create,
	)
	router.PUT("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "update_commission"),
		commissionH.Update,
	)
	router.DELETE("/commissions/:id",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "delete_commission"),
		commissionH.Delete,
	)

	// 18. Reporting routes
	router.GET("/reports/commissions/beneficiary",
		RequirePermission(userRepo, "view_commissions_report"),
		reportH.TotalCommissionByBeneficiary,
	)

	router.GET("/reports/installments/outstanding",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_installments_report"),
		reportH.OutstandingInstallmentsByPlan,
	)

	router.GET("/reports/sales/monthly",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_sales_report"),
		reportH.MonthlySalesVolume,
	)

	router.GET("/reports/lettings/rentroll",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_lettings_report"),
		reportH.ActiveLettingsRentRoll,
	)

	router.GET("/reports/properties/top-payments",
		AuthMiddleware(authSvc, userRepo),
		RequirePermission(userRepo, "view_property_payments_report"),
		reportH.TopPropertiesByPaymentVolume,
	)

	router.GET("/export/properties.csv", ExportPropertiesCSV(db))

	router.GET("/export/properties.csv", ImportPropertiesCSV(db))

}
