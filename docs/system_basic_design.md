# Tài liệu Thiết kế Cơ bản (System Basic Design) - Meow Application

Tài liệu này mô tả chi tiết kiến trúc, luồng hoạt động, cấu trúc cơ sở dữ liệu và các giao thức API của **Meow Application** — Hệ thống giám sát tài nguyên máy chủ và nhật ký hệ thống thời gian thực, tích hợp cơ chế xác thực người dùng.

---

## 1. Giới thiệu (Introduction)

### 1.1. Mục tiêu hệ thống
Meow Application được xây dựng với mục tiêu cung cấp một bảng điều khiển (Dashboard) trực quan giúp lập trình viên và quản trị viên hệ thống giám sát sức khỏe của máy chủ Backend (được phát triển bằng Go) theo thời gian thực.

### 1.2. Phạm vi hệ thống
Hệ thống bao gồm hai thành phần chính:
- **Backend (Go / Gin Framework)**: Thu thập thông tin tài nguyên hệ thống (CPU, RAM, Uptime, Goroutines), quản lý thông tin tài khoản người dùng, xác thực bằng JWT, kiểm tra kết nối cơ sở dữ liệu, đọc nhật ký hoạt động (`app.log`) và cung cấp API RESTful.
- **Frontend (SvelteKit / Svelte 5 / Vite)**: Gọi API giám sát định kỳ, trực quan hóa dữ liệu bộ nhớ RAM dưới dạng biểu đồ đường (Area Line Chart) và biểu đồ tròn (Donut Gauge Chart), hiển thị log trực tiếp dạng Console Terminal, đồng thời cung cấp giao diện Đăng ký/Đăng nhập cho người dùng.

---

## 2. Kiến trúc hệ thống (System Architecture)

Hệ thống được thiết kế theo mô hình client-server tách biệt rõ ràng, giao tiếp thông qua giao thức HTTP RESTful API, hỗ trợ cấu hình chia sẻ tài nguyên nguồn gốc chéo (CORS) và xác thực dựa trên JSON Web Token (JWT).

```mermaid
graph TD
    A["Client (Web Browser)"] <-->|Hiển thị & Tương tác| B["Frontend (SvelteKit)"]
    B <-->|HTTP RESTful API với JWT (CORS)| C["Backend (Go / Gin)"]
    C <-->|Đọc/Ghi dữ liệu users| D[("Database (PostgreSQL)")]
    C <-->|Ghi & Đọc logs| E["Local File (app.log)"]
```

### Các thành phần chính:
1. **Frontend**: Chạy trên môi trường Node.js (cổng mặc định `5173`). Sử dụng Svelte 5 Runes để quản lý state phản ứng (reactive state) hiệu quả, biểu diễn biểu đồ động thông qua SVG và lưu trữ JWT Token trong Client (ví dụ LocalStorage hoặc Cookie) để xác thực các request.
2. **Backend**: Chạy dưới dạng ứng dụng Go biên dịch (cổng mặc định `8080`). Sử dụng Gin framework để định tuyến API và xử lý HTTP requests.
3. **Database**: PostgreSQL dùng để lưu trữ dữ liệu lâu dài (bao gồm thông tin tài khoản người dùng trong bảng `users` và kết nối thông qua Connection Pool bằng thư viện `pgx/v5`).
4. **Log Storage**: File cục bộ `app.log` đóng vai trò ghi nhận lịch sử hoạt động và sự cố phục vụ công tác gỡ lỗi.

---

## 3. Luồng dữ liệu & Sơ đồ Sequence (Data Flow & Sequence Diagram)

### 3.1. Luồng Giám sát sức khỏe (Health Check)
Quy trình cốt lõi của ứng dụng là **Giám sát sức khỏe (Health Check)**, được tự động kích hoạt định kỳ mỗi 10 giây từ phía Frontend hoặc chạy thủ công khi người dùng click nút làm mới (refresh).

