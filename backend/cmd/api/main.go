package main

import (
	"go.uber.org/fx"

	"g_gen/internal/di"
	"g_gen/internal/server"
)

// Swagger メタデータ
// @title           農業災害支援システム API
// @version         1.0
// @description     農業災害の報告と支援申請を管理するためのAPI
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	app := fx.New(
		di.Provider(),
		fx.Invoke(server.RegisterRoutes),
	)

	// Run the application
	app.Run()
}
