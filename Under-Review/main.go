package main

import (
	"fmt"

	"example.com/seen-eg/CMS/DBManager"
	"example.com/seen-eg/CMS/Middlewares"
	"example.com/seen-eg/CMS/Routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.StorageConditionRoute(app.Group("/storage_condition"))
	Routes.SafetyHazardRoute(app.Group("/safety_hazard"))
	Routes.LabDivisionRoute(app.Group("/lab_division"))
	Routes.ContactRoute(app.Group("/Contact"))
	Routes.UnitsOfMeasurementRoute(app.Group("/uom"))
	Routes.ChemicalTemplateRoute(app.Group("/chemical_template"))
	Routes.BatchRecevingRoute(app.Group("/batch_receving"))
	Routes.ChemicalInventoryRecordRoute(app.Group("/chemical_inventory_record"))
	Routes.ChemicalSolutionTemplateRoute(app.Group("/chemical_solution_template"))
	Routes.EquipmentCalibration(app.Group("/equipment_calibration"))
	Routes.TestingMethodsRoute(app.Group("/testing_methods"))
	Routes.UserRoute(app.Group("/user"))
	Routes.UserRoleRoute(app.Group("/user_roles"))
	Routes.ApprovalCycleRoute(app.Group("/approval_cycle"))
	Routes.SettingsRoute(app)
	Routes.ChemicalSolutionPreparationRoute(app.Group("/chemical_solution_preparation"))
	Routes.ChemicalSolutionInventoryRecordRoute(app.Group("/chemical_solution_inventory_record"))
	Routes.ChemicalInsightsRoute(app.Group("/chemical_insights"))
	Routes.StepRoute(app.Group("/step"))
	Routes.StandardizationLogEventRoute(app.Group("/chemical_solution_inventory_record/standard_logic_event"))
}

func main() {
	fmt.Println("Hello CMS")

	fmt.Print("Initializing DataBase Connections ... ")
	initState := DBManager.InitCMSCollections()
	if initState {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")
	app := fiber.New()
	app.Use(cors.New())
	app.Use(Middlewares.Auth)
	app.Static("/Resources", "./Resources")
	app.Use(pprof.New())

	// Create new group route on path "/debug/pprof"
	/*zft := app.Group("/debug")
	zft.Get("/pprof", func(c *fiber.Ctx) error {
		pprofhandler.PprofHandler(&fasthttp.RequestCtx{})
		return nil
	})*/

	// bench marking
	// TestCases.ReqCntTest()

	SetupRoutes(app)
	fmt.Println("[OK]")

	app.Listen(":8080")

}
