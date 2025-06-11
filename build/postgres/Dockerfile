# Dockerfile
FROM postgres:17

# 必要なビルドツールをインストール（Debian/Ubuntu用）
RUN apt-get update && apt-get install -y \
    build-essential \
    postgresql-server-dev-17 \
    wget \
    tar \
    gzip \
    && rm -rf /var/lib/apt/lists/*

# pg_bigmのダウンロードとビルド（PostgreSQL 17対応版）
RUN cd /tmp && \
    wget https://github.com/pgbigm/pg_bigm/archive/refs/tags/v1.2-20240606.tar.gz && \
    tar -xzf v1.2-20240606.tar.gz && \
    cd pg_bigm-1.2-20240606 && \
    make USE_PGXS=1 && \
    make USE_PGXS=1 install && \
    cd / && rm -rf /tmp/pg_bigm-1.2-20240606*

# ビルドツールをクリーンアップ（システム必須パッケージは保持）
RUN apt-get remove -y build-essential postgresql-server-dev-17 wget && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 初期化SQLスクリプトをコピー
COPY backend/migrations /docker-entrypoint-initdb.d/

# デフォルトの設定を上書き
COPY build/postgres/postgresql.conf /usr/local/share/postgresql/postgresql.conf.sample