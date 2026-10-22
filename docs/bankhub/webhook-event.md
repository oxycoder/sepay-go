# Tham chiếu Webhook Events

## Tham chiếu đầy đủ 14 webhook event của SePay Bank Hub gồm liên kết, hủy liên kết, khôi phục tài khoản, trạng thái phiên cùng schema payload chi tiết và quy tắc retry.

---

**API Overview:**

API tích hợp Bank Hub - Nền tảng kết nối tài khoản ngân hàng và nhận thông báo giao dịch realtime.

## Base URL
- **Sandbox**: `https://bankhub-api-sandbox.sepay.vn`
- **Production**: `https://bankhub-api.sepay.vn`

## Authentication
- Sử dụng **Basic Authentication** với `client_id:client_secret` để lấy access token
- Sử dụng **Bearer Token** cho các API khác


---

<Callout type="info" title="Cấu hình Webhook">
Để nhận webhook events, bạn cần cung cấp webhook URL cho SePay Bank Hub. 
API thêm webhook
</Callout>

<Callout type="info" title="Lưu ý">
Đây không phải webhook để nhận biến động số dư. Để nhận thông báo biến động số dư dùng IPN 
Xem chi tiết IPN
</Callout>

***

## Cấu trúc chung Webhook Event

Mỗi webhook event có cấu trúc chuẩn như sau:

<Response title="Webhook Event">
```json
{
  "timestamp": 1768312100,
  "xid": "e61cdae9-505e-4349-9713-a3e5b8b80b5e",
  "version": "1.0",
  "event": "EVENT_NAME",
  "metadata": {

  }
}
```
</Response>

**Các trường chung:**

<Params
  rows={[
{ name: "timestamp", type: "integer", required: true, description: "Unix timestamp (giây) khi event được tạo" },
{ name: "xid", type: "string", required: true, description: "UUID duy nhất của event" },
{ name: "version", type: "string", required: true, description: "Phiên bản API webhook (hiện tại: \"1.0\")" },
{ name: "event", type: "string", required: true, description: "Tên event (ví dụ: LINK_TOKEN_CREATED, BANK_ACCOUNT_LINKED)" },
{ name: "metadata", type: "object", required: true, description: "Dữ liệu chi tiết của event, khác nhau tùy loại event" }
]}
/>

## Luồng liên kết ngân hàng

* Khi vừa tạo link token cho luồng liên kết ngân hàng

<Response title="LINK_TOKEN_CREATED">
```json
{
  "timestamp": 1768312100,
  "xid": "e61cdae9-505e-4349-9713-a3e5b8b80b5e",
  "version": "1.0",
  "event": "LINK_TOKEN_CREATED",
  "metadata": {
    "purpose": "LINK_BANK_ACCOUNT",
    "link_token_xid": "628af450-6f83-443f-b2f1-7d8ae4664e45",
    "link_session_xid": null
  }
}
```
</Response>

* Khi người dùng bắt đầu phiên liên kết và chuẩn bị chọn ngân hàng.

<Response title="LINK_SESSION_INITIALIZED">
```json
{
  "timestamp": 1768311902,
  "xid": "d5df6e97-689e-49ae-b380-08126c6da950",
  "version": "1.0",
  "event": "LINK_SESSION_INITIALIZED",
  "metadata": {
    "link_session_xid": "0e296f70-cf41-4e8f-a89f-00cee248b887",
    "link_token_xid": "ef4cb765-c297-4a51-a9eb-6f7a751978ae"
  }
}
```
</Response>

* Khi người dùng đã chọn ngân hàng và chuyển đến màn hình nhập thông tin.

<Response title="LINK_SESSION_STATE_CHANGED - REQUEST_BANK_ACCOUNT_LINK">
```json
{
  "timestamp": 1768312140,
  "xid": "e117249c-a827-4f89-951d-433c3a384785",
  "version": "1.0",
  "event": "LINK_SESSION_STATE_CHANGED",
  "metadata": {
    "state": "REQUEST_BANK_ACCOUNT_LINK",
    "brand_name": "MBBank",
    "account_type": "personal",
    "link_token_xid": "628af450-6f83-443f-b2f1-7d8ae4664e45",
    "link_session_xid": "a9e6d327-5888-449a-b82a-356fa80d6211"
  }
}
```
</Response>

