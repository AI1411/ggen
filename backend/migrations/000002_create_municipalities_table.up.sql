-- 市町村マスタテーブル
-- 全国の都道府県・市区町村の基本情報を管理する
DROP TABLE IF EXISTS municipalities CASCADE;
CREATE TABLE IF NOT EXISTS municipalities
(
    id                      BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,     -- 内部ID（主キー、自動採番）
    prefecture_code         VARCHAR(2)   NOT NULL REFERENCES prefectures (code), -- 都道府県ID（外部キー）
    organization_code       VARCHAR(6)   NOT NULL UNIQUE,                        -- 団体コード（総務省地方公共団体コード）
    prefecture_name_kanji   VARCHAR(10)  NOT NULL,                               -- 都道府県名（漢字表記）
    municipality_name_kanji VARCHAR(50)  NOT NULL,                               -- 市区町村名（漢字表記）
    prefecture_name_kana    VARCHAR(20)  NOT NULL,                               -- 都道府県名（カタカナ表記）
    municipality_name_kana  VARCHAR(100) NOT NULL,                               -- 市区町村名（カタカナ表記）
    is_active               BOOLEAN      NOT NULL DEFAULT TRUE                   -- 有効フラグ（デフォルトは有効）合併などで無効化する場合に使用
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_municipalities_prefecture_id ON municipalities (prefecture_code);
CREATE INDEX IF NOT EXISTS idx_municipalities_organization_code ON municipalities (organization_code);
CREATE INDEX IF NOT EXISTS idx_municipalities_prefecture_kanji ON municipalities (prefecture_name_kanji);
CREATE INDEX IF NOT EXISTS idx_municipalities_municipality_kanji ON municipalities (municipality_name_kanji);
CREATE INDEX IF NOT EXISTS idx_municipalities_pref_muni_kanji ON municipalities (prefecture_name_kanji, municipality_name_kanji);
CREATE INDEX IF NOT EXISTS idx_municipalities_prefecture_kana ON municipalities (prefecture_name_kana);
CREATE INDEX IF NOT EXISTS idx_municipalities_municipality_kana ON municipalities (municipality_name_kana);
CREATE INDEX IF NOT EXISTS idx_municipalities_is_active ON municipalities (is_active);

-- テーブルコメント
COMMENT ON TABLE municipalities IS '市町村マスタテーブル - 全国の都道府県・市区町村の基本情報を管理';

-- カラムコメント
COMMENT ON COLUMN municipalities.id IS '内部ID（主キー、自動採番）';
COMMENT ON COLUMN municipalities.prefecture_code IS '都道府県ID（外部キー、都道府県マスタのID）';
COMMENT ON COLUMN municipalities.organization_code IS '団体コード（総務省地方公共団体コード、6桁）';
COMMENT ON COLUMN municipalities.prefecture_name_kanji IS '都道府県名（漢字表記）';
COMMENT ON COLUMN municipalities.municipality_name_kanji IS '市区町村名（漢字表記）';
COMMENT ON COLUMN municipalities.prefecture_name_kana IS '都道府県名（カタカナ表記）';
COMMENT ON COLUMN municipalities.municipality_name_kana IS '市区町村名（カタカナ表記）';
COMMENT ON COLUMN municipalities.is_active IS '有効フラグ（TRUE: 有効、FALSE: 無効）';
