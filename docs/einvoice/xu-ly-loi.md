# Xử lý lỗi khi tích hợp eInvoice API

## Cấu trúc response lỗi, bảng mã HTTP và mã lỗi thường gặp, lịch retry tự động, cách chẩn đoán khi hóa đơn không được phát hành và cách kiểm tra trước khi gửi lại.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

Trang này tập hợp toàn bộ thông tin về lỗi của eInvoice API: cấu trúc response lỗi, ý nghĩa từng mã lỗi, cơ chế retry tự động và cách xử lý khi hóa đơn không được phát hành.

## Hai loại lỗi

|                          | Lỗi đồng bộ                     | Lỗi bất đồng bộ                                                         |
| ------------------------ | ------------------------------- | ----------------------------------------------------------------------- |
| Nhận ở đâu               | Ngay trong response khi gọi API | Trong response của API theo dõi trạng thái, khi `status` là `Failed`    |
| Hóa đơn đã được tạo chưa | **Chưa được tạo**               | **Chưa xác định**, cần [kiểm tra trước khi gửi lại](#truoc-khi-gui-lai) |
| Có tự động retry không   | Không, request đã kết thúc      | Có, nếu là lỗi tạm thời                                                 |

## Cấu trúc response lỗi

Mọi endpoint trả lỗi theo cùng một cấu trúc:

```json
{
  "success": false,
  "error": {
    "code": "EINVOICE_DOCUMENT_EXISTED",
    "message": "Hệ thống đã ghi nhận hoá đơn cho mã tham chiếu này."
  }
}
```

Xử lý lỗi dựa trên `error.code`, không khớp theo chuỗi `error.message`. `message` là văn bản tiếng Việt dùng để hiển thị hoặc ghi log, nội dung có thể thay đổi.

## Phân loại lỗi

### Lỗi HTTP

| Status | Ý nghĩa                                               | Cách xử lý                                                                                 |
| ------ | ----------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| 400    | Sai tham số hoặc thiếu trường bắt buộc                | Đọc `error.code` để xác định trường bị lỗi, sửa rồi gọi lại                                |
| 401    | Token thiếu, sai, hoặc đã hết hạn                     | Lấy token mới bằng [POST /v1/token](./tao-token), token có hiệu lực 24 giờ                 |
| 404    | Không tìm thấy đối tượng                              | Kiểm tra `reference_code`, `tracking_code` hoặc ID trên URL                                |
| 405    | Sai HTTP method                                       | Endpoint tạo và phát hành dùng POST, endpoint tra cứu dùng GET                             |
| 409    | Xung đột: đối tượng đã tồn tại hoặc đang được sử dụng | Kiểm tra trước, xem [Trước khi gửi lại](#truoc-khi-gui-lai)                                |
| 422    | Dữ liệu đúng cú pháp nhưng không hợp lệ về nghiệp vụ  | Thường do hóa đơn gốc không hợp lệ, Cơ quan Thuế chưa duyệt tờ khai, hoặc sai ngày hóa đơn |
| 500    | Lỗi hệ thống, hoặc lỗi trạng thái tài khoản           | Xử lý theo `error.code`, nếu không có mã nhận diện được thì liên hệ SePay                  |

<Callout type="warn" title="Lỗi 500 có thể do trạng thái tài khoản">
API này trả 
`500`
 cho cả các lỗi trạng thái tài khoản như 
`QUOTA_HAS_BEEN_USERD_UP`
, 
`REGISTRATION_NOT_COMPLETE`
 và 
`BILLING_IS_UNPAID`
. Những trường hợp này phải xử lý trên my.sepay.vn, gọi lại không làm thay đổi kết quả. Chỉ 
`500`
 không kèm mã lỗi nhận diện được mới là lỗi hệ thống, khi đó liên hệ SePay kèm 
`tracking_code`
 hoặc 
`reference_code`
.
</Callout>

### Mã lỗi thường gặp

| `error.code`                       | Cách xử lý                                                                     |
| ---------------------------------- | ------------------------------------------------------------------------------ |
| `VALIDATION_ERROR`                 | Thiếu trường bắt buộc hoặc sai định dạng. `message` nêu trường bị từ chối      |
| `TOKEN_EXPIRED`                    | Token quá 24 giờ. Lấy token mới, xem [Xác thực Bearer Token](./tao-token)      |
| `EINVOICE_DOCUMENT_EXISTED`        | `reference_code` đã có hóa đơn. Xem [Trước khi gửi lại](#truoc-khi-gui-lai)    |
| `QUOTA_HAS_BEEN_USERD_UP`          | Hết hạn mức. Gia hạn rồi gọi lại, xem [Kiểm tra hạn mức](./kiem-tra-han-ngach) |
| `REGISTRATION_NOT_TAX_APPROVED`    | Cơ quan Thuế chưa duyệt tờ khai. Không xử lý được qua API                      |
| `INVOICE_DATE_BEFORE_TAX_APPROVAL` | Ngày hóa đơn trước ngày duyệt tờ khai. Xem `tax_authority_approved_date`       |
| `REGISTRATION_NOT_COMPLETE`        | Chưa hoàn tất đăng ký gói. Xử lý trên my.sepay.vn                              |
| `BILLING_IS_UNPAID`                | Chưa thanh toán gói. Xử lý trên my.sepay.vn                                    |
| `INVALID_ORIGIN_INVOICE`           | Hóa đơn gốc phải đã phát hành, chưa bị thay thế hoặc xóa bỏ                    |
| `ISSUED_DATE_INVALID`              | `issued_date` phải là ngày hiện tại                                            |
| `INVALID_ISSUED_DATE`              | `issued_date` phải theo định dạng `Y-m-d H:i:s`                                |
| `SELLER_STORE_NOT_FOUND`           | `seller_store_xid` không tồn tại                                               |
| `SELLER_STORE_INACTIVE`            | Địa điểm đã bị hủy kích hoạt, không dùng để xuất hóa đơn mới                   |
| `EINVOICE_ACCOUNT_NOT_FOUND`       | Sai `provider_account_id`                                                      |
| `EINVOICE_DRAFT_NOT_FOUND`         | Không tìm thấy hóa đơn nháp theo `reference_code`                              |

`QUOTA_HAS_BEEN_USERD_UP` viết đúng như API trả về, kể cả lỗi chính tả trong mã.

Mỗi endpoint còn có bảng **Xử lý lỗi** riêng, liệt kê đầy đủ mã lỗi của endpoint đó, ví dụ [Xuất hóa đơn](./xuat-hoa-don-dien-tu#xu-ly-loi) và [Phát hành hóa đơn](./phat-hanh-hoa-don-dien-tu#xu-ly-loi).

## Lịch retry

API tạo và phát hành hóa đơn chạy bất đồng bộ: response trả về `tracking_code`, phần xử lý còn lại chạy nền. Khi phần chạy nền gặp lỗi **tạm thời**, SePay tự thử lại với khoảng cách giữa các lần tăng dần theo Fibonacci:

| Lần         | Chờ     | Tổng    |
| ----------- | ------- | ------- |
| 1 (ban đầu) | ngay    | 0 phút  |
| 2           | 1 phút  | 1 phút  |
| 3           | 1 phút  | 2 phút  |
| 4           | 2 phút  | 4 phút  |
| 5           | 3 phút  | 7 phút  |
| 6           | 5 phút  | 12 phút |
| 7           | 8 phút  | 20 phút |
| 8 (cuối)    | 13 phút | 33 phút |

Tổng cộng 8 lần xử lý (1 lần đầu cộng 7 lần retry), kéo dài khoảng 33 phút nếu tất cả đều thất bại.

Ngoài số lần, còn một giới hạn thứ hai: yêu cầu tạo quá **2 giờ** sẽ không được retry nữa, kể cả khi chưa dùng hết 8 lần. Đây là vùng an toàn dự phòng cho trường hợp một lần xử lý bị treo. Thông thường không chạm tới giới hạn này vì chuỗi 33 phút đã kết thúc trước đó.

Chỉ lỗi hạ tầng được retry: nhà cung cấp hóa đơn trả 5xx, timeout, mất kết nối, hoặc lần xử lý bị ngắt giữa chừng. Lỗi do dữ liệu hoặc do trạng thái tài khoản luôn cho cùng một kết quả ở mọi lần thử, nên hệ thống trả `Failed` ngay.

<Callout type="info" title="Thời điểm là sớm nhất, không phải chính xác">
Cron quét chạy theo chu kỳ 30 giây, nên một lần retry có thể muộn hơn 
`next_retry_at`
 vài chục giây. Không dùng 
`next_retry_at`
 làm mốc timeout cứng.
</Callout>

### Khối `retry` trong response trạng thái

Khi `status` là `Pending` và request **đã gặp lỗi tạm thời ít nhất một lần**, hai API theo dõi trạng thái trả thêm khối `retry`:

```json
{
  "status": "Pending",
  "retry": {
    "retries_count": 2,
    "max_retries": 7,
    "next_retry_at": "2026-09-04 15:12:30"
  }
}
```

| Trường          | Ý nghĩa                                               |
| --------------- | ----------------------------------------------------- |
| `retries_count` | Số lần đã retry                                       |
| `max_retries`   | Số lần retry tối đa, hiện là 7                        |
| `next_retry_at` | Thời điểm dự kiến của lần retry kế tiếp, giờ Việt Nam |

Hai điểm dễ hiểu nhầm:

* **Không có khối `retry`**: request chưa gặp lỗi lần nào và đang chờ xử lý lần đầu, không phải retry đã bị tắt.
* **`next_retry_at` là `null`**: không còn lần retry nào được lên lịch, do đã hết 7 lần hoặc quá 2 giờ, và hệ thống đang chốt kết quả cuối cùng. Không phải sắp retry ngay.

## Chẩn đoán khi hóa đơn không được phát hành

Kiểm tra lần lượt theo thứ tự dưới đây. Nguyên nhân thường nằm ở các bước đầu.

### 1. Kiểm tra response

`success: false` là lỗi đồng bộ, đọc `error.code` theo các bảng ở trên. `success: true` kèm `data.tracking_code` nghĩa là yêu cầu đã được tiếp nhận, chuyển sang bước 2.

### 2. Kiểm tra môi trường

Production là `https://einvoice-api.sepay.vn`, Sandbox là `https://einvoice-api-sandbox.sepay.vn`. Token chỉ dùng được cho môi trường đã cấp nó, gửi token Sandbox tới Production sẽ nhận `401`. Xem [Sandbox](./sandbox).

### 3. Kiểm tra hạn token

Token có hiệu lực 24 giờ. Tiến trình chạy dài dùng lại token đã cache từ hôm trước sẽ nhận `401 TOKEN_EXPIRED` giữa chừng. Lấy token mới khi gặp `401`, không cần lấy lại ở mỗi request.

### 4. Kiểm tra `status`

Gọi [Trạng thái xuất hóa đơn](./theo-doi-trang-thai-xuat-hoa-don) hoặc [Trạng thái phát hành](./theo-doi-trang-thai-phat-hanh-hoa-don) với `tracking_code`.

| `status`  | Ý nghĩa                        | Cách xử lý                                                                     |
| --------- | ------------------------------ | ------------------------------------------------------------------------------ |
| `Pending` | Đang xử lý hoặc đang chờ retry | Tiếp tục poll, xem khối `retry` để biết số lần còn lại                         |
| `Success` | Đã hoàn tất                    | Lấy hóa đơn qua [Chi tiết hóa đơn](./chi-tiet-hoa-don)                         |
| `Failed`  | Đã dừng                        | Đọc `message` để biết nguyên nhân, xem [Trước khi gửi lại](#truoc-khi-gui-lai) |

### 5. Hóa đơn đang ở dạng nháp

Gửi `is_draft: true` khi tạo thì hóa đơn dừng ở dạng nháp, chưa gửi lên Cơ quan Thuế và không tính vào hạn mức. Cần gọi thêm [Phát hành hóa đơn](./phat-hanh-hoa-don-dien-tu) để phát hành.

### 6. Kiểm tra hạn mức

Khi hết hạn mức, mọi yêu cầu phát hành đều trả `Failed` với cùng một lý do. Gọi [API kiểm tra hạn mức](./kiem-tra-han-ngach) để xem số lượt còn lại.

### 7. Kiểm tra phê duyệt của Cơ quan Thuế

Khi tờ khai chưa được duyệt, không hóa đơn nào phát hành được. Xem `tax_authority_approved_date` trong [Danh sách tài khoản nhà cung cấp](./danh-sach-tai-khoan). Ngày hóa đơn cũng phải từ ngày này trở đi.

## Trước khi gửi lại

`Failed` không đồng nghĩa với việc chưa có hóa đơn nào được tạo. Kiểm tra trước khi gửi lại:

1. Gọi [Chi tiết hóa đơn](./chi-tiet-hoa-don) với `reference_code` của bạn.
2. Trả về `404`: chưa có hóa đơn nào được ghi nhận, có thể gửi lại.
3. Trả về dữ liệu hóa đơn: hóa đơn **đã tồn tại**, không gửi lại. Xem `status` trong response để biết hóa đơn cần phát hành tiếp hay đã hoàn tất.

<Callout type="warn" title="Không gửi lại khi hệ thống chưa xác nhận được kết quả">
Khi 
`message`
 của một yêu cầu 
`Failed`
 là 
"Hệ thống chưa xác nhận được kết quả phát hành với nhà cung cấp hóa đơn"
, hệ thống 
không xác định được
 hóa đơn đã lên Cơ quan Thuế hay chưa, không phải là hóa đơn chưa lên.
Gửi lại trong trường hợp này có thể tạo ra 
hóa đơn trùng đã ký và đã nộp thuế
. Thủ tục hủy một hóa đơn đã nộp thuế tốn kém hơn nhiều so với việc chờ SePay kiểm tra. Giữ nguyên 
`tracking_code`
 và liên hệ SePay.
</Callout>

API [Xuất hóa đơn](./xuat-hoa-don-dien-tu) chặn sẵn trường hợp này: gửi lại cùng một `reference_code` mà hệ thống đã ghi nhận hóa đơn sẽ trả `409 EINVOICE_DOCUMENT_EXISTED`. Đây là lớp bảo vệ có sẵn, không dùng thay cho bước kiểm tra ở trên. Mỗi hóa đơn cần một `reference_code` duy nhất do bạn sinh ra, và khi gửi lại phải dùng đúng mã đó.

## Câu hỏi thường gặp

<AccordionGroup>
  <Accordion title="Nhận `Pending` có nên tự gửi lại không?">
    Không. Hệ thống đã retry theo lịch ở trên. Một request nữa chỉ tạo thêm luồng xử lý song song, không rút ngắn thời gian phát hành. Việc cần làm là poll API trạng thái, dùng `retry.next_retry_at` để giãn nhịp gọi, và chỉ dừng khi `status` chuyển sang `Success` hoặc `Failed`.
  </Accordion>

  <Accordion title="Vì sao `status` vẫn `Pending` mà `next_retry_at` lại là `null`?">
    Không còn lần retry nào được lên lịch và hệ thống đang chốt kết quả cuối cùng. Trạng thái sẽ chuyển sang `Success` hoặc `Failed`, không nằm mãi ở `Pending`. Tiếp tục gọi **API theo dõi trạng thái** với nhịp thưa hơn, ví dụ mỗi 1-2 phút. Không gửi lại yêu cầu tạo hoặc phát hành hóa đơn.
  </Accordion>

  <Accordion title="Hết hạn mức hóa đơn thì có được retry không?">
    Không. Hết hạn mức không phải lỗi tạm thời, mọi lần thử đều cho cùng kết quả nên hệ thống trả `Failed` ngay. Gia hạn hạn mức rồi gọi lại. Kiểm tra số lượt còn lại bằng [API kiểm tra hạn mức](./kiem-tra-han-ngach) trước khi xuất hóa đơn để tránh tình huống này.
  </Accordion>

  <Accordion title="Một yêu cầu có thể `Pending` bao lâu?">
    Tối đa khoảng 2 giờ. Trường hợp thông thường hoàn tất trong vài giây, chuỗi retry đầy đủ mất khoảng 33 phút, và mốc 2 giờ là giới hạn cuối cùng, sau đó hệ thống bắt buộc chốt kết quả. Không đặt timeout phía bạn ngắn hơn rồi coi yêu cầu là thất bại.
  </Accordion>

  <Accordion title="Nên xử lý lỗi theo `error.code` hay theo `message`?">
    Theo `error.code`. Đây là chuỗi cố định, dùng để rẽ nhánh xử lý. `message` là văn bản tiếng Việt để hiển thị cho người dùng hoặc ghi log, nội dung có thể thay đổi mà không báo trước. Ngoại lệ duy nhất là response bất đồng bộ khi `status` là `Failed`: response này không có `error.code`, nguyên nhân nằm trong `message`.
  </Accordion>
</AccordionGroup>

## Tiếp theo

* [Trạng thái xuất hóa đơn](./theo-doi-trang-thai-xuat-hoa-don): chi tiết response của API theo dõi
* [Trạng thái phát hành](./theo-doi-trang-thai-phat-hanh-hoa-don): chi tiết response của API theo dõi
* [Chi tiết hóa đơn](./chi-tiet-hoa-don): tra hóa đơn theo `reference_code`
* [Kiểm tra hạn mức](./kiem-tra-han-ngach): số lượt phát hành còn lại
* [Sandbox](./sandbox): môi trường thử nghiệm trước khi chạy thật