```mermaid
sequenceDiagram
    autonumber
    actor User as Người dùng
    participant FE as "Frontend (SvelteKit)"
    participant BE as "Backend (Go/Gin)"
    participant DB as "Database (PostgreSQL)"
    participant Log as "File: app.log"

    Note over FE, BE: Tự động kích hoạt mỗi 10 giây hoặc khi click nút Refresh
    User->>FE: Yêu cầu cập nhật chỉ số
    FE->>BE: GET /api/health
    
    Note over BE: Xử lý song song bằng Goroutines (sync.WaitGroup)
    par Đọc Log ứng dụng
        BE->>Log: Đọc 50 dòng log cuối cùng
        Log-->>BE: Trả về danh sách logs dạng mảng chuỗi
    and Kiểm tra Database
        BE->>DB: Ping kết nối (pgxpool.Ping)
        DB-->>BE: Trả về kết quả (CONNECTED/DISCONNECTED)
    end
    
    Note over BE: Thu thập thông tin runtime hệ thống<br/>(runtime.MemStats, Uptime, CPU cores)
    BE-->>FE: Trả về dữ liệu JSON dạng HealthResponse
    
    alt Xử lý khi thành công (200 OK)
        FE->>FE: Cập nhật thông số (OS, CPU cores, Uptime)
        FE->>FE: Tính phần trăm RAM đã dùng = (Alloc / Sys) * 100
        FE->>FE: Thêm điểm dữ liệu mới vào lịch sử (giữ tối đa 15 điểm để vẽ biểu đồ SVG)
        FE->>FE: Định dạng cú pháp màu cho logs (GET/POST, INFO/ERROR...) và cuộn xuống cuối terminal
    else Xử lý khi lỗi (500 Error hoặc Mất kết nối)
        FE->>FE: Đánh dấu trạng thái OFFLINE và hiển thị cảnh báo lỗi kết nối
    end
```

### 3.2. Luồng Đăng ký & Đăng nhập (Authentication Flow)
Luồng này thực hiện đăng ký và đăng nhập tài khoản người dùng, cấp phát JWT Access Token để định danh phiên làm việc.

```mermaid
sequenceDiagram
    autonumber
    actor User as Người dùng
    participant FE as "Frontend (SvelteKit)"
    participant BE as "Backend (Go/Gin)"
    participant DB as "Database (PostgreSQL)"

    rect rgb(240, 248, 255)
        Note over User, DB: Luồng Đăng ký (Register)
        User->>FE: Nhập Name, Email, Password và nhấn Đăng ký
        FE->>BE: POST /api/auth/register (JSON Payload)
        BE->>BE: Validate dữ liệu đầu vào (tên, email, độ dài mật khẩu)
        BE->>BE: Mã hóa băm mật khẩu (Hash Password)
        BE->>DB: INSERT INTO users (name, email, hash_password)
        DB-->>BE: Trả về ID của người dùng mới tạo
        BE->>BE: Tạo JWT Access Token chứa Claim UserID
        BE-->>FE: Trả về 200 OK với JWT Token
        FE->>FE: Lưu JWT Token và chuyển hướng về Trang chủ (Dashboard)
    end

    rect rgb(245, 245, 245)
        Note over User, DB: Luồng Đăng nhập (Login)
        User->>FE: Nhập Email, Password và nhấn Đăng nhập
        FE->>BE: POST /api/auth/login (JSON Payload)
        BE->>BE: Validate định dạng dữ liệu đầu vào
        BE->>DB: SELECT id, hash_password FROM users WHERE email = ?
        DB-->>BE: Trả về ID và mật khẩu đã băm
        BE->>BE: So khớp mật khẩu đầu vào với mật khẩu đã băm
        BE->>BE: Tạo JWT Access Token chứa Claim UserID
        BE-->>FE: Trả về 200 OK với JWT Token
        FE->>FE: Lưu JWT Token và cập nhật trạng thái đã đăng nhập
    end
```

---

## 4. Thiết kế Cơ sở dữ liệu (Database Schema Design)

