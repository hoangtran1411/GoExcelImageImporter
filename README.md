# Golang Excel Image Importer

Một công cụ hiệu chỉnh Excel mạnh mẽ được viết bằng Go, giúp tự động chèn hình ảnh vào bảng tính dựa trên mã sản phẩm. Công cụ này được tối ưu hóa cho hiệu suất cao, xử lý hàng ngàn hình ảnh một cách nhanh chóng nhờ cơ chế Worker Pool.

## ✨ Tính năng nổi bật

- **🚀 Hiệu suất vượt trội:** Sử dụng Worker Pool để xử lý song song các tác vụ tải và nén ảnh.
- **📱 Giao diện thân thiện:** Được xây dựng với Fyne, cung cấp giao diện GUI dễ sử dụng trên Windows.
- **💾 Tối ưu bộ nhớ:** Sử dụng Iterator để đọc file Excel lớn mà không tốn nhiều RAM.
- **🔍 Tìm kiếm thông minh:** Tự động khớp tên file ảnh với mã sản phẩm trong cột Excel được chỉ định.
- **📏 Tự động căn chỉnh:** Ảnh được tự động điều chỉnh tỷ lệ để vừa vặn trong ô Excel.

## 🛠️ Công nghệ sử dụng

- **Ngôn ngữ:** [Go (Golang)](https://golang.org/)
- **Thư viện Excel:** [Excelize v2](https://github.com/xuri/excelize)
- **Framework GUI:** [Fyne v2](https://fyne.io/)

## 🚀 Hướng dấn khởi động

### Yêu cầu hệ thống
- Go 1.20 trở lên.
- Cài đặt các thư viện cần thiết cho Fyne (trên Windows yêu cầu C compiler như msys2 nếu build từ source).

### Cài đặt và Chạy
1. Clone dự án:
   ```bash
   git clone <repository-url>
   cd ImageToExcel
   ```

2. Cài đặt dependencies:
   ```bash
   go mod tidy
   ```

3. Chạy ứng dụng:
   ```bash
   go run main.go
   ```

### Build và Release tự động
Dự án đã được thiết lập **GitHub Actions**. Mỗi khi bạn `push` code lên nhánh `main`, hệ thống sẽ tự động:
- Kiểm tra lỗi (Linting/Testing).
- Build file `.exe` cho Windows 64-bit.
- Bạn có thể tải file thực thi mới nhất trong phần **Actions** của repository.

### 🔨 Build thủ công (.exe)
Để build ứng dụng mà không hiện cửa sổ console trên Windows:
```powershell
go build -ldflags="-s -w -H=windowsgui" -o ImageToExcel.exe
```

## 📖 Hướng dẫn sử dụng

1. **Chọn file Excel:** Chọn file nguồn chứa danh sách dữ liệu.
2. **Chọn thư mục ảnh:** Chọn thư mục chứa các file ảnh (định dạng .jpg, .png, .gif). Tên file phải khớp với mã sản phẩm.
3. **Cấu hình cột:**
   - **Product Code Column:** Cột chứa mã sản phẩm (Ví dụ: A, B, C...).
   - **Image Target Column:** Cột mà bạn muốn chèn ảnh vào (Ví dụ: F, G...).
4. **Bắt đầu:** Nhấn **Start Processing** và đợi quá trình hoàn tất. Kết quả sẽ được lưu thành một file mới có đuôi `_output.xlsx`.

## 📂 Cấu trúc thư mục

- `main.go`: Điểm khởi đầu của ứng dụng.
- `internal/gui`: Mã nguồn giao diện người dùng.
- `internal/engine`: Logic xử lý hình ảnh và Excel.
- `docs/`: Tài liệu hướng dẫn và kế hoạch phát triển.

## 📝 Roadmap & Cải thiện tương lai
Vui lòng xem trong thư mục `docs/` để biết thêm chi tiết về các kế hoạch nâng cấp ứng dụng.

---
Phát triển bởi [Antigravity]
Khởi tạo vào tháng 1/2026.
