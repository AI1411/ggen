package main

import (
	"context"

	"gorm.io/gen"

	"g_gen/internal/infra/db"
	applogger "g_gen/internal/infra/logger"
)

func main() {
	ctx := context.Background()

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/domain/query", // 出力パス
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		FieldNullable:     true,
	})

	sqlHandler, err := db.NewSQLHandler(
		db.DefaultDatabaseConfig(),
		applogger.New(applogger.DefaultConfig()),
	)
	if err != nil {
		panic(err)
	}

	g.UseDB(sqlHandler.Conn(ctx))

	// schema_migrationsを除外してテーブル生成
	tables, err := sqlHandler.Conn(ctx).Migrator().GetTables()
	if err != nil {
		panic(err)
	}

	// schema_migrationsを除外
	var filteredTables []interface{}

	for _, tableName := range tables {
		if tableName != "schema_migrations" {
			model := g.GenerateModel(tableName)
			filteredTables = append(filteredTables, model)
		}
	}

	g.ApplyBasic(filteredTables...)

	// Generate the code
	g.Execute()
}
