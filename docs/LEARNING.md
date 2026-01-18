# 🎓 Learning & Interview Guide

Tài liệu này giúp bạn tổng hợp kiến thức từ project **ImageToExcel Importer** và chuẩn bị cho các buổi phỏng vấn Golang/Backend.

## 🧠 Core Concepts (Kiến thức Cốt lõi)

### 1. Concurrency (Đa luồng)
Project sử dụng mô hình **Worker Pool** để xử lý ảnh.
- **Tại sao dùng Worker Pool?** Để kiểm soát số lượng Goroutines chạy đồng thời. Nếu tạo 1 Goroutine cho mỗi file ảnh (ví dụ 10,000 ảnh), hệ thống sẽ bị quá tải (thrashing) và tốn RAM.
- **Channels**: Dùng `jobs channel` để gửi task và `results channel` để nhận kết quả. Đây là mô hình "Fan-out / Fan-in".
- **Synchronization**: Dùng `sync.WaitGroup` để chờ tất cả workers hoàn thành trước khi đóng channel kết quả.

### 2. Memory Management (Quản lý bộ nhớ)
- **Streaming Excel**: Sử dụng `rows.Next()` của thư viện `excelize` thay vì đọc toàn bộ file vào RAM. Điều này giúp xử lý file Excel hàng triệu dòng mà RAM vẫn ổn định.
- **Lazy Loading**: Chỉ load decode config của ảnh (`image.DecodeConfig`) để lấy kích thước trước khi load toàn bộ pixel data.

### 3. Application Architecture
- **Wails Framework**: Kết hợp sức mạnh của Go (Backend performance) và Web Tech (Frontend UI).
- **Frontend-Backend Bridge**: Giao tiếp bất đồng bộ qua JSON bridge.

---

## 🎤 Interview Questions (Câu hỏi Phỏng vấn)

Dưới đây là các câu hỏi nhà tuyển dụng có thể hỏi dựa trên project này:

### Level: Junior / Fresher

**Q1: Tại sao bạn chọn Go cho project này thay vì Python hay C#?**
> *Gợi ý:* Go có tốc độ khởi động nhanh, compile ra native binary nhỏ gọn (không cần runtime nặng như .NET/JVM), và đặc biệt là mô hình Concurrency (Goroutines) rất mạnh mẽ để xử lý I/O bound tasks (đọc/ghi file ảnh) và CPU bound tasks (nén ảnh) cùng lúc.

**Q2: Làm sao để đảm bảo thread-safe khi ghi vào file Excel?**
> *Gợi ý:* Thư viện `excelize` không an toàn tuyệt đối khi ghi song song. Trong project này, tôi dùng pattern "Single Consumer": Nhiều workers xử lý ảnh song song, nhưng kết quả được đẩy vào 1 channel duy nhất. Channel này được 1 loop (main thread) đọc và ghi vào Excel tuần tự. Điều này loại bỏ race conditions mà không cần dùng Mutex phức tạp.

**Q3: `defer` hoạt động như thế nào? Tại sao dùng `defer wg.Done()`?**
> *Gợi ý:* `defer` đẩy hàm vào stack và thực thi theo thứ tự LIFO khi hàm bao quanh return. Dùng `defer wg.Done()` đảm bảo rằng dù worker có bị panic hay return sớm ở đâu, `WaitGroup` vẫn được giảm đếm, tránh deadlock (treo chương trình mãi mãi).

### Level: Mid / Senior

**Q4: Bạn xử lý việc cập nhật UI (Progress Bar) từ Backend Go như thế nào trong Wails?**
> *Gợi ý:* Wails cung cấp cơ chế `EventsEmit`. Từ Go backend, tôi emit sự kiện kềm số % tiến độ. Frontend lắng nghe sự kiện này (`runtime.EventsOn`) và update DOM. Đây là mô hình Event-Driven, giúp decouple logic backend và giao diện.

**Q5: Nếu file Excel có 1 triệu dòng, project hiện tại có xử lý được không? Có bị OOM (Out Of Memory) không?**
> *Gợi ý:* Có thể xử lý được nhờ dùng `Iterator` (`rows.Next()`) lấy từng dòng một thay vì `GetRows()` load cả cục. Tuy nhiên, map `productMap` lưu mã sản phẩm vẫn nằm trong RAM. Với 1 triệu dòng, map này tốn khoảng vài chục đến trăm MB RAM, vẫn nằm trong giới hạn cho phép của máy tính hiện đại. Nếu cần tối ưu hơn, có thể dùng database nhẹ (SQLite/BadgerDB) thay vì map in-memory.

**Q6: Làm sao để tối ưu hóa tốc độ build Docker/CI cho project Go?**
> *Gợi ý:* Sử dụng Cache cho `go mod download` và `go build` (như đã config trong GitHub Actions `setup-go` với `cache: true`). Dùng multi-stage build trong Dockerfile (build ở stage 1, copy binary sang alpine/scratch ở stage 2) để giảm kích thước image.

**Q7: Bạn thiết kế tính năng Auto-Update như thế nào để an toàn?**
> *Gợi ý:*
> 1. Check checksum/hash của file tải về (hiện tại project chưa làm, là điểm cần cải thiện).
> 2. Sử dụng cơ chế thay thế nguyên tử (atomic replacement) hoặc script batch đệm để tránh lỗi "file đang sử dụng" trên Windows.
> 3. Versioning rõ ràng (Semantic Versioning) và inject version lúc build bằng `ldflags` để tránh hardcode sai sót.

---

## 📚 Bài tập mở rộng (Challenge)

Để nắm chắc kiến thức, hãy thử tự thực hiện các task sau:
1. **Thêm Checksum Validation**: Khi tải update về, kiểm tra mã SHA256 xem có khớp với file trên GitHub Release không.
2. **Stop/Resume**: Thêm nút "Pause" để tạm dừng worker pool và nút "Resume" để chạy tiếp.
3. **Benchmarking**: Viết benchmark so sánh tốc độ xử lý khi dùng `WorkerCount = 1` vs `WorkerCount = 10`.
