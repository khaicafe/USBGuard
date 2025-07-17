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

```bash
garble -literals -tiny build -o secure_app.exe secure_app.go
GOOS=windows GOARCH=amd64 garble -literals -tiny build -o secure_app.exe secure_app.go

```

<!-- sercurity -->

```bash
✅ 1. Dữ liệu ghi vào sector có hiện ra như file bình thường không?
Không.
Khi bạn ghi vào sector bằng cách dùng \\.\E: và Seek(sector \* 512) như bạn đang làm:

Bạn ghi trực tiếp vào vùng raw disk (sector vật lý)

Không qua hệ thống file (NTFS, FAT32...)

❌ Không hiện file nào trong File Explorer

❌ Không hiện khi dùng dir trong cmd

❌ Không bị xoá khi người dùng “Format USB”

➡️ Người dùng bình thường không hề biết có dữ liệu ở đó.

⚠️ 2. Có đọc được dữ liệu không?
Có – nếu hacker dùng công cụ phân tích raw disk, ví dụ:

Tool đọc sector Mục đích
HxD (Windows) Đọc và chỉnh sector trực tiếp
WinHex Xem ổ đĩa ở cấp byte
dd (Linux) Copy sector từ ổ đĩa
DiskGenius Phục hồi, xem partition, đọc sector

→ Nếu dữ liệu bạn ghi không mã hóa, thì có thể bị đọc rõ ràng.

🔒 3. Nếu bạn đã mã hóa bằng AES-GCM như trong code, thì sao?
Dữ liệu sẽ trông như chuỗi byte rác (hex hoặc binary)

Không có chuỗi UUID hay số serial nào hiện ra

Dù hacker copy toàn bộ sector → vẫn không giải mã được nếu không có AES key chính xác

→ Vô nghĩa với hacker nếu không biết key + thuật toán

✅ 4. Gợi ý bảo mật cao nhất bạn đang dùng:
Thành phần Bảo vệ ra sao?
Ghi dữ liệu vào sector 2048 ✅ Không hiện file
Dữ liệu mã hóa bằng AES-GCM ✅ Không thể đọc hiểu
So sánh UUID + Serial ✅ Không giả lập được
Code dùng garble obfuscate ✅ Khó patch bypass

➡️ Bạn đang áp dụng chuẩn “software dongle” mạnh mẽ nhất bằng USB vật lý.

```
