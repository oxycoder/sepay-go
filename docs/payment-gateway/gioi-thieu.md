# Tích hợp cổng thanh toán SePay VietQR, NAPAS và thẻ quốc tế

## Tích hợp cổng thanh toán SePay một lần để nhận VietQR, thẻ Visa/Mastercard/JCB và NAPAS QR qua API, SDK PHP/Node.js.

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

## Cổng thanh toán SePay là gì?

**Cổng thanh toán SePay** cho phép doanh nghiệp và lập trình viên tích hợp nhiều phương thức thanh toán vào website hoặc ứng dụng chỉ với một lần tích hợp, không cần kết nối riêng với từng ngân hàng hay tổ chức thẻ.

SePay đóng vai trò trung gian giữa website/ứng dụng của bạn với hệ thống ngân hàng, mạng thẻ quốc tế và cổng NAPAS, xử lý toàn bộ luồng thanh toán từ khởi tạo đơn hàng đến xác nhận giao dịch.

<Callout type="info" title="Thời gian kích hoạt theo từng phương thức">
QR chuyển khoản ngân hàng (VietQR)
: kích hoạt online trong 
~30 phút
, nhận thanh toán ngay sau khi hoàn tất cài đặt.
Cổng NAPAS (VietQRPay) & thẻ tín dụng/ghi nợ quốc tế
: thời gian duyệt 
3-5 ngày đến 2 tuần
 tùy hồ sơ đơn vị bán hàng, liên hệ bộ phận kinh doanh để biết chi tiết.
Trong thời gian chờ duyệt Production, bạn có thể đăng ký và 
thử ngay trên môi trường Sandbox
 để hoàn thiện tích hợp. Đăng ký tại 
my.sepay.vn
.
</Callout>

<ButtonLink href="/vi/thanh-toan-demo" variant="primary">Trải nghiệm thanh toán thử</ButtonLink>

***

## Phương thức thanh toán được hỗ trợ

Cổng thanh toán SePay hỗ trợ ba nhóm phương thức thanh toán chính:

<FeatureGrid
  iconPosition="inline"
  items={[
  { icon: "zap", title: "QR Code chuyển khoản ngân hàng (VietQR)", description: "Khách hàng quét mã QR VietQR bằng ứng dụng ngân hàng để chuyển tiền trực tiếp vào tài khoản của bạn. Tiền về tài khoản ngân hàng mà không qua bất kỳ ví trung gian nào. Hỗ trợ tài khoản cá nhân và doanh nghiệp, xác nhận giao dịch tự động theo thời gian thực, phí cố định theo giao dịch." },
  { icon: "shield", title: "QR Code cổng NAPAS (VietQRPay)", description: "Cổng thanh toán QR theo chuẩn NAPAS, SePay là đối tác ủy quyền chính thức của NAPAS. Sử dụng công nghệ VietQRPay và VietQRGlobal. Tiền về ngân hàng theo thời gian thực, dòng tiền nguyên vẹn, dễ đối soát." },
  { icon: "lock", title: "Thẻ tín dụng/ghi nợ quốc tế", description: "Chấp nhận thanh toán bằng thẻ Visa, Mastercard và JCB. Sản phẩm hợp tác giữa SePay, VPBank và MPGS của Mastercard. Hỗ trợ thẻ tín dụng và thẻ ghi nợ quốc tế, bảo mật theo chuẩn 3D Secure." }
]}
/>

***

## Tính năng chính

Cổng thanh toán SePay được xây dựng để giải quyết các bài toán thanh toán phổ biến của doanh nghiệp Việt Nam:

