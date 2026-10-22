# Sandbox Cổng thanh toán SePay

## Sandbox Cổng thanh toán SePay là môi trường thử nghiệm cách ly với Production: test toàn bộ luồng thanh toán bằng dữ liệu giả lập, không cần ngân hàng thật.

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

**Sandbox** là môi trường thử nghiệm của Cổng thanh toán SePay, cách ly hoàn toàn với Production. Bạn test toàn bộ luồng thanh toán (tạo form, chuyển sang trang thanh toán, nhận IPN, redirect) bằng dữ liệu giả lập, không chạm tiền thật và không cần liên kết ngân hàng thật. Vì vậy Sandbox phù hợp cho phát triển, kiểm thử tự động và giảng dạy tích hợp thanh toán.

## URL môi trường Sandbox

Tích hợp Sandbox dùng hai host, mỗi host một vai trò:

| Vai trò                              | Base URL Sandbox                                |
| ------------------------------------ | ----------------------------------------------- |
| Gửi form thanh toán (Checkout)       | `https://pay-sandbox.sepay.vn/v1/checkout/init` |
| REST API tra cứu đơn hàng, giao dịch | `https://pgapi-sandbox.sepay.vn`                |

Code gửi form thanh toán tới host **Checkout**; host **API** chỉ dùng khi bạn gọi REST API tra cứu đơn hàng hoặc giao dịch. REST API xác thực bằng **Basic Auth**: `Authorization: Basic base64(merchant_id:secret_key)`. Lưu ý `merchant_id` và `secret_key` của Sandbox và Production là hai bộ riêng, mỗi bộ chỉ hoạt động trên đúng môi trường của nó.

## Lấy thông tin Sandbox

Đăng ký tài khoản tại [my.sepay.vn](https://my.sepay.vn/register?onboarding=payment-gateway), sau đó kích hoạt Cổng thanh toán để nhận thông tin Sandbox, không cần chờ duyệt hồ sơ Production.

1. Vào **Cổng thanh toán** chọn **Đăng ký**, tại màn hình phương thức thanh toán chọn **Bắt đầu ngay**.

<Image src="/images/quick_start/step_1_1.png" alt="Màn hình phương thức thanh toán Cổng thanh toán SePay" caption="Màn hình phương thức thanh toán" />

2. Chọn bắt đầu với Sandbox rồi bấm **Bắt đầu hướng dẫn tích hợp**.

<Image src="/images/quick_start/step_1_8.png" alt="Bắt đầu hướng dẫn tích hợp Sandbox Cổng thanh toán SePay" caption="Bắt đầu hướng dẫn tích hợp với Sandbox" />

3. Sao chép `MERCHANT ID` và `SECRET KEY` của Sandbox để dùng cho các bước tích hợp bên dưới.

<Image src="/images/quick_start/step_1_10.png" alt="Thông tin tích hợp Sandbox gồm MERCHANT ID và SECRET KEY" caption="Thông tin tích hợp Sandbox: MERCHANT ID và SECRET KEY" />

## Trỏ tích hợp sang Sandbox

Khi dùng SDK, đặt môi trường là `sandbox`; khi dùng form HTML thuần, trỏ `action` tới host Checkout Sandbox `pay-sandbox.sepay.vn`:

<CodeTabs>
  <Code label="SDK PHP">
    ```php
    use SePay\SePayClient;
    
    // Tham số thứ ba là môi trường: 'sandbox' hoặc 'production'
    $sepay = new SePayClient('YOUR_SANDBOX_MERCHANT_ID', 'YOUR_SANDBOX_SECRET_KEY', 'sandbox');
    ```
  </Code>
  <Code label="SDK NodeJS">
    ```js
    import { SePayPgClient } from 'sepay-pg-node-sdk';
    
    const client = new SePayPgClient({
      env: 'sandbox',
      merchant_id: 'YOUR_SANDBOX_MERCHANT_ID',
      secret_key: 'YOUR_SANDBOX_SECRET_KEY'
    });
    ```
  </Code>
  <Code label="Form HTML">
    ```html
    <form method="POST" action="https://pay-sandbox.sepay.vn/v1/checkout/init">
      <!-- các trường đơn hàng và signature -->
    </form>
    ```
  </Code>
</CodeTabs>

Chi tiết tạo form, ký signature và cấu hình IPN có trong [Bắt đầu nhanh](/vi/cong-thanh-toan/bat-dau).

## Dữ liệu Sandbox cách ly với Production

* Giao dịch trong Sandbox là giả lập, không phát sinh tiền thật và không liên quan tới tài khoản ngân hàng Production.
* `merchant_id` và `secret_key` của Sandbox chỉ xác thực trên môi trường Sandbox.
* IPN và Callback URL trỏ tới hệ thống của bạn vẫn hoạt động như thật, giúp kiểm thử trọn vòng đời đơn hàng.

## Test thử trên trình duyệt

Trải nghiệm một giao dịch Sandbox hoàn chỉnh ngay trên trình duyệt mà không cần viết code:

<ButtonLink href="/vi/thanh-toan-demo" variant="primary">Trải nghiệm thanh toán thử</ButtonLink>

## Chuyển sang Production khi sẵn sàng

Sau khi test xong trên Sandbox, liên kết tài khoản ngân hàng thật, chọn **Chuyển sang Production** để nhận `MERCHANT ID` và `SECRET KEY` chính thức. Sau đó cập nhật sang base URL Production: host Checkout `pay.sepay.vn` và REST API `pgapi.sepay.vn`, cùng IPN và Callback URL Production. Các bước chi tiết nằm ở mục **Go live** trong [Bắt đầu nhanh](/vi/cong-thanh-toan/bat-dau).