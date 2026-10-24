package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// globalDataDefaultLimit mirrors the legacy global-data behaviour of returning
// the whole list when the caller does not page.
const globalDataDefaultLimit = "1000"

// globalDataModels maps legacy mobile-app model names to the REST list route
// that serves the same records.
var globalDataModels = map[string]string{
	"members":              "members",
	"member":               "members",
	"member_types":         "member-types",
	"member-types":         "member-types",
	"transporters":         "transporters",
	"transporter_vehicles": "transporter-vehicles",
	"transporter-vehicles": "transporter-vehicles",
	"routes":               "routes",
	"route_centers":        "route-centers",
	"route-centers":        "route-centers",
	"centers":              "route-centers",
	"shifts":               "milk-delivery-shifts",
	"milk_delivery_shifts": "milk-delivery-shifts",
	"milk-delivery-shifts": "milk-delivery-shifts",
	"cans":                 "milk-cans",
	"milk_cans":            "milk-cans",
	"milk-cans":            "milk-cans",
	"measuring_cans":       "milk-cans",
	"customers":            "customers",
	"stores":               "stores",
	"store_stocks":         "store-stocks",
	"store-stocks":         "store-stocks",
	"employees":            "employees",
	"suppliers":            "suppliers",
	"banks":                "banks",
	"milk_deliveries":      "milk-deliveries",
	"milk-deliveries":      "milk-deliveries",
	"store_sales":          "store-sales",
	"bank_branches":        "bank-branches",
	"store-sales":          "store-sales",
}

// registerGlobalDataRoute serves GET /api/global-data?model=<name> by
// re-dispatching to the mapped REST route. It is registered outside the
// authenticated group on purpose: the re-dispatched request runs the full
// auth + RBAC chain for the target resource (e.g. members.view), so callers
// need exactly the permission they would need to call that route directly.
func registerGlobalDataRoute(engine *gin.Engine) {
	engine.GET("/api/global-data", func(c *gin.Context) {
		model := strings.ToLower(strings.TrimSpace(c.Query("model")))
		target, ok := globalDataModels[model]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error":            "unsupported global-data model: " + model,
				"supported_models": supportedGlobalDataModels(),
			})
			return
		}

		query := c.Request.URL.Query()
		query.Del("model")

		if search := query.Get("search"); search != "" && query.Get("q") == "" {
			query.Set("q", search)
		}

		limit := firstNonEmpty(query.Get("limit"), query.Get("Limit"), globalDataDefaultLimit)
		page := firstNonEmpty(query.Get("page"), query.Get("Page"), "1")
		query.Set("limit", limit)
		query.Set("Limit", limit)
		query.Set("page", page)
		query.Set("Page", page)

		c.Request.URL.Path = "/api/" + target
		c.Request.URL.RawQuery = query.Encode()
		engine.HandleContext(c)
	})
}

func supportedGlobalDataModels() []string {
	names := make([]string, 0, len(globalDataModels))
	for name := range globalDataModels {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
