# Xuất hóa đơn điện tử cho hộ kinh doanh, doanh nghiệp

## Xuất hóa đơn điện tử cho hộ kinh doanh và doanh nghiệp qua SePay eInvoice API: tạo, phát hành, tra cứu và tải hóa đơn theo Nghị định 123/2020/NĐ-CP.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

## SePay eInvoice API là gì?

**SePay eInvoice API** giúp hộ kinh doanh và doanh nghiệp xuất hóa đơn điện tử tuân thủ **Nghị định 123/2020/NĐ-CP** và Thông tư 78/2021/TT-BTC của Bộ Tài chính. Qua một giao diện API RESTful duy nhất, hệ thống của bạn có thể tạo hóa đơn, phát hành lên cơ quan thuế (CQT) và quản lý toàn bộ vòng đời hóa đơn một cách tự động.

API cung cấp đầy đủ vòng đời hóa đơn: từ tạo nháp, phát hành lên CQT, kiểm tra trạng thái theo thời gian thực, cho đến tải file PDF/XML và tra cứu hạn mức sử dụng, tất cả thông qua giao diện RESTful HTTP chuẩn.

<ButtonLink href="/vi/einvoice-demo" variant="primary">Xem hóa đơn demo</ButtonLink>

***

## Đối tượng sử dụng

eInvoice API phù hợp cho các tổ chức và cá nhân cần tự động hóa quy trình xuất hóa đơn điện tử:

<Features
  items={[
  { icon: "home", title: "Hộ kinh doanh & doanh nghiệp", description: "Xuất hóa đơn điện tử cho từng giao dịch bán hàng, đáp ứng yêu cầu hóa đơn theo Nghị định 123/2020/NĐ-CP." },
  { icon: "package", title: "Nền tảng SaaS & ERP", description: "Tích hợp phát hành hóa đơn điện tử trực tiếp vào quy trình bán hàng, kế toán hoặc quản lý đơn hàng của phần mềm." },
  { icon: "trending-up", title: "Doanh nghiệp xuất hóa đơn số lượng lớn", description: "Tự động hóa hoàn toàn việc tạo và phát hành hóa đơn theo lô, giảm thiểu thao tác thủ công và sai sót." },
  { icon: "sliders", title: "Phần mềm kế toán", description: "Kết nối trực tiếp với CQT để phát hành và đồng bộ trạng thái hóa đơn, hỗ trợ đối soát và lưu trữ theo quy định." },
  { icon: "store", title: "Sàn thương mại điện tử", description: "Tự động xuất hóa đơn cho từng đơn hàng hoàn thành, đáp ứng yêu cầu xuất hóa đơn theo giao dịch." },
  { icon: "code", title: "Developer xây dựng giải pháp", description: "Tích hợp tính năng phát hành hóa đơn điện tử vào sản phẩm phục vụ khách hàng doanh nghiệp." }
]}
/>

***

## Tính năng chính

<Features
  items={[
  { icon: "zap", title: "Xuất hóa đơn điện tử", description: "Tạo hóa đơn nháp hoặc phát hành trực tiếp qua POST v1/invoices/create, hỗ trợ nhiều mẫu và ký hiệu hóa đơn." },
  { icon: "shield", title: "Phát hành lên cơ quan thuế", description: "Gửi hóa đơn đến CQT qua nhà cung cấp hóa đơn điện tử đã cấu hình bằng POST v1/invoices/issue." },
  { icon: "zap", title: "Xuất hóa đơn điều chỉnh", description: "Điều chỉnh hóa đơn đã phát hành thông qua API POST v1/invoices/adjust, bao gồm điều chỉnh tăng, giảm hoặc điều chỉnh thông tin hóa đơn." },
  { icon: "zap", title: "Xuất hóa đơn thay thế", description: "Xuất hóa đơn mới thay thế cho hóa đơn đã phát hành có sai sót thông qua API POST v1/invoices/replace. Hóa đơn gốc sẽ được đánh dấu là đã bị thay thế." },
  { icon: "refresh", title: "Theo dõi trạng thái bất đồng bộ", description: "Kiểm tra kết quả xử lý create/issue theo thời gian thực thông qua API /check/{tracking_code}." },
  { icon: "database", title: "Tra cứu và tải hóa đơn", description: "Lấy chi tiết hóa đơn đã phát hành kèm URL tải file PDF, XML theo chuẩn CQT." },
  { icon: "clock", title: "Kiểm tra hạn mức sử dụng", description: "Theo dõi số lượt phát hành còn lại qua GET v1/usage để lên kế hoạch sử dụng dịch vụ." },
  { icon: "sliders", title: "Quản lý tài khoản nhà cung cấp", description: "Xem danh sách và cấu hình chi tiết của các tài khoản nhà cung cấp hóa đơn điện tử đã đăng ký." },
  { icon: "store", title: "Địa điểm kinh doanh người bán", description: "Khai báo danh mục địa điểm (mã số, tên, địa chỉ) và gắn seller_store_xid khi tạo hóa đơn để truyền mã số và tên địa điểm lên hóa đơn." }
]}
/>