<FeatureGrid
  iconPosition="inline"
  items={[
  { icon: "zap", title: "Chấp nhận thanh toán trực tuyến", description: "Tích hợp một lần duy nhất qua API hoặc SDK để tiếp nhận thanh toán từ nhiều phương thức khác nhau. Không cần ký hợp đồng riêng với từng ngân hàng." },
  { icon: "refresh", title: "Thông báo giao dịch tức thì (IPN)", description: "Khi có giao dịch thành công, SePay tự động gửi thông báo đến endpoint IPN của bạn theo thời gian thực, giúp hệ thống tự động xác nhận đơn hàng mà không cần người dùng thao tác thêm." },
  { icon: "link", title: "Chuyển hướng sau thanh toán", description: "Tùy chỉnh URL để điều hướng khách hàng sau khi thanh toán thành công, thất bại hoặc hủy bỏ, phù hợp với luồng trải nghiệm của từng hệ thống." },
  { icon: "code", title: "Webhook & tự lập trình cổng thanh toán", description: "SePay cung cấp webhook thông báo biến động số dư, cho phép lập trình viên tự xây dựng cổng thanh toán tùy chỉnh cho website mà không phụ thuộc vào form thanh toán mặc định." },
  { icon: "sliders", title: "Môi trường Sandbox", description: "Toàn bộ tính năng đều có thể kiểm thử trên môi trường Sandbox trước khi đưa vào production, đảm bảo tích hợp chính xác mà không ảnh hưởng đến giao dịch thật." },
  { icon: "trending-up", title: "Kích hoạt online", description: "Đăng ký và kích hoạt cổng thanh toán hoàn toàn online, không cần giấy tờ phức tạp. Với VietQR, bạn có thể bắt đầu nhận tiền chỉ sau khoảng 30 phút." },
  { icon: "shield", title: "Bảo mật theo chuẩn quốc tế", description: "Thanh toán thẻ tuân thủ chuẩn 3D Secure (3DS) và PCI DSS thông qua đối tác MPGS của Mastercard. Mọi giao dịch đều được mã hoá và xác thực trước khi xử lý." },
  { icon: "lock", title: "Đối tác chính thức của NAPAS", description: "SePay là đối tác ủy quyền của NAPAS, đảm bảo tuân thủ chuẩn kỹ thuật và pháp lý của hệ thống thanh toán quốc gia." }
]}
/>

***

## Cách tích hợp cổng thanh toán SePay

Cổng thanh toán SePay phục vụ nhiều nhu cầu tích hợp, từ lập trình viên tự build bằng API đến người dùng không code:

<FeatureGrid
  items={[
  { icon: "code", title: "Tự build với API", description: "Lập trình viên tích hợp trực tiếp qua RESTful API của SePay, toàn quyền tùy biến giao diện thanh toán, phù hợp với hệ thống tùy chỉnh, ứng dụng di động hoặc kiến trúc SaaS multi-tenant." },
  { icon: "package", title: "Dùng SDK PHP / Node.js", description: "SePay cung cấp SDK chính thức cho PHP và Node.js, giúp rút ngắn thời gian tích hợp. Phù hợp với đội backend muốn thư viện sẵn sàng thay vì tự xử lý chữ ký, gửi lại và bắt lỗi." },
  { icon: "store", title: "Dùng plugin CMS (WooCommerce, Haravan)", description: "Cửa hàng TMĐT chạy trên WordPress/WooCommerce, Haravan có thể cài plugin chính thức và cấu hình merchant_id/secret_key, không cần viết code." },
  { icon: "users", title: "Không code → Payment link", description: "Cá nhân, cửa hàng nhỏ hoặc freelancer không muốn code có thể tạo QR payment link trực tiếp tại <a href=\"https://my.sepay.vn\" target=\"_blank\" rel=\"noopener\">my.sepay.vn</a> và gửi link cho khách hàng để nhận thanh toán." }
]}
/>

<Callout type="tip" title="Không biết nên chọn phương thức nào?">
Nếu bạn là cá nhân hoặc doanh nghiệp nhỏ, hãy bắt đầu với 
QR chuyển khoản ngân hàng (VietQR)
: kích hoạt nhanh nhất, phí thấp nhất. Thanh toán thẻ quốc tế và cổng NAPAS phù hợp hơn cho doanh nghiệp có nhu cầu lớn hoặc phục vụ khách hàng quốc tế.
</Callout>