### 4.1. Hiện trạng kết nối
Hệ thống sử dụng **PostgreSQL** làm cơ sở dữ liệu chính. Việc kết nối được quản lý thông qua Connection Pool bằng thư viện `github.com/jackc/pgx/v5/pgxpool` giúp tối ưu tài nguyên và quản lý vòng đời kết nối an toàn.

- **Biến môi trường**: `DATABASE_URL` dạng `postgres://username:password@host:port/database`
- **Khởi tạo Pool**: Thực hiện một lần duy nhất (`sync.Once`) khi server bắt đầu chạy.
- **Thời gian chờ kết nối (Timeout)**: 5 giây. Nếu quá thời gian này mà không kết nối được hoặc không `Ping` thành công, ứng dụng Backend sẽ dừng khởi chạy ngay lập tức (`log.Fatalf`).

### 4.2. Danh sách các bảng
Hệ thống chứa bảng nghiệp vụ `users` để phục vụ xác thực người dùng.

#### Bảng `users`
Bảng lưu trữ thông tin tài khoản người dùng của hệ thống.

| Tên Cột (Column Name) | Kiểu Dữ Liệu (Data Type) | Ràng Buộc (Constraints) | Mô Tả / Mục Đích (Description / Purpose) |
| :--- | :--- | :--- | :--- |
| `id` | `bigint` | `PRIMARY KEY`, `GENERATED ALWAYS AS IDENTITY` | Khóa chính tự động tăng của người dùng (sử dụng GENERATED ALWAYS AS IDENTITY). |
| `name` | `text` | `NOT NULL` | Tên của người dùng, được đánh chỉ mục GIN (`idx_users_name_trgm` sử dụng trigram) để tìm kiếm nhanh, không phân biệt hoa/thường. |
| `email` | `text` | `UNIQUE`, `NOT NULL` | Địa chỉ email dùng để đăng nhập và định danh duy nhất. |
| `hash_password` | `text` | `NOT NULL` | Mật khẩu của người dùng đã được băm (hash). |
| `created_at` | `timestamptz` | `DEFAULT CURRENT_TIMESTAMP` | Thời điểm tài khoản được tạo ra. |
| `updated_at` | `timestamptz` | `DEFAULT CURRENT_TIMESTAMP` | Thời điểm tài khoản được cập nhật lần cuối, tự động cập nhật qua trigger `trigger_update_users_updated_at`. |

