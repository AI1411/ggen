package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"g_gen/internal/domain/model"
	"g_gen/internal/infra/db"
	applogger "g_gen/internal/infra/logger"
)

const filePath = "cmd/seed/municipality/municipalities.csv"

// CSVRecord CSVの1行を表現する構造体
type CSVRecord struct {
	OrganizationCode      string
	PrefectureNameKanji   string
	MunicipalityNameKanji string
	PrefectureNameKana    string
	MunicipalityNameKana  string
}

// MunicipalityImporter 市町村データインポーター
type MunicipalityImporter struct {
	db *gorm.DB
}

// NewMunicipalityImporter コンストラクタ
func NewMunicipalityImporter(ctx context.Context, db db.Client) *MunicipalityImporter {
	return &MunicipalityImporter{db: db.Conn(ctx)}
}

// ImportFromCSV CSVファイルから市町村データをインポート
func (m *MunicipalityImporter) ImportFromCSV() error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("CSVファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ヘッダー行をスキップ
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("ヘッダー行の読み込みに失敗しました: %w", err)
	}

	var municipalities []*model.Municipality
	lineNum := 1 // ヘッダー行の次から開始
	skippedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("CSV読み込みエラー (行: %d): %w", lineNum+1, err)
		}
		lineNum++

		// CSVレコードのバリデーション
		if len(record) < 5 {
			log.Printf("警告: 行 %d - カラム数が不足しています (期待値: 5以上, 実際: %d)", lineNum, len(record))
			skippedCount++
			continue
		}

		csvRecord := CSVRecord{
			OrganizationCode:      strings.TrimSpace(record[0]),
			PrefectureNameKanji:   strings.TrimSpace(record[1]),
			MunicipalityNameKanji: strings.TrimSpace(record[2]),
			PrefectureNameKana:    strings.TrimSpace(record[3]),
			MunicipalityNameKana:  strings.TrimSpace(record[4]),
		}

		// 市区町村名が空の場合はスキップ
		if csvRecord.MunicipalityNameKanji == "" || csvRecord.MunicipalityNameKana == "" {
			log.Printf("情報: 行 %d - 市区町村名が空のためスキップしました (団体コード: %s)", lineNum, csvRecord.OrganizationCode)
			skippedCount++
			continue
		}

		// データの検証
		if err := m.validateCSVRecord(csvRecord, lineNum); err != nil {
			log.Printf("警告: 行 %d - %v", lineNum, err)
			skippedCount++
			continue
		}

		municipality, err := m.convertToMunicipality(csvRecord)
		if err != nil {
			log.Printf("警告: 行 %d - データ変換エラー: %v", lineNum, err)
			skippedCount++
			continue
		}

		municipalities = append(municipalities, municipality)
	}

	if len(municipalities) == 0 {
		return fmt.Errorf("インポート可能なデータがありませんでした")
	}

	log.Printf("処理結果: %d件をインポート予定、%d件をスキップしました", len(municipalities), skippedCount)

	// バッチでデータベースに挿入
	return m.batchInsert(municipalities)
}

// validateCSVRecord CSVレコードのバリデーション
func (m *MunicipalityImporter) validateCSVRecord(record CSVRecord, lineNum int) error {
	// 団体コードの検証
	if record.OrganizationCode == "" {
		return fmt.Errorf("団体コードが空です")
	}

	if len(record.OrganizationCode) != 6 {
		return fmt.Errorf("団体コードは6桁である必要があります (実際: %d桁)", len(record.OrganizationCode))
	}

	if _, err := strconv.Atoi(record.OrganizationCode); err != nil {
		return fmt.Errorf("団体コードは数値である必要があります: %v", err)
	}

	// 都道府県名の検証
	if record.PrefectureNameKanji == "" {
		return fmt.Errorf("都道府県名（漢字）が空です")
	}

	if record.PrefectureNameKana == "" {
		return fmt.Errorf("都道府県名（カナ）が空です")
	}

	// 市区町村名の検証（都道府県レベルの場合は空でも可）
	if record.MunicipalityNameKanji == "" && record.MunicipalityNameKana != "" {
		return fmt.Errorf("市区町村名（漢字）が空なのに（カナ）が設定されています")
	}

	if record.MunicipalityNameKanji != "" && record.MunicipalityNameKana == "" {
		return fmt.Errorf("市区町村名（漢字）が設定されているのに（カナ）が空です")
	}

	return nil
}

