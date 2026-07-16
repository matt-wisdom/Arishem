package api

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router) {
	scans := app.Group("/scans")
	scans.Post("/code", CreateCodeScan)
	scans.Post("/webapp", CreateWebappScan)
	scans.Get("/", ListScans)
	scans.Get("/:id", GetScan)
	scans.Delete("/:id", CancelScan)

	llmpentest := app.Group("/llmpentest")
	llmpentest.Post("/", CreateLLMPentest)
	llmpentest.Get("/", ListLLMPentests)
	llmpentest.Get("/:id", GetLLMPentest)
	llmpentest.Post("/:id/rerun", RerunLLMPentest)

	reports := app.Group("/reports")
	reports.Get("/", ListReports)
	reports.Get("/:id", GetReport)
	reports.Get("/:id/download", GetReportDownload)

	integrations := app.Group("/integrations")
	integrations.Post("/github", ConnectGitHub)
	integrations.Get("/github", GetGitHubIntegration)
	integrations.Get("/github/repos", ListGitHubRepos)
	integrations.Delete("/github", DisconnectGitHub)

	alerts := app.Group("/alerts")
	alerts.Get("/", ListAlerts)
	alerts.Post("/", CreateAlert)
	alerts.Put("/:id", UpdateAlert)
	alerts.Delete("/:id", DeleteAlert)
	alerts.Post("/test/:id", TestAlert)
}