* Khi người dùng đã nhập thông tin tài khoản và chuyển đến màn hình nhập OTP.

<Response title="LINK_SESSION_STATE_CHANGED">
```json
{
  "timestamp": 1768311997,
  "xid": "a67f03a1-1bdf-46a5-8c6e-4af6aec28aca",
  "version": "1.0",
  "event": "LINK_SESSION_STATE_CHANGED",
  "metadata": {
    "state": "CONFIRM_BANK_ACCOUNT_LINK",
    "brand_name": "MBBank",
    "account_type": "personal",
    "account_number": "0000000000011111",
    "account_holder_name": "NGUYEN VAN A",
    "link_token_xid": "ef4cb765-c297-4a51-a9eb-6f7a751978ae",
    "link_session_xid": "0e296f70-cf41-4e8f-a89f-00cee248b887"
  }
}
```
</Response>

* Khi quá trình liên kết hoàn tất thành công sẽ gửi 3 webhook events:

<Response title="LINK_SESSION_STATE_CHANGED - FINISHED_BANK_ACCOUNT_LINK">
```json
{
  "timestamp": 1768312172,
  "xid": "21e2445d-9696-4edd-816c-746b2e6ec110",
  "version": "1.0",
  "event": "LINK_SESSION_STATE_CHANGED",
  "metadata": {
    "state": "FINISHED_BANK_ACCOUNT_LINK",
    "brand_name": "MBBank",
    "account_type": "personal",
    "account_number": "0000000000011111",
    "account_holder_name": "NGUYEN VAN A",
    "link_token_xid": "628af450-6f83-443f-b2f1-7d8ae4664e45",
    "link_session_xid": "a9e6d327-5888-449a-b82a-356fa80d6211"
  }
}
```
</Response>

<Response title="BANK_ACCOUNT_LINKED">
```json
{
  "timestamp": 1768312023,
  "xid": "e343bdec-97dc-4c60-b953-f81b183e0096",
  "version": "1.0",
  "event": "BANK_ACCOUNT_LINKED",
  "metadata": {
    "bank_account_xid": "587eabfb-f086-11f0-b16e-52c7e9b4f41b",
    "brand_name": "MBBank",
    "account_type": "personal",
    "account_number": "0000000000011111",
    "link_token_xid": "ef4cb765-c297-4a51-a9eb-6f7a751978ae",
    "link_session_xid": "0e296f70-cf41-4e8f-a89f-00cee248b887"
  }
}
```
</Response>

<Response title="LINK_SESSION_COMPLETED">
```json
{
  "timestamp": 1768312173,
  "xid": "19b3bbb2-41a3-4b27-a301-db91a036fb18",
  "version": "1.0",
  "event": "LINK_SESSION_COMPLETED",
  "metadata": {
    "bank_account_xid": "b1293181-f086-11f0-b16e-52c7e9b4f41b",
    "link_token_xid": "628af450-6f83-443f-b2f1-7d8ae4664e45",
    "link_session_xid": "a9e6d327-5888-449a-b82a-356fa80d6211"
  }
}
```
</Response>

***

## Luồng huỷ liên kết

* Khi vừa tạo thành công link token cho luồng huỷ liên kết ngân hàng:

<Response title="LINK_TOKEN_CREATED">
```json
{
  "timestamp": 1768443727,
  "xid": "3406b960-71e4-4fdd-9bb3-41bba8fe6107",
  "version": "1.0",
  "event": "LINK_TOKEN_CREATED",
  "metadata": {
    "purpose": "UNLINK_BANK_ACCOUNT",
    "link_token_xid": "378c3028-a366-4aaa-90d4-56e0047d92ce",
    "link_session_xid": null
  }
}
```
</Response>

* Khi người dùng bắt đầu huỷ liên kết và chuyển đến màn hình xác nhận (OTP).

