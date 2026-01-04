# Cấu trúc dự án (Project Architecture)

Dự án được tổ chức theo mô hình **Clean Architecture** đơn giản hóa để dễ bảo trì và mở rộng.

## 🏗️ Sơ đồ cấu trúc

```text
ImageToExcel/
├── main.go               # Entry point, khởi tạo GUI và chạy App
├── internal/             # Code logic chính của ứng dụng
│   ├── gui/              # Xử lý giao diện người dùng (Fyne)
│   │   └── app.go        # Định nghĩa các widget, event handlers
│   └── engine/           # Xử lý logic nghiệp vụ (Core Logic)
│       ├── processor.go  # mapping Excel, worker pool, chèn ảnh
│       └── processor_test.go
├── docs/                 # Tài liệu hướng dẫn & kế hoạch
└── go.mod                # Quản lý dependencies
```

## ⚙️ Luồng xử lý chính (Main Flow)

1.  **GUI Phase**: Người dùng nhập đường dẫn Excel, thư mục ảnh và các cột cấu hình.
2.  **Mapping Phase**: `Processor` đọc cột mã sản phẩm từ Excel và tạo một `map[string]int` (Key: Mã sản phẩm, Value: Vị trí dòng).
3.  **Dispatching Phase**: Ứng dụng quét thư mục ảnh, tìm các file khớp với mã sản phẩm trong map và đẩy vào `jobs channel`.
4.  **Worker Phase**: Các Goroutines (Workers) lấy job, giải mã ảnh, lấy kích thước và chuẩn bị dữ liệu.
5.  **Collecting Phase**: Kết quả từ workers trả về `results channel`. Luồng chính đọc kết quả và gọi hàm `AddPictureFromBytes` của Excelize để chèn vào Excel (được thực hiện tuần tự để tránh xung đột file).
6.  **Finalization**: Lưu file Excel mới với hậu tố `_output.xlsx`.

## 🔒 Lưu ý về Luồng (Thread Safety)
- Vì thư viện `excelize` không đảm bảo thread-safe hoàn toàn khi ghi dữ liệu đồng thời vào cùng một file, chúng ta sử dụng cơ chế **Single Consumer** (chỉ có 1 luồng duy nhất thực hiện việc chèn ảnh vào Excel) trong khi việc đọc và giải mã ảnh được thực hiện song song bởi nhiều Workers.
