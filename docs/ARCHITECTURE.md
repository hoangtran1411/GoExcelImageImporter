# Cấu trúc dự án (Project Architecture)

Dự án được tổ chức theo mô hình **Wails Architecture**, kết hợp giữa Go (Backend) và Web Technologies (Frontend).

## 🏗️ Sơ đồ cấu trúc

```text
ImageToExcel/
├── main.go               # Wails Entry point, cấu hình cửa sổ
├── app.go                # Backend Logic (Exposed methods cho JS)
├── wails.json            # File cấu hình dự án Wails
├── frontend/             # Giao diện người dùng
│   └── dist/             # HTML/CSS/JS Assets (được embed vào binary)
│       ├── index.html
│       ├── style.css
│       └── app.js
├── internal/             # Code logic sâu (Core Business Logic)
│   └── engine/           # Xử lý logic nghiệp vụ
│       ├── processor.go  # mapping Excel, worker pool, chèn ảnh
│       └── processor_test.go
├── build/                # Thư mục chứa file build output
└── go.mod                # Quản lý dependencies (Go)
```

## ⚙️ Luồng xử lý chính (Main Flow)

1.  **Frontend (JS)**: Người dùng tương tác với giao diện HTML/CSS. Khi nhấn "Start", JS gọi method `Process()` được expose từ Backend.
2.  **Bridge**: Wails Bridge chuyển lời gọi từ JS sang Go method `Process` trong `app.go`.
3.  **App Logic**: `app.go` nhận cấu hình, khởi tạo `Processor` từ `internal/engine`.
4.  **Processor Phase**:
    - **Mapping**: Đọc cột mã sản phẩm từ Excel -> Map.
    - **Dispatching**: Quét thư mục ảnh, tạo Jobs.
    - **Workers**: Xử lý ảnh song song (Scaling, Decode).
    - **Collection**: Gom kết quả và chèn vào Excel (Single Thread safe).
5.  **Feedback**: Trong quá trình, Backend gửi event `progress` ngược lại Frontend. Khi hoàn tất, Frontend hiển thị **Toast Notification** thông báo kết quả chi tiết.

## 🔄 Auto Update Mechanism

Hệ thống cập nhật tự động hoạt động như sau:
1.  **Check**: Khi khởi động, Backend gọi GitHub API kiểm tra latest release.
2.  **Notify**: Nếu có phiên bản mới, gửi tín hiệu cho Frontend hiển thị nút Update.
3.  **Update Action**: Người dùng nhấn Update -> Backend tải file `.exe` mới về thư mục tạm.
4.  **Swap**: Chạy script batch đệm để: Tắt app hiện tại -> Xóa exe cũ -> Move exe mới vào vị trí -> Chạy app mới.

## 🔒 Lưu ý Kỹ thuật
- **Wails Bridge**: Giao tiếp giữa JS và Go là bất đồng bộ (Promise-based).
- **Concurrency**: Sử dụng Goroutines cho việc xử lý ảnh nặng, nhưng ghi file Excel phải tuần tự.
