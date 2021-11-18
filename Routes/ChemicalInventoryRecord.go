package Routes

import (
	"example.com/seen-eg/CMS/Controllers"
	"github.com/gofiber/fiber/v2"
)

func ChemicalInventoryRecordRoute(route fiber.Router) {
	route.Post("/check_move_validity/:id/:consider_opened", Controllers.ChemicalInventoryRecordCheckIfValidLocation)
	route.Post("/op_open/:id", Controllers.ChemicalInventoryRecordOpenOp)
	route.Post("/op_set_status/:id", Controllers.ChemicalInventoryRecordSetStatusOp)
	route.Post("/op_move_to/:id/:consider_opened", Controllers.ChemicalInventoryRecordMoveOp)
	route.Post("/get_all", Controllers.ChemicalInventoryRecordGetAll)
	route.Post("/get_all/populated", Controllers.ChemicalInventoryRecordGetAllPopulated)
	route.Get("/get_all/valid", Controllers.ChemicalInventoryRecordGetAllValid)
	route.Post("/get_all_less_than_safety_stock", Controllers.ChemicalInventoryRecordGetLessThanSafetyStock)
	route.Get("/get_all/aggregate", Controllers.ChemicalInventoryRecordGetAllPopulatedAggregated)
	route.Get("/get_count_alarm_action_compliance", Controllers.ChemicalInventoryRecordCountAlarmActionCompliance)
}
