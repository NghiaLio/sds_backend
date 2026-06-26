# Hướng Dẫn Chạy Backend - Windows & macOS

Dự án này là một REST API Backend đơn giản viết bằng Golang, quản lý danh sách Task (Công việc), sử dụng cơ sở dữ liệu SQLite (`tasks.db`) và cơ chế xác thực JWT Access Token.

Dự án đã được cấu hình biên dịch thuần Go (Pure-Go SQLite) nên **không cần cài đặt GCC compiler** khi build hoặc chạy.

---

## 🚀 1. Cách Chạy Ngay (Sử dụng File Đã Build Sẵn)

Các file mã máy đã được biên dịch sẵn nằm trong thư mục `bin/`. Bạn có thể chạy trực tiếp mà không cần cài đặt Go trên máy.

### 💻 Trên Windows:
1. Mở **PowerShell** hoặc **CMD** tại thư mục dự án `d:\Backend`.
2. Chạy lệnh sau để khởi động server:
   ```powershell
   .\bin\backend-windows-amd64.exe
   ```

### 🍎 Trên macOS (MacBook):
Khi chạy một file thực thi chia sẻ từ máy khác, hệ điều hành macOS sẽ chặn mặc định (Lỗi Security / Unidentified Developer). Hãy làm theo các bước sau để chạy:

1. Mở **Terminal** và di chuyển vào thư mục dự án.
2. Cấp quyền chạy (chọn đúng phiên bản chip của máy bạn):
   * **Nếu là máy chip Apple Silicon (M1/M2/M3/M4)**:
     ```bash
     chmod +x ./bin/backend-darwin-arm64
     ```
   * **Nếu là máy chip Intel cũ**:
     ```bash
     chmod +x ./bin/backend-darwin-amd64
     ```
3. **Mở khóa chặn bảo mật của macOS (Bắt buộc)**:
   * **Với chip Apple Silicon**:
     ```bash
     xattr -d com.apple.quarantine ./bin/backend-darwin-arm64
     ```
   * **Với chip Intel**:
     ```bash
     xattr -d com.apple.quarantine ./bin/backend-darwin-amd64
     ```
4. Khởi chạy server:
   * **Với chip Apple Silicon**:
     ```bash
     ./bin/backend-darwin-arm64
     ```
   * **Với chip Intel**:
     ```bash
     ./bin/backend-darwin-amd64
     ```

---

## 🛠️ 2. Cách Tự Biên Dịch Lại (Compile From Source)

Nếu bạn thay đổi mã nguồn Go và muốn tự build ra các file binary mới:

### 💻 Trên Windows:
Chạy file script batch:
```powershell
.\build.bat
```
Script sẽ tự động dọn dẹp thư viện (`go mod tidy`) và tạo mới các file nhị phân cho cả Windows, Linux, macOS trong thư mục `bin/`.

### 🍎 Trên macOS / Linux:
1. Mở **Terminal** và di chuyển vào thư mục dự án.
2. Cấp quyền chạy file script shell:
   ```bash
   chmod +x ./build.sh
   ```
3. Chạy lệnh build:
   ```bash
   ./build.sh
   ```

---

## ⚙️ 3. Chạy Chế Độ Phát Triển (Development Mode)

Yêu cầu máy bạn đã cài đặt **Go SDK** (Khuyên dùng phiên bản Go 1.22 trở lên).

1. Di chuyển vào thư mục dự án:
   ```bash
   cd d:\Backend
   ```
2. Khởi chạy trực tiếp:
   ```bash
   go run main.go
   ```

---

## 🧪 4. Kiểm tra & Sử dụng API

Khi server hoạt động, địa chỉ mặc định là: `http://localhost:8080`.

1. **Khởi tạo dữ liệu mẫu**: 
   Khi chạy lần đầu tiên, hệ thống sẽ tự động khởi tạo file database `tasks.db`, tạo tài khoản mặc định `admin` (mật khẩu: `admin123`) và tự động seed 10 tasks mẫu liên kết với tài khoản này.
   
2. **Kiểm tra tài liệu chi tiết**:
   Xem file [api_docs.md](file:///d:/Backend/api_docs.md) để biết chi tiết các tham số đầu vào và định dạng JSON phản hồi của từng API.

3. **Import vào Postman**:
   Import file [postman_collection.json](file:///d:/Backend/postman_collection.json) vào Postman để gửi yêu cầu test. Bộ collection này đã được chia làm 2 thư mục:
   * **Auth**: Chứa các API `Register User` và `Login User`. Khi chạy API `Login`, Postman sẽ tự động lưu token vào biến môi trường.
   * **Tasks**: Chứa các API CRUD Task. Thư mục này được cấu hình tự động kế thừa và đính kèm `Bearer Token` từ biến môi trường của Postman vào header.