***

## Luồng xuất hóa đơn điện tử

<Mermaid title="Luồng xuất hóa đơn điện tử">
sequenceDiagram
  participant App as Merchant App
  participant API as SePay eInvoice API
  participant CQT as Cơ quan thuế (CQT)

  App->>API: POST v1/token
  API-->>App: access_token

  App->>API: GET v1/provider-accounts
  API-->>App: Danh sách tài khoản nhà cung cấp

  App->>API: GET v1/provider-accounts/{id}
  API-->>App: Chi tiết tài khoản (mẫu/ký hiệu HĐ)

  App->>API: POST v1/invoices/create
  API-->>App: tracking_code (tạo hóa đơn)

  App->>API: GET v1/invoices/create/check/{tracking_code}
  API-->>App: status, message (tạo hóa đơn)

  App->>API: POST v1/invoices/issue
  API-->>App: tracking_code (phát hành)

  App->>API: GET v1/invoices/issue/check/{tracking_code}
  API-->>App: status, message (phát hành)

  API->>CQT: Gửi hóa đơn lên cơ quan thuế
  CQT-->>API: Xác nhận tiếp nhận

  App->>API: GET v1/invoices/{reference_code}
  API-->>App: Chi tiết hóa đơn + URL file PDF/XML

  App->>API: GET v1/usage
  API-->>App: quota_remaining
</Mermaid>

Luồng xử lý eInvoice gồm các bước tuần tự từ xác thực, kiểm tra tài khoản đến tạo và phát hành hóa đơn:

1. **Lấy access token**: Gọi `POST v1/token` để lấy `access_token` dùng cho tất cả các API tiếp theo.

2. **Danh sách tài khoản nhà cung cấp**: Gọi `GET v1/provider-accounts` để xem các tài khoản hóa đơn điện tử khả dụng và trạng thái từng tài khoản.

3. **Chi tiết tài khoản**: Gọi `GET v1/provider-accounts/{id}` để lấy cấu hình chi tiết: mẫu hóa đơn, ký hiệu, trạng thái hoạt động.

4. **Tạo hóa đơn (Create)**: Gửi dữ liệu hóa đơn qua `POST v1/invoices/create`. API trả về `tracking_code` để theo dõi quá trình xử lý.

5. **Kiểm tra trạng thái tạo hóa đơn**: Gọi `GET v1/invoices/create/check/{tracking_code}` để xác nhận hóa đơn được tạo thành công hay thất bại.

6. **Phát hành hóa đơn (Issue)**: Gửi yêu cầu phát hành qua `POST v1/invoices/issue`. API trả về `tracking_code` riêng cho bước phát hành.

7. **Kiểm tra trạng thái phát hành**: Gọi `GET v1/invoices/issue/check/{tracking_code}` để xác nhận kết quả phát hành lên CQT.

8. **Lấy chi tiết hóa đơn**: Sau khi phát hành thành công, gọi `GET v1/invoices/{reference_code}` để nhận thông tin đầy đủ và URL tải file PDF, XML.

9. **Kiểm tra hạn ngạch**: Gọi `GET v1/usage` để theo dõi số lượt phát hành còn lại.

10. **Danh sách hóa đơn**: Gọi `GET v1/invoices` để lấy danh sách hóa đơn có phân trang phục vụ đối soát và quản lý.

