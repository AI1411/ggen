package main

import (
	"context"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"g_gen/internal/domain/model"
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

	// Generate the code
	g.Execute()

	// 生成したModelにRelation情報を手動追加（これだけは手動対応が必要）
	allModels := []any{
		g.GenerateModel(
			model.TableNameMunicipality,
			gen.FieldRelateModel(field.BelongsTo, "Prefecture", model.Prefecture{}, &field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"PrefectureCode"},
					"references": []string{"Code"},
				},
			}),
		),
		g.GenerateModel(
			model.TableNamePrefecture,
			gen.FieldRelateModel(field.HasMany, "Municipalities", model.Municipality{}, &field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"PrefectureCode"},
					"references": []string{"Code"},
				},
			}),
		),
		g.GenerateModel(model.TableNameWorkCategory),
	}

	g.ApplyBasic(allModels...)

	// Generate the code
	g.Execute()
}
