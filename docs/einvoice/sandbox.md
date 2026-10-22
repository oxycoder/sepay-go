# Sandbox SePay eInvoice API

## Sandbox SePay eInvoice API là môi trường thử nghiệm cách ly với Production: test tạo và phát hành hóa đơn điện tử bằng dữ liệu giả lập trước khi go-live.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

**Sandbox** là môi trường thử nghiệm của SePay eInvoice API, cách ly hoàn toàn với Production. Bạn test trọn vòng đời hóa đơn (lấy token, tạo nháp, phát hành, kiểm tra trạng thái, tải file) bằng dữ liệu giả lập, không phát hành hóa đơn thật lên cơ quan thuế. Vì vậy Sandbox phù hợp cho phát triển, kiểm thử tự động và giảng dạy tích hợp hóa đơn điện tử.

## Base URL Sandbox

Base URL của môi trường Sandbox eInvoice API:

| Môi trường  | Base URL                                |
| ----------- | --------------------------------------- |
| **Sandbox** | `https://einvoice-api-sandbox.sepay.vn` |

Xác thực bằng Bearer token trong header: `Authorization: Bearer <ACCESS_TOKEN>`.

## Lấy token Sandbox

Gọi `POST /v1/token` trên base URL Sandbox để lấy `access_token` dùng cho các API tiếp theo. Chi tiết tham số và response có trong [Xác thực Bearer Token](/vi/einvoice-api/v1/tao-token).

## Gọi API trên Sandbox

Đặt base URL Sandbox `einvoice-api-sandbox.sepay.vn` và gửi Bearer token Sandbox trong header. Ví dụ lấy danh sách nhà cung cấp hóa đơn:

<CodeTabs>
  <Code label="cURL">
    ```bash
    curl https://einvoice-api-sandbox.sepay.vn/v1/provider-accounts \
      -H 'Authorization: Bearer YOUR_SANDBOX_TOKEN'
    ```
  </Code>
  <Code label="PHP">
    ```php
    $ch = curl_init('https://einvoice-api-sandbox.sepay.vn/v1/provider-accounts');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer YOUR_SANDBOX_TOKEN',
    ]);
    $response = curl_exec($ch);
    ```
  </Code>
  <Code label="NodeJS">
    ```js
    const res = await fetch('https://einvoice-api-sandbox.sepay.vn/v1/provider-accounts', {
      headers: { Authorization: 'Bearer YOUR_SANDBOX_TOKEN' }
    });
    const data = await res.json();
    ```
  </Code>
</CodeTabs>

Toàn bộ luồng tạo và phát hành hóa đơn có trong [Bắt đầu nhanh](/vi/einvoice-api/v1/bat-dau-nhanh).

## Test thử trên trình duyệt

Xem cấu trúc một hóa đơn điện tử và preview layout ngay trên trình duyệt mà không cần viết code:

<ButtonLink href="/vi/einvoice-demo" variant="primary">Xem hóa đơn demo</ButtonLink>

## Chuyển sang Production khi sẵn sàng

Sau khi test xong trên Sandbox, lấy token Production và đổi base URL sang `https://einvoice-api.sepay.vn`. Phần còn lại của tích hợp giữ nguyên.