11. **Lịch sử gửi email**: Gọi `GET v1/email-logs` để tra cứu hóa đơn đã được gửi email tới đâu, trạng thái và thời điểm gửi (hỗ trợ cả Production và Sandbox; trên Sandbox dữ liệu được giả lập từ hóa đơn đã phát hành).

12. **Địa điểm kinh doanh người bán**: Gọi `GET v1/business-locations` (kèm các endpoint tạo, cập nhật, xóa) để quản lý danh mục địa điểm, rồi gắn `seller_store_xid` khi tạo hóa đơn để truyền mã số và tên địa điểm lên hóa đơn.

<Callout type="info" title="Xử lý bất đồng bộ">
Các bước 
Create
 và 
Issue
 được xử lý bất đồng bộ. Sau khi gọi mỗi bước, bạn phải gọi API 
`/check`
 tương ứng để xác nhận trạng thái trước khi chuyển sang bước tiếp theo.
</Callout>

***

## Điều chỉnh hóa đơn đã phát hành

Hóa đơn đã gửi lên cơ quan thuế không thể sửa trực tiếp. Khi phát hiện sai sót, bạn cần lập một hóa đơn mới để điều chỉnh hoặc thay thế hóa đơn gốc:

|                  | **[Hóa đơn điều chỉnh](/vi/einvoice-api/v1/xuat-hoa-don-dieu-chinh)** | **[Hóa đơn thay thế](/vi/einvoice-api/v1/xuat-hoa-don-thay-the)** |
| ---------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------- |
| Endpoint         | **`POST v1/invoices/adjust`**                                         | **`POST v1/invoices/replace`**                                    |
| Hóa đơn gốc      | Vẫn còn hiệu lực, được đánh dấu **bị điều chỉnh**                     | Hết hiệu lực, được đánh dấu **bị thay thế**                       |
| Hóa đơn mới      | Chỉ ghi nhận phần chênh lệch hoặc thông tin cần sửa                   | Ghi lại toàn bộ nội dung mới                                      |
| Số lần thực hiện | Có thể thực hiện nhiều lần                                            | Chỉ thực hiện một lần                                             |

Cả hai API đều xử lý bất đồng bộ như API Create. Sau khi gọi API, hệ thống trả về `tracking_code`. Bạn có thể dùng [API theo dõi trạng thái xuất hóa đơn](/vi/einvoice-api/v1/theo-doi-trang-thai-xuat-hoa-don) để kiểm tra kết quả.

Các loại hóa đơn điều chỉnh được hỗ trợ gồm:

* `increase` - Điều chỉnh tăng
* `decrease` - Điều chỉnh giảm
* `combined` - Điều chỉnh tăng và giảm
* `info` - Điều chỉnh thông tin

Muốn lọc riêng các hóa đơn điều chỉnh hoặc thay thế trong danh sách, dùng tham số `document_type` của [API lấy danh sách hóa đơn](/vi/einvoice-api/v1/danh-sach-hoa-don).

***

## Môi trường

| Môi trường     | Base URL                                |
| -------------- | --------------------------------------- |
| **Production** | `https://einvoice-api.sepay.vn`         |
| **Sandbox**    | `https://einvoice-api-sandbox.sepay.vn` |

Xác thực bằng Bearer token trong header: `Authorization: Bearer <ACCESS_TOKEN>`

Nếu bạn mới bắt đầu, hãy thử nghiệm trên Sandbox trước khi chuyển sang Production. Cách lấy token Sandbox, gọi API và dữ liệu cách ly có trong [Sandbox SePay eInvoice API](/vi/einvoice-api/v1/sandbox).

***

## Bước tiếp theo

Để bắt đầu tích hợp eInvoice API, thực hiện theo thứ tự sau:

1. **[Xác thực API hóa đơn điện tử (Bearer Token)](/vi/einvoice-api/v1/tao-token)**: Lấy Bearer token để xác thực các API tiếp theo
2. **[Danh sách nhà cung cấp hóa đơn điện tử](/vi/einvoice-api/v1/danh-sach-tai-khoan)**: Xem các tài khoản nhà cung cấp hóa đơn điện tử khả dụng
3. **[Xuất hóa đơn điện tử](/vi/einvoice-api/v1/xuat-hoa-don-dien-tu)**: Bắt đầu tạo hóa đơn đầu tiên