> [!NOTE]
> Để xem chi tiết các Chỉ mục (Indexes), Trigger và Hàm tự động cập nhật, vui lòng tham khảo tài liệu chi tiết tại [users_table.md](file:///home/hungg562002/Project/meow/docs/db_schema/users_table.md). Toàn bộ danh sách các bảng dữ liệu được quản lý và chỉ mục tại [db_schema-lock.json](file:///home/hungg562002/Project/meow/docs/db_schema/db_schema-lock.json).

---

## 5. Định nghĩa API (API Contracts)

### 5.1. API Health Check

API duy nhất được sử dụng để lấy trạng thái hệ thống, tài nguyên RAM/CPU và log hoạt động.

- **HTTP Method**: `GET`
- **Path**: `/api/health`
- **Headers**:
  - `Accept: application/json`

#### Response thành công (HTTP Status: `200 OK`)

Trả về toàn bộ thông tin tài nguyên và log ứng dụng.

**Payload mẫu (JSON):**
```json
{
  "status": "UP",
  "database_status": "CONNECTED",
  "server_info": {
    "os": "linux",
    "architecture": "amd64",
    "go_version": "go1.26.3",
    "num_cpu": 8,
    "num_goroutine": 4,
    "uptime": "15m32s",
    "memory_usage": {
      "alloc": "2.50 MB",
      "total_alloc": "10.20 MB",
      "sys": "15.30 MB",
      "num_gc": 1
    }
  },
  "app_logs": [
    "2026/05/30 16:14:05 Server starting on port 8080",
    "[GIN-debug] Listening and serving HTTP on :8080",
    "[GIN-debug] GET /api/health --> controllers.(*HealthController).GetHealth-fm (3 handlers)"
  ]
}
```

#### Response lỗi (HTTP Status: `500 Internal Server Error`)

Xảy ra khi Backend gặp sự cố trong quá trình đọc file log `app.log`.

**Payload mẫu (JSON):**
```json
{
  "error": "Không thể tải log ứng dụng: open app.log: permission denied"
}
```

### 5.2. API Đăng ký tài khoản (Register)

API dùng để đăng ký tài khoản người dùng mới và tự động sinh JWT Access Token để đăng nhập.

- **HTTP Method**: `POST`
- **Path**: `/api/auth/register`
- **Headers**:
  - `Content-Type: application/json`

#### Request Body
| Thuộc tính | Kiểu dữ liệu | Bắt buộc | Mô tả |
| :--- | :--- | :--- | :--- |
| `name` | String | Có | Họ tên hiển thị của người dùng |
| `email` | String | Có | Địa chỉ email người dùng (phải đúng định dạng email) |
| `password` | String | Có | Mật khẩu truy cập (tối thiểu 6 ký tự) |

**Payload mẫu (JSON):**
```json
{
  "name": "Meow Admin",
  "email": "admin@meow.io",
  "password": "secretpassword"
}
```

#### Response thành công (HTTP Status: `200 OK`)
Trả về JWT Access Token để người dùng sử dụng xác thực trong các phiên làm việc tiếp theo.

**Payload mẫu (JSON):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTgzMTIwMDB9.xxxxxx"
}
```

#### Response lỗi (HTTP Status: `400 Bad Request`)
Xảy ra khi dữ liệu đầu vào bị thiếu hoặc không đúng định dạng (ví dụ mật khẩu ngắn hơn 6 ký tự hoặc email sai định dạng).

**Payload mẫu (JSON):**
```json
{
  "error": "Dữ liệu đầu vào không hợp lệ: Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}
```

#### Response lỗi (HTTP Status: `500 Internal Server Error`)
Xảy ra khi có lỗi ghi cơ sở dữ liệu (ví dụ trùng email do cột `email` có ràng buộc UNIQUE) hoặc lỗi hệ thống khác.

**Payload mẫu (JSON):**
```json
{
  "error": "failed to create user: ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)"
}
```

### 5.3. API Đăng nhập tài khoản (Login)

API dùng để xác thực thông tin tài khoản và sinh JWT Access Token cho phiên làm việc hiện tại.

- **HTTP Method**: `POST`
- **Path**: `/api/auth/login`
- **Headers**:
  - `Content-Type: application/json`

#### Request Body
| Thuộc tính | Kiểu dữ liệu | Bắt buộc | Mô tả |
| :--- | :--- | :--- | :--- |
| `email` | String | Có | Địa chỉ email tài khoản (phải đúng định dạng email) |
| `password` | String | Có | Mật khẩu tài khoản |

**Payload mẫu (JSON):**
```json
{
  "email": "admin@meow.io",
  "password": "secretpassword"
}
```

#### Response thành công (HTTP Status: `200 OK`)
Trả về JWT Access Token để người dùng sử dụng xác thực trong các phiên làm việc tiếp theo.

**Payload mẫu (JSON):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTgzMTIwMDB9.xxxxxx"
}
```

#### Response lỗi (HTTP Status: `400 Bad Request`)
Xảy ra khi dữ liệu đầu vào thiếu trường bắt buộc hoặc sai định dạng email.

**Payload mẫu (JSON):**
```json
{
  "error": "Dữ liệu đầu vào không hợp lệ: Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

#### Response lỗi (HTTP Status: `401 Unauthorized`)
Xảy ra khi sai địa chỉ email hoặc mật khẩu không chính xác.

**Payload mẫu (JSON):**
```json
{
  "error": "Sai địa chỉ email hoặc mật khẩu"
}
```

#### Response lỗi (HTTP Status: `500 Internal Server Error`)
Xảy ra khi có lỗi hệ thống hoặc kết nối cơ sở dữ liệu bị ngắt quãng.

**Payload mẫu (JSON):**
```json
{
  "error": "failed to get hash password from email: connection refused"
}
```
