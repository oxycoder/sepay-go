# Luồng thanh toán Cổng Thanh Toán SePay

## Tìm hiểu luồng thanh toán one-time và recurring trong Cổng thanh toán SePay, chọn mô hình tích hợp phù hợp cho website hoặc app của bạn.

---

**API Overview:**

API cổng thanh toán SePay hỗ trợ nhiều phương thức thanh toán bao gồm chuyển khoản ngân hàng qua QR code, NAPAS QR và thẻ quốc tế.

**Base URLs:**
- Production API: `https://pgapi.sepay.vn`
- Sandbox API: `https://pgapi-sandbox.sepay.vn`
- Production Checkout: `https://pay.sepay.vn`
- Sandbox Checkout: `https://pay-sandbox.sepay.vn`

**Xác thực:** Tất cả API sử dụng Basic Authentication với `merchant_id` và `secret_key`.


---

## Thanh toán một lần (One-time Payment)

Luồng thanh toán một lần cho phép khách hàng thanh toán ngay lập tức cho một đơn hàng cụ thể.

<Mermaid title="Trình tự xử lý thanh toán">
sequenceDiagram
  participant C as Khách hàng
  participant M as Website Merchant
  participant S as SePay Gateway
  participant B as Ngân hàng/Thẻ

  C->>M: 1. Chọn sản phẩm & thanh toán
  M->>M: 2. Tạo đơn hàng
  M->>M: 3. Tạo form checkout với signature
  M->>S: 4. POST /v1/checkout/init
  S->>S: 5. Xác thực signature
  S->>C: 6. Chuyển hướng đến trang thanh toán

  C->>S: 7. Chọn phương thức thanh toán
  S->>B: 8. Xử lý thanh toán
  B->>S: 9. Kết quả thanh toán

  alt Thanh toán thành công
      S->>M: 10a. Callback success_url
      S->>M: 11a. IPN notification
      S->>C: 12a. Chuyển hướng về trang thành công
  else Thanh toán thất bại
      S->>M: 10b. Callback error_url
 S->>C: 12b. Chuyển hướng về trang lỗi
  else Khách hàng hủy
      S->>M: 10c. Callback cancel_url
      S->>C: 11c. Chuyển hướng về trang hủy
  end

  M->>C: 13. Hiển thị kết quả cuối cùng
</Mermaid>

### Các bước chi tiết:

1. **Khách hàng chọn sản phẩm**: Thêm vào giỏ hàng và bấm "Thanh toán"
2. **Website tạo đơn hàng**: Lưu thông tin đơn hàng vào database
3. **Tạo form checkout**: Chuẩn bị dữ liệu và tạo signature HMAC-SHA256
4. **Gửi request đến SePay**: POST form đến endpoint `/v1/checkout/init`
5. **Xác thực signature**: SePay kiểm tra tính hợp lệ của chữ ký
6. **Chuyển hướng khách hàng**: Đến trang thanh toán của SePay
7. **Chọn phương thức thanh toán**: Thẻ, QR Banking, QR NAPAS
8. **Xử lý thanh toán**: SePay giao tiếp với ngân hàng/thẻ
9. **Nhận kết quả**: Từ hệ thống ngân hàng
10. **Callback về website**: Gọi các URL callback tương ứng
11. **IPN notification**: Thông báo kết quả thanh toán qua IPN
12. **Chuyển hướng khách hàng**: Về trang kết quả trên website
13. **Hiển thị kết quả**: Trang thành công/lỗi/hủy

***

## Thanh toán định kỳ (Recurring Payment)

<Callout type="info" title="Tính năng sắp ra mắt">
Thanh toán định kỳ hiện đang trong giai đoạn hoàn thiện và sẽ được phát hành trong thời gian tới.
Vui lòng theo dõi cập nhật từ SePay để biết thời điểm ra mắt chính thức.
</Callout>