package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func stockCheckerHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy dữ liệu từ URL (ví dụ: ?productId=1&storeId=1)
	productId := r.URL.Query().Get("productId")
	storeId := r.URL.Query().Get("storeId")

	if productId == "" || storeId == "" {
		http.Error(w, "Thiếu tham số productId hoặc storeId", http.StatusBadRequest)
		return
	}

	// ⚠️ ĐÂY LÀ ĐOẠN CODE VULNERABLE (CỐ TÌNH GÂY LỖI)
	// Giả lập việc gọi một script hệ thống. Chúng ta dùng lệnh 'echo 89' để
	// đại diện cho việc stock_checker.sh trả về số lượng hàng tồn kho là 89.
	commandString := fmt.Sprintf("echo 89 %s %s", productId, storeId)

	// Gọi trực tiếp thông qua shell "sh -c". Điều này khiến Linux biên dịch
	// các ký tự đặc biệt như |, ;, &&.
	cmd := exec.Command("cmd", "/c", commandString)

	// Lấy kết quả trả về từ terminal
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(w, "Lỗi chạy lệnh: %v\nOutput: %s", err, string(output))
		return
	}

	// Trả kết quả text thô về cho trình duyệt (giống hệt bài lab)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func main() {
	http.HandleFunc("/stock", stockCheckerHandler)

	fmt.Println("🔥 Vulnerable Server đang chạy tại cổng 8080...")
	fmt.Println("👉 Truy cập bình thường: http://localhost:8080/stock?productId=1&storeId=1")

	// Khởi động server
	http.ListenAndServe(":8080", nil)
}
