# 🔐 USBGuard

**USBGuard** là một hệ thống bảo vệ phần mềm chạy trên Windows, đảm bảo rằng ứng dụng chỉ hoạt động khi được khởi chạy từ đúng USB đã đăng ký. Dữ liệu nhận diện được mã hóa và ghi trực tiếp vào sector đặc biệt trên USB, kết hợp với kiểm tra UUID nhằm **ngăn chặn việc sao chép hoặc giả lập thiết bị**.

---

## 🚀 Tính năng chính

- ✅ Ràng buộc phần mềm với một USB cụ thể bằng UUID
- ✅ Ghi dữ liệu mã hóa vào **sector đặc biệt** (ẩn khỏi hệ thống file)
- ✅ Mã hóa dữ liệu bằng **AES** để chống đọc thô
- ✅ Obfuscate code với [`garble`](https://github.com/burrowers/garble) để chống dịch ngược
- ✅ Không phụ thuộc vào file cấu hình ngoài → không thể sao chép dễ dàng

---

## 🧱 Cơ chế hoạt động

### 1. `usb_writer.go`

- Lấy **UUID** của USB (qua PowerShell)
- Mã hóa UUID bằng AES
- Ghi nội dung mã hóa vào **sector số 2048** của ổ USB

### 2. `secure_app.go`

- Khi phần mềm khởi chạy:
  - Xác định ổ đĩa hiện tại (ví dụ: E:)
  - Đọc dữ liệu ở sector 2048
  - Giải mã và so sánh với UUID thực tế
  - Nếu hợp lệ → chạy, nếu không → thoát

---

## 📂 Cấu trúc thư mục

USBGuard/
├── usb_writer.go // Tool ghi UUID vào USB
├── secure_app.go // Phần mềm chính, kiểm tra bảo mật
├── README.md

---

## 🛠️ Hướng dẫn sử dụng

### 🔧 1. Build & chạy tool ghi UUID

```bash
go build -o usb_writer.exe usb_writer.go
usb_writer.exe

```
