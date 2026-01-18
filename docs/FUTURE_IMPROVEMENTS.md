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
- [x] **Wails Migration:** Chuyển đổi từ Fyne sang Wails để có giao diện đẹp và nhẹ hơn.
- [x] **Auto Update:** Tích hợp cơ chế tự động cập nhật qua GitHub Releases.
- [x] **CI/CD:** Thiết lập GitHub Actions để tự động build file `.exe` mỗi khi có release mới.

## 🎨 Giao diện (GUI Enhancements)

- **[x] Dark Mode:** Giao diện tối hiện đại, dễ nhìn.
- **[x] Biểu tượng (Icon):** Ứng dụng có icon riêng và giao diện chuyên nghiệp.
- **[x] Cải thiện thanh tiến trình (Progress Bar):** Hiển thị chi tiết phần trăm và trạng thái.
- **[x] Responsive UI:** Giao diện tự động co giãn phù hợp kích thước cửa sổ.

## 📚 Tài liệu lưu trữ
- `ARCHITECTURE.md`: Mô tả chi tiết cấu trúc code Wails.
- `SKILL.md`: Hướng dẫn tái sử dụng kỹ năng Wails.