***

## Luồng thanh toán

Ở cấp độ nghiệp vụ, một giao dịch qua SePay trải qua **5 bước chính**: từ khi khách hàng bấm "Thanh toán" cho đến khi đơn vị bán hàng ghi nhận đơn hàng thành công:

<Mermaid title="Luồng thanh toán SePay (nghiệp vụ)">
flowchart LR
  A[Khách hàng<br/>chọn thanh toán] --> B[Đơn vị bán hàng<br/>tạo đơn hàng]
  B --> C[Khách hàng<br/>thực hiện thanh toán]
  C --> D[Ngân hàng / NAPAS<br/>xác nhận]
  D --> E[Đơn vị bán hàng cập nhật<br/>đơn hàng qua IPN]
</Mermaid>

1. **Khách hàng chọn thanh toán**: Trên website/app của đơn vị bán hàng, khách hàng chọn sản phẩm và bấm nút thanh toán.

2. **Đơn vị bán hàng tạo đơn hàng**: Hệ thống của đơn vị bán hàng khởi tạo đơn hàng trên SePay với số tiền, mô tả và URL callback; nhận lại link thanh toán hoặc mã QR.

3. **Khách hàng thực hiện thanh toán**: Khách hàng chọn phương thức (QR VietQR, QR NAPAS, hoặc thẻ quốc tế) và hoàn tất giao dịch.

4. **Ngân hàng / NAPAS xác nhận**: Giao dịch được xử lý và xác nhận bởi ngân hàng phát hành hoặc cổng NAPAS.

5. **Đơn vị bán hàng cập nhật đơn hàng qua IPN**: SePay gửi IPN webhook đến endpoint của đơn vị bán hàng để cập nhật trạng thái đơn hàng tự động. Khách hàng được chuyển hướng về `return_url` để xem kết quả.

<Callout type="info" title="Xem sequence API chi tiết">
Tài liệu này mô tả luồng ở cấp độ nghiệp vụ. Để xem 
sequence diagram các API call cụ thể
 (endpoint, payload, chữ ký, flow 
`checkout/init`
), vui lòng xem 
Tạo đơn hàng, Tạo form thanh toán
.
</Callout>

<Callout type="info" title="Xác nhận đơn hàng qua IPN">
Không nên chỉ dựa vào 
`return_url`
 để xác nhận đơn hàng vì khách hàng có thể đóng trình duyệt trước khi được chuyển hướng. Luôn dùng 
IPN Webhook
 để cập nhật trạng thái đơn hàng phía server.
</Callout>

***

## Môi trường

| Môi trường     | Base URL                         |
| -------------- | -------------------------------- |
| **Production** | `https://pgapi.sepay.vn`         |
| **Sandbox**    | `https://pgapi-sandbox.sepay.vn` |

Xác thực bằng **Basic Auth**: `Authorization: Basic base64(merchant_id:secret_key)`

Nếu bạn mới bắt đầu, hãy thử nghiệm trên Sandbox trước khi chuyển sang Production. Cách lấy thông tin Sandbox, trỏ tích hợp và dữ liệu cách ly có trong [Sandbox Cổng Thanh Toán SePay](/vi/cong-thanh-toan/sandbox).

***

## Bước tiếp theo

Sau khi đã hiểu về Cổng thanh toán SePay:

1. **[Bắt đầu tích hợp](/vi/cong-thanh-toan/bat-dau)**: Hướng dẫn từng bước để tích hợp cổng thanh toán vào website hoặc ứng dụng của bạn
2. **[Xác thực API](/vi/cong-thanh-toan/API/tong-quan)**: Tìm hiểu cách xác thực với API SePay và các base URL cho Sandbox và Production

<Callout type="info" title="Bắt đầu với Sandbox">
Hãy bắt đầu bằng môi trường Sandbox để kiểm thử tích hợp mà không ảnh hưởng đến giao dịch thật, không cần ngân hàng thật. Xem 
Sandbox Cổng Thanh Toán SePay
.
</Callout>