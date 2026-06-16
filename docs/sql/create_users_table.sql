-- Kích hoạt extension pg_trgm để hỗ trợ tìm kiếm một phần (substring search)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Tạo function cập nhật cột updated_at tự động khi dòng được update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Tạo bảng users (sử dụng số nhiều để tránh trùng với từ khóa 'user' dành riêng trong PostgreSQL)
CREATE TABLE IF NOT EXISTS users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    hash_password text NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

-- Tạo trigger tự động cập nhật thời gian cho cột updated_at trước khi UPDATE
CREATE TRIGGER trigger_update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Tạo index GIN sử dụng trigram để tối ưu hóa truy vấn tìm kiếm một phần không phân biệt chữ hoa/thường (ILIKE) theo tên
CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING gin (name gin_trgm_ops);
