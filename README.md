# Golang Excel Image Importer

Một công cụ hiệu chỉnh Excel mạnh mẽ được viết bằng Go và Wails, giúp tự động chèn hình ảnh vào bảng tính dựa trên mã sản phẩm. Công cụ này kết hợp sức mạnh xử lý của Go với giao diện hiện đại của Web (HTML/CSS/JS).

## ✨ Tính năng nổi bật

- **🚀 Hiệu suất vượt trội:** Backend Go xử lý ảnh và Excel cực nhanh với Worker Pool.
- **🎨 Giao diện hiện đại:** Dark Mode cao cấp, **Toast Notification** mượt mà & Responsive.
- **� Auto Update:** Tự động kiểm tra và cập nhật phiên bản mới nhất từ GitHub Releases.
- **�💾 Tối ưu bộ nhớ:** Stream dữ liệu Excel để xử lý file lớn mà không tốn nhiều RAM.
- **🔍 Tìm kiếm thông minh:** Tự động khớp tên file ảnh với mã sản phẩm linh hoạt.
- **📦 Nhẹ và Nhanh:** Ứng dụng Wails sử dụng WebView2 có sẵn trên Windows, file thực thi nhỏ gọn (~10MB).

## 🛠️ Công nghệ sử dụng

- **Backend:** [Go (Golang)](https://golang.org/)
- **Framework:** [Wails v2](https://wails.io/)
- **Frontend:** HTML, CSS (Custom Premium Theme), JavaScript
- **Thư viện Excel:** [Excelize v2](https://github.com/xuri/excelize)

## 💻 Sự tương thích & Yêu cầu hệ thống

Công cụ này được tối ưu hóa cho môi trường Windows. Dưới đây là chi tiết về khả năng tương thích:

### Hệ điều hành hỗ trợ
| Phiên bản | Trạng thái | Ghi chú |
| :--- | :--- | :--- |
| **Windows 11** | ✅ Tốt nhất | Hoạt động hoàn hảo, WebView2 đã có sẵn. |
| **Windows 10** | ✅ Tốt nhất | Hoạt động hoàn hảo, WebView2 thường đã có sẵn (hoặc qua Windows Update). |
| **Windows 7 / 8 / 8.1** | ⚠️ Hạn chế | Yêu cầu cài đặt [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) bản dành cho Win 7/8. Microsoft đã ngừng hỗ trợ chính thức. |
| **Windows Server** | ✅ Hỗ trợ | Hoạt động tốt trên Windows Server 2016 trở lên (cần WebView2). |

### Yêu cầu phần mềm & Phần cứng
- **Kiến trúc:** Windows 64-bit (x64) là bắt buộc.
- **WebView2:** Yếu tố then chốt để hiển thị giao diện.
- **RAM:** Tối thiểu 2GB (Khuyến nghị 4GB+ để xử lý mượt mà hàng nghìn ảnh).
- **Bộ nhớ:** Khoảng 50MB cho ứng dụng và file tạm.

### Dành cho nhà phát triển (Build từ nguồn)
- **Go:** 1.20 trở lên (Project đang dùng 1.25.5).
- **Wails CLI:** Chạy lệnh `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

### Cài đặt Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Cài đặt và Chạy dự án
1. Clone dự án:
   ```bash
   git clone <repository-url>
   cd ImageToExcel
   ```

2. Cài đặt dependencies và chạy Dev Mode:
   ```bash
   wails dev
   ```
   Lệnh này sẽ tự động cài đặt Go modules và Frontend assets, sau đó mở ứng dụng.

### 🔨 Build bản Release
Để tạo file `.exe` cho Windows:
```bash
wails build
```
File thực thi sẽ nằm trong thư mục `build/bin/`.

Để nén nhỏ file (yêu cầu UPX):
```bash
wails build -upx
```

## 🧪 Unit Testing & Makefile
Dự án đạt độ phủ code (test coverage) > 80% cho phần lõi xử lý.

```bash
# Chạy Unit Test
go test ./... -v

# Nếu có 'make' (Windows cài qua Chocolatey/Scoop hoặc dùng Git Bash)
make test
```

## 📖 Hướng dẫn sử dụng

1. **Chọn file Excel:** Chọn file nguồn chứa danh sách dữ liệu.
2. **Chọn thư mục ảnh:** Chọn thư mục chứa ảnh (hỗ trợ .jpg, .png, .webp...).
3. **Cấu hình:**
   - **Sheet Name:** Chọn Sheet cần xử lý.
   - **Cột Mã:** Cột chứa mã sản phẩm (VD: A).
   - **Cột Ảnh:** Cột đích để chèn ảnh (VD: F).
   - **Kích thước:** Điều chỉnh chiều cao dòng và độ rộng cột.
4. **Bắt đầu:** Nhấn **Start Processing** và theo dõi tiến trình.

## 📂 Cấu trúc thư mục

- `main.go`: Cấu hình cửa sổ và Wails entry.
- `app.go`: Backend logic (Go methods exposed to JS).
- `frontend/`: Mã nguồn giao diện (HTML/CSS/JS).
- `internal/engine`: Core logic xử lý Excel và Ảnh.
- `wails.json`: Cấu hình dự án Wails.
- `build/`: Thư mục chứa file thực thi sau khi build.

## 📝 Roadmap & Cải thiện
Xem thư mục `docs/` để biết thêm chi tiết.

---
Phát triển bởi [Antigravity]
Khởi tạo vào tháng 1/2026.