// convertToMunicipality CSVレコードをMunicipalityモデルに変換
func (m *MunicipalityImporter) convertToMunicipality(record CSVRecord) (*model.Municipality, error) {
	// 都道府県コードは団体コードの上2桁
	prefectureCode := record.OrganizationCode[:2]

	// 市区町村名が空の場合は都道府県レベルのデータ
	municipalityNameKanji := record.MunicipalityNameKanji
	municipalityNameKana := record.MunicipalityNameKana

	// 都道府県レベルの場合は空文字列をセット
	if municipalityNameKanji == "" {
		municipalityNameKanji = ""
	}
	if municipalityNameKana == "" {
		municipalityNameKana = ""
	}

	return &model.Municipality{
		PrefectureCode:        prefectureCode,
		OrganizationCode:      record.OrganizationCode,
		PrefectureNameKanji:   record.PrefectureNameKanji,
		MunicipalityNameKanji: municipalityNameKanji,
		PrefectureNameKana:    record.PrefectureNameKana,
		MunicipalityNameKana:  municipalityNameKana,
		IsActive:              true,
	}, nil
}

// batchInsert バッチでデータベースに挿入
func (m *MunicipalityImporter) batchInsert(municipalities []*model.Municipality) error {
	const batchSize = 100

	// トランザクション開始
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 既存データを削除（必要に応じて）
	if err := tx.Where("1 = 1").Delete(&model.Municipality{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("既存データの削除に失敗しました: %w", err)
	}

	// バッチでインサート
	for i := 0; i < len(municipalities); i += batchSize {
		end := i + batchSize
		if end > len(municipalities) {
			end = len(municipalities)
		}

		batch := municipalities[i:end]
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("バッチ挿入に失敗しました (batch %d-%d): %w", i, end-1, err)
		}

		log.Printf("バッチ挿入完了: %d-%d / %d", i+1, end, len(municipalities))
	}

	// トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %w", err)
	}

	log.Printf("インポート完了: %d件のレコードを挿入しました", len(municipalities))
	return nil
}

// GetStatistics インポート統計を取得
func (m *MunicipalityImporter) GetStatistics() (map[string]*int64, error) {
	stats := make(map[string]*int64)

	// 総件数
	if err := m.db.Model(&model.Municipality{}).Count(stats["total"]).Error; err != nil {
		return nil, fmt.Errorf("総件数の取得に失敗しました: %w", err)
	}

	// 都道府県レベルのデータ件数
	if err := m.db.Model(&model.Municipality{}).Where("municipality_name_kanji = ?", "").Count(stats["prefecture_level"]).Error; err != nil {
		return nil, fmt.Errorf("都道府県レベルデータ件数の取得に失敗しました: %w", err)
	}

	// 市区町村レベルのデータ件数
	if err := m.db.Model(&model.Municipality{}).Where("municipality_name_kanji != ?", "").Count(stats["municipality_level"]).Error; err != nil {
		return nil, fmt.Errorf("市区町村レベルデータ件数の取得に失敗しました: %w", err)
	}

	// 有効なデータ件数
	if err := m.db.Model(&model.Municipality{}).Where("is_active = ?", true).Count(stats["active"]).Error; err != nil {
		return nil, fmt.Errorf("有効データ件数の取得に失敗しました: %w", err)
	}

	return stats, nil
}

// 使用例
func main() {
	ctx := context.Background()
	client, err := db.NewSQLHandler(db.DefaultDatabaseConfig(), applogger.New(applogger.DefaultConfig()))
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	importer := NewMunicipalityImporter(ctx, client)

	if err := importer.ImportFromCSV(); err != nil {
		log.Fatal("インポートに失敗しました:", err)
	}
}
