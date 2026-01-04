# Kế hoạch phát triển và Cải thiện (Roadmap)

Tài liệu này ghi lại các ý tưởng và kế hoạch để nâng cấp ứng dụng **Golang Excel Image Importer**.

## 🚀 Các tính năng ưu tiên (High Priority)

1.  **[x] Hỗ trợ nhiều định dạng hơn:** Thêm hỗ trợ cho file `.webp`.
2.  **Xem trước (Preview):** Cho phép xem trước danh sách các mã sản phẩm không tìm thấy ảnh trước khi chạy.
3.  **[x] Tùy chỉnh kích thước ảnh:** Cho phép người dùng nhập kích thước ô Excel hoặc kích thước ảnh mong muốn trực tiếp từ GUI.
4.  **[x] Logging:** Xuất file log (`_missing.log`) cho các trường hợp mã sản phẩm bị thiếu.

## 🛠️ Cải tiến kỹ thuật (Technical Improvements)

- [ ] **Unit Tests:** Bổ sung thêm test case cho `internal/engine` (đặc biệt là logic mapping và scaling).
- [x] **Concurrency Tuning:** Cho phép người dùng điều chỉnh số lượng "Workers" từ giao diện để tối ưu theo cấu hình máy.
- [ ] **Memory Management:** Thử nghiệm cơ chế `AddPicture` thay vì `AddPictureFromBytes` nếu file ảnh quá lớn để giảm tải bộ nhớ đệm.
- [x] **CI/CD:** Thiết lập GitHub Actions để tự động build file `.exe` mỗi khi có release mới.

## 🎨 Giao diện (GUI Enhancements)

- Thêm chế độ Tối (Dark Mode) / Sáng (Light Mode).
- Thêm biểu tượng (Icon) cho ứng dụng.
- **[x] Cải thiện thanh tiến trình (Progress Bar):** Hiển thị chi tiết số lượng file thiếu qua thông báo kết thúc.

## 📚 Tài liệu lưu trữ
- `ARCHITECTURE.md`: Mô tả chi tiết cấu trúc code (dự kiến).
- `USER_GUIDE.md`: Hướng dẫn chi tiết cho người dùng cuối (dự kiến).
