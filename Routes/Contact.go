package Routes

import (
	"example.com/seen-eg/CMS/Controllers"
	"github.com/gofiber/fiber/v2"
)

func ContactRoute(route fiber.Router) {
	route.Post("/new", Controllers.ContactNew)
	route.Post("/get_all/", Controllers.ContactGet)
	route.Post("/get_all/populated/", Controllers.ContactGetPopulated)
	route.Put("/set_status/:id/:new_status", Controllers.ContactSetStatus)
	route.Put("/modify", Controllers.ContactModify)
	route.Post("/get_all_grouped/", Controllers.ContactGetAggregated)

}