<Response title="LINK_SESSION_STATE_CHANGED">
```json
{
  "timestamp": 1768443767,
  "xid": "a3d04fca-226c-4376-ac31-189283801f1b",
  "version": "1.0",
  "event": "LINK_SESSION_STATE_CHANGED",
  "metadata": {
    "state": "REQUEST_BANK_ACCOUNT_UNLINK",
    "brand_name": "ACB",
    "account_type": "individual",
    "account_number": "7777777777",
    "link_token_xid": "378c3028-a366-4aaa-90d4-56e0047d92ce",
    "link_session_xid": "53f42c84-ada6-442a-8732-6d8ec686cdc0"
  }
}
```
</Response>

* Khi tài khoản ngân hàng được huỷ liên kết thành công (huỷ liên kết đến ngân hàng)

<Response title="LINK_SESSION_COMPLETED">
```json
{
  "timestamp": 1768443806,
  "xid": "76e32fa6-5ed2-42b6-9f57-7cb8acbd12dc",
  "version": "1.0",
  "event": "LINK_SESSION_COMPLETED",
  "metadata": {
    "bank_account_xid": "54a488a6-f1b8-11f0-b21a-a6006ab65aca",
    "unlinked": true,
    "link_token_xid": "378c3028-a366-4aaa-90d4-56e0047d92ce",
    "link_session_xid": "53f42c84-ada6-442a-8732-6d8ec686cdc0"
  }
}
```
</Response>

<Response title="BANK_ACCOUNT_UNLINKED | BANK_ACCOUNT_INACTIVED">
```json
{
  "timestamp": 1768443805,
  "xid": "3c02781a-e061-44a0-b035-0cab9d00f73c",
  "version": "1.0",
  "event": "BANK_ACCOUNT_UNLINKED | BANK_ACCOUNT_INACTIVED",
  "metadata": {
    "bank_account_xid": "54a488a6-f1b8-11f0-b21a-a6006ab65aca",
    "brand_name": "ACB",
    "account_type": "individual",
    "account_number": "7777777777",
    "link_token_xid": "378c3028-a366-4aaa-90d4-56e0047d92ce",
    "link_session_xid": "53f42c84-ada6-442a-8732-6d8ec686cdc0"
  }
}
```
</Response>

<Callout type="warning" title="Lưu ý">
BANK_ACCOUNT_UNLINKED:
 Huỷ liên kết trực tiếp với phía ngân hàng
BANK_ACCOUNT_INACTIVED:
 Chỉ vô hiệu hoá mà không huỷ trực tiếp với ngân hàng
</Callout>

***

## Luồng khôi phục tài khoản

Purpose `REACTIVATE_BANK_ACCOUNT` dùng để bật lại tài khoản đang ngừng kích hoạt (trước đó bị hủy liên kết kiểu vô hiệu hóa). Luồng phát webhook đối xứng luồng liên kết/hủy liên kết.

* Khi vừa tạo thành công link token cho luồng khôi phục:

<Response title="LINK_TOKEN_CREATED">
```json
{
  "timestamp": 1768520100,
  "xid": "b7e4d2f0-1c3a-4e58-9d21-6f0a2b8c4e11",
  "version": "1.0",
  "event": "LINK_TOKEN_CREATED",
  "metadata": {
    "purpose": "REACTIVATE_BANK_ACCOUNT",
    "link_token_xid": "9a1f7c62-5b84-4d0e-a3f1-2c7e9b4f41bd",
    "link_session_xid": null
  }
}
```
</Response>

* Khi người dùng mở trang khôi phục và phiên được khởi tạo:

<Response title="LINK_SESSION_INITIALIZED">
```json
{
  "timestamp": 1768520140,
  "xid": "c1a8e6b3-7d42-4f19-8e05-3a9c1d2b6f74",
  "version": "1.0",
  "event": "LINK_SESSION_INITIALIZED",
  "metadata": {
    "link_session_xid": "4d7b0a91-2e63-4c58-b9f2-8a1c6e5d3f20",
    "link_token_xid": "9a1f7c62-5b84-4d0e-a3f1-2c7e9b4f41bd"
  }
}
```
</Response>

* Khi tài khoản được khôi phục thành công sẽ gửi 3 webhook events:

