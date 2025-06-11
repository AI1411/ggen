-- 工種区分マスタ
DROP TABLE IF EXISTS work_categories CASCADE;
CREATE TABLE work_categories
(
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- 工種区分ID（主キー、自動掲番）
    category_name VARCHAR(20) NOT NULL,                            -- 工種区分名（漢字表記）
    icon_name     VARCHAR(50) NOT NULL,                            -- アイコンファイル名
    sort_order    INTEGER     NOT NULL DEFAULT 0,                  -- 表示順序（小さいほど上に表示）
    is_active     BOOLEAN     NOT NULL DEFAULT true                -- 有効フラグ（TRUE: 有効、FALSE: 無効）
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_work_categories_category_name ON work_categories (category_name);
CREATE INDEX IF NOT EXISTS idx_work_categories_icon_name ON work_categories (icon_name);
CREATE INDEX IF NOT EXISTS idx_work_categories_sort_order ON work_categories (sort_order);
CREATE INDEX IF NOT EXISTS idx_work_categories_is_active ON work_categories (is_active);

-- テーブルコメント
COMMENT ON TABLE work_categories IS '工種区分マスタ - 工種区分の基本情報を管理';

-- カラムコメント
COMMENT ON COLUMN work_categories.id IS '工種区分ID（主キー、自動掲番）';
COMMENT ON COLUMN work_categories.category_name IS '工種区分名（漢字表記）';
COMMENT ON COLUMN work_categories.icon_name IS 'アイコンファイル名';
COMMENT ON COLUMN work_categories.sort_order IS '表示順序';
COMMENT ON COLUMN work_categories.is_active IS '有効フラグ（TRUE: 有効、FALSE: 無効）';

-- 工種区分マスタに初期データを挿入
INSERT INTO work_categories (category_name, icon_name, sort_order)
VALUES ('農地', '', 10),
       ('水路', '', 20),
       ('農道', '', 30),
       ('ため池', '', 40),
       ('頭首工', '', 50),
       ('揚水機', '', 60),
       ('堤防', '', 70),
       ('橋梁', '', 80),
       ('農地保全施設', '', 90);