<Response title="LINK_SESSION_STATE_CHANGED - FINISHED_BANK_ACCOUNT_REACTIVATE">
```json
{
  "timestamp": 1768520180,
  "xid": "e2c9f4a7-3b81-4d60-9c15-7f0a2e8b1d43",
  "version": "1.0",
  "event": "LINK_SESSION_STATE_CHANGED",
  "metadata": {
    "state": "FINISHED_BANK_ACCOUNT_REACTIVATE",
    "brand_name": "OCB",
    "account_type": "individual",
    "link_token_xid": "9a1f7c62-5b84-4d0e-a3f1-2c7e9b4f41bd",
    "link_session_xid": "4d7b0a91-2e63-4c58-b9f2-8a1c6e5d3f20"
  }
}
```
</Response>

<Response title="BANK_ACCOUNT_REACTIVATED">
```json
{
  "timestamp": 1768520181,
  "xid": "f3d0a5b8-4c92-4e71-8d26-1a2b3c4d5e6f",
  "version": "1.0",
  "event": "BANK_ACCOUNT_REACTIVATED",
  "metadata": {
    "bank_account_xid": "54a488a6-f1b8-11f0-b21a-a6006ab65aca",
    "brand_name": "OCB",
    "account_type": "individual",
    "link_token_xid": "9a1f7c62-5b84-4d0e-a3f1-2c7e9b4f41bd",
    "link_session_xid": "4d7b0a91-2e63-4c58-b9f2-8a1c6e5d3f20"
  }
}
```
</Response>

<Response title="LINK_SESSION_COMPLETED">
```json
{
  "timestamp": 1768520182,
  "xid": "a4e1b6c9-5d03-4f82-9e37-2b3c4d5e6f70",
  "version": "1.0",
  "event": "LINK_SESSION_COMPLETED",
  "metadata": {
    "bank_account_xid": "54a488a6-f1b8-11f0-b21a-a6006ab65aca",
    "reactivated": true,
    "link_token_xid": "9a1f7c62-5b84-4d0e-a3f1-2c7e9b4f41bd",
    "link_session_xid": "4d7b0a91-2e63-4c58-b9f2-8a1c6e5d3f20"
  }
}
```
</Response>

<Callout type="warning" title="Lưu ý khi khôi phục tài khoản">
Tài khoản định danh (VA) đi kèm được khôi phục cùng tài khoản chính (áp dụng cho các ngân hàng có VA: OCB, KienLongBank, MSB; TPBank không có VA).
Giao dịch phát sinh trong thời gian tài khoản ngừng kích hoạt 
không được đồng bộ bù
 sau khi khôi phục.
</Callout>

***

## Gọi lại Webhook

SePay sẽ gọi lại webhooks nếu trạng thái kết nối mạng đến webhook url thất bại. Ngoài ra, bạn có thể tùy chọn các điều kiện SePay hỗ trợ sẵn để có thể gọi lại webhooks. Thời gian gọi cách nhau bằng phút, tăng dần theo dãy số **[Fibonacci](https://en.wikipedia.org/wiki/Fibonacci_sequence)**

* Số lần gọi lại tối đa là 7 lần
* Tối đa là 5 giờ kể từ khi gọi lần đầu thất bại
* Network connect timeout của SePay là 5 giây
* Thời gian chờ phản hồi tối đa của SePay là 8 giây

<Callout type="info" title="Lưu ý">
SePay sẽ 
KHÔNG
 gọi lại webhooks nếu trạng thái là thất bại nhưng kết nối mạng là thành công, trừ khi webhook đó được thiết lập điều kiện gọi lại và chỉ gọi lại khi thỏa mãn một trong số các điều kiện được thiết lập.
</Callout>

***

## Bước tiếp theo

Sau khi cấu hình webhook events, bạn có thể:

1. **[Thông báo biến động số dư (IPN)](/vi/bankhub/thong-bao-bien-dong-so-du)** - Nhận thông báo giao dịch realtime
2. **[API Cập nhật Webhook](/vi/bankhub/api/api-webhook/cap-nhat-webhook)** - Cập nhật URL webhook qua API
3. **[API Danh sách tài khoản](/vi/bankhub/api/tai-khoan-ngan-hang/danh-sach-tai-khoan)** - Xem danh sách tài khoản đã liên kết
4. **[API Danh sách giao dịch](/vi/bankhub/api/api-giao-dich/danh-sach-giao-dich)** - Truy vấn lịch sử giao dịch