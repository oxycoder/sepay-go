# API xuất hóa đơn điều chỉnh

## Tạo hóa đơn điều chỉnh tăng, giảm, tăng/giảm hoặc điều chỉnh thông tin cho hóa đơn điện tử đã phát hành qua SePay eInvoice API.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

## Hóa đơn điều chỉnh là gì?

**Hóa đơn điều chỉnh** được lập khi hóa đơn điện tử đã phát hành và được cơ quan thuế (CQT) tiếp nhận có sai sót cần sửa. Hóa đơn điều chỉnh **không thay thế hoặc hủy bỏ hóa đơn gốc**. Hóa đơn gốc vẫn còn hiệu lực, còn hóa đơn điều chỉnh chỉ ghi nhận **phần chênh lệch** hoặc thông tin cần sửa.

Ví dụ: Hóa đơn gốc có giá trị 10.000.000đ, nhưng giá đúng là 11.000.000đ. Khi đó, bạn lập hóa đơn điều chỉnh tăng với giá trị **1.000.000đ**, là phần chênh lệch cần điều chỉnh, không phải lập lại toàn bộ hóa đơn.

### Điều kiện của hóa đơn gốc

Hóa đơn gốc được xác định bằng `original_invoice_reference_code` (chính là `reference_code` của hóa đơn đã tạo trước đó, cùng `provider_account_id`).

Hóa đơn gốc phải:

* Đã phát hành thành công và được cơ quan thuế tiếp nhận. Hóa đơn nháp hoặc đang xử lý không thể điều chỉnh.
* Là hóa đơn mới hoặc hóa đơn đã được điều chỉnh. **Hóa đơn thay thế không điều chỉnh được.**
* Chưa bị thay thế hoặc xóa bỏ.

Một hóa đơn đã được điều chỉnh vẫn có thể tiếp tục điều chỉnh nhiều lần.

Lần sửa đầu tiên chọn hình thức cho tất cả các lần sau. Hóa đơn đã được điều chỉnh thì chỉ điều chỉnh tiếp, không lập được hóa đơn thay thế. Ngược lại, hóa đơn thay thế chỉ được [thay thế](/vi/einvoice-api/v1/xuat-hoa-don-thay-the) tiếp, không lập được hóa đơn điều chỉnh.

### Các loại hóa đơn điều chỉnh được hỗ trợ

Tham số `adjustment_type` xác định loại hóa đơn điều chỉnh:

* `increase` - Điều chỉnh tăng
* `decrease` - Điều chỉnh giảm
* `combined` - Điều chỉnh tăng/giảm trên cùng một hóa đơn
* `info` - Điều chỉnh thông tin, không thay đổi giá trị tiền

Chi tiết cách nhập dòng hàng cho từng loại được trình bày tại mục **Các loại điều chỉnh (adjustment\_type)** bên dưới.

### Cách hoạt động trong SePay eInvoice API

API `POST v1/invoices/adjust` nhận `original_invoice_reference_code`, loại điều chỉnh và các dòng cần điều chỉnh.

API trả về `tracking_code` ngay sau khi nhận yêu cầu. Việc ký số và gửi hóa đơn lên CQT được xử lý **bất đồng bộ**. Bạn dùng `tracking_code` để theo dõi trạng thái qua [API theo dõi trạng thái xuất hóa đơn](/vi/einvoice-api/v1/theo-doi-trang-thai-xuat-hoa-don) và lấy kết quả cuối cùng.

Một hóa đơn có thể được điều chỉnh nhiều lần.

***

## API Endpoint

<Endpoint>
  <Method>POST</Method>

  <Path>https://einvoice-api.sepay.vn/v1/invoices/adjust</Path>

  <Description>
    Xuất hóa đơn điều chỉnh
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## API Request

<Params>
  <RequestBody>
    <Fields>
      <Field name="provider_account_id" type="string" required="true">
        ID tài khoản nhà cung cấp (UUID) - phải là tài khoản đã phát hành hóa đơn gốc
      </Field>
      <Field name="original_invoice_reference_code" type="string" required="true">
        Mã tham chiếu (reference_code) của hóa đơn gốc đã phát hành cần điều chỉnh
      </Field>
      <Field name="adjustment_type" type="enum: increase, decrease, combined, info" required="true">
        Loại điều chỉnh:
- `increase`: Điều chỉnh tăng
- `decrease`: Điều chỉnh giảm
- `combined`: Điều chỉnh tăng/giảm
- `info`: Điều chỉnh thông tin

      </Field>
      <Field name="issued_date" type="string" required="false">
        Ngày lập hóa đơn (YYYY-MM-DD HH:mm:ss). Không bắt buộc - mặc định là thời điểm hiện tại. Nếu truyền thì phải là ngày hiện tại.
      </Field>
      <Field name="reference_code" type="string" required="false">
        Mã tham chiếu của hóa đơn điều chỉnh, phải là duy nhất. Nếu không truyền, hệ thống tự sinh UUID.
      </Field>
      <Field name="payment_method" type="enum: TM, CK, TM/CK, KHAC" required="false">
        Phương thức thanh toán. Không bắt buộc - nếu bỏ trống, hệ thống dùng lại phương thức thanh toán của hóa đơn gốc (TM/CK khi hóa đơn gốc cũng không có).
- TM: Tiền mặt
- CK: Chuyển khoản
- TM/CK: Tiền mặt và chuyển khoản
- KHAC: Khác

      </Field>
      <Field name="is_draft" type="boolean" required="true">
        Phải là boolean thật (true/false):
- `true`: Xuất nháp (phát hành sau qua API phát hành hóa đơn)
- `false`: Xuất và phát hành luôn

      </Field>
      <Field name="buyer" type="object" required="false">
        <Fields>
          <Field name="type" type="enum: personal, company" required="false">
            Loại người mua (personal, company)
          </Field>
          <Field name="name" type="string" required="false">
            Tên người/đơn vị mua
          </Field>
          <Field name="legal_name" type="string" required="false">
            Tên pháp lý (dùng khi buyer.type là company)
          </Field>
          <Field name="tax_code" type="string" required="false">
            Mã số thuế
          </Field>
          <Field name="address" type="string" required="false">
            Địa chỉ
          </Field>
          <Field name="email" type="string (email)" required="false">
            Email nhận hóa đơn
          </Field>
          <Field name="phone" type="string" required="false">
            Số điện thoại
          </Field>
          <Field name="buyer_code" type="string" required="false">
            Mã khách hàng (mã người mua hàng)
          </Field>
          <Field name="national_id" type="string" required="false">
            Căn cước công dân / Số CCCD / Số định danh cá nhân
          </Field>
        </Fields>
      </Field>
      <Field name="items" type="array" required="true">
        <Description>Danh sách dòng điều chỉnh.

`increase`/`decrease`/`combined`: cần ít nhất 1 mặt hàng điều chỉnh.

`info`: cần ít nhất 1 dòng nội dung.
</Description>
        <ArrayItems>
          <Fields>
            <Field name="line_number" type="integer" required="false">
              Số thứ tự dòng
            </Field>
            <Field name="line_type" type="enum: 1, 2, 3, 4" required="false">
              Loại dòng hàng:
- 1: Hàng hóa/dịch vụ bình thường
- 2: Hàng khuyến mại
- 3: Chiết khấu thương mại
- 4: Ghi chú

            </Field>
            <Field name="item_code" type="string" required="false">
              Mã hàng hóa/dịch vụ
            </Field>
            <Field name="item_name" type="string" required="true">
              Tên hàng hóa/dịch vụ
            </Field>
            <Field name="unit" type="string" required="false">
              Đơn vị tính
            </Field>
            <Field name="quantity" type="number" required="false">
              Số lượng - chấp nhận số âm với điều chỉnh tăng/giảm tiền
            </Field>
            <Field name="unit_price" type="number" required="false">
              Đơn giá - chấp nhận số âm với điều chỉnh tăng/giảm tiền
            </Field>
            <Field name="tax_rate" type="enum: -2, -1, 0, 5, 8, 10" required="false">
              Thuế suất (%):
- -2: Không chịu thuế
- -1: Không kê khai, tính nộp thuế GTGT
- 0: 0%
- 5: 5%
- 8: 8%
- 10: 10%

            </Field>
            <Field name="discount_tax" type="number" required="false">
              Phần trăm chiết khấu trên sản phẩm (%), không âm
            </Field>
            <Field name="discount_amount" type="number" required="false">
              Số tiền chiết khấu trên sản phẩm
            </Field>
            <Field name="before_discount_and_tax_amount" type="number" required="false">
              Số tiền trước chiết khấu và thuế (dùng cho line_type=3)
            </Field>
            <Field name="sub_amount" type="number" required="false">
              Thành tiền trước thuế (sau chiết khấu)
            </Field>
            <Field name="tax_amount" type="integer" required="false">
              Tiền thuế của dòng. **Không bắt buộc.** Nếu truyền thì phải là **số nguyên** (không có phần thập phân).

- **Không truyền:** hệ thống tự tính từ thành tiền trước thuế và thuế suất, làm tròn về số nguyên.
- **Có truyền:** bạn tự tính và **tự làm tròn** về số nguyên; hệ thống dùng **đúng giá trị bạn gửi**, **KHÔNG tính lại / đối chiếu** với `tax_rate`. Gửi giá trị có phần thập phân sẽ bị từ chối với lỗi `400`.

Lưu ý: hóa đơn gốc là hóa đơn bán hàng (template_code = 2) - hệ thống tự bỏ qua `tax_rate` và `tax_amount` truyền lên, tiền thuế luôn bằng 0.

            </Field>
            <Field name="amount" type="number" required="false">
              Thành tiền sau thuế của dòng
            </Field>
          </Fields>
        </ArrayItems>
      </Field>
      <Field name="notes" type="string" required="false">
        Ghi chú hóa đơn, tối đa 255 ký tự
      </Field>
      <Field name="total_amount" type="integer" required="false">
        Tổng tiền thanh toán cuối cùng của hóa đơn (đã gồm thuế). **Không bắt buộc.** Nếu truyền thì phải là **số nguyên** (không có phần thập phân).

- **Không truyền:** hệ thống tự tính từ các dòng hàng và làm tròn về số nguyên.
- **Có truyền:** hệ thống dùng **đúng giá trị bạn gửi**, **KHÔNG kiểm tra / đối chiếu** với các dòng hàng.

**Cảnh báo (khi có truyền):** Bạn tự chịu trách nhiệm về giá trị này. Nếu không khớp tiền hàng + thuế, hóa đơn sẽ hiển thị sai và cơ quan thuế có thể từ chối.

      </Field>
    </Fields>

    <Example>
      {
        "provider_account_id": "0aea3134-da40-11f0-aef4-52c7e9b4f41b",
        "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
        "adjustment_type": "increase",
        "items": [
          {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Điều chỉnh tăng mặt hàng Sản phẩm A",
            "unit": "cái",
            "quantity": 1,
            "unit_price": 500000,
            "tax_rate": 10
          }
        ],
        "notes": "Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
        "is_draft": false
      }
    </Example>
  </RequestBody>
</Params>

## API Response

<Responses>
  <Response status="202">
    <Description>
      Yêu cầu xuất hóa đơn điều chỉnh đã được tiếp nhận
    </Description>

    <Example>
      {
        "success": true,
        "data": {
          "tracking_code": "5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b",
          "tracking_url": "https://einvoice-api.sepay.vn/v1/invoices/adjust/check/5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b",
          "message": "Đã tạo yêu cầu xuất hóa đơn điều chỉnh"
        }
      }
    </Example>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="202">
    <Fields>
      <Field name="success" type="boolean" required="false">
      </Field>
      <Field name="data" type="object" required="false">
        <Fields>
          <Field name="tracking_code" type="string" required="true">
            Mã tracking để theo dõi trạng thái (đồng thời là reference_code của hóa đơn mới)
          </Field>
          <Field name="tracking_url" type="string (uri)" required="true">
            URL để check trạng thái xử lý
          </Field>
          <Field name="message" type="string" required="false">
            Thông điệp phản hồi
          </Field>
        </Fields>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Bad Request", description: "Thiếu hoặc sai định dạng các trường bắt buộc (VALIDATION_ERROR)." },
  { code: 401, name: "Unauthorized", description: "Thiếu hoặc sai Bearer token, hoặc token đã hết hạn (TOKEN_EXPIRED)." },
  { code: 400, name: "INVALID_ADJUSTMENT_TYPE", description: "adjustment_type là bắt buộc và phải là một trong: increase, decrease, combined, info." },
  { code: 400, name: "INVALID_IS_DRAFT_TYPE", description: "is_draft phải là boolean thật (true/false), không phải chuỗi." },
  { code: 400, name: "VALIDATION_ERROR", description: "Dữ liệu dòng hàng không hợp lệ: giá trị âm ngoài loại điều chỉnh tiền, hóa đơn điều chỉnh thông tin có tổng khác 0, không có dòng chênh lệch nào khác 0..." },
  { code: 400, name: "ISSUED_DATE_INVALID", description: "issued_date không phải là ngày hiện tại." },
  { code: 400, name: "INVALID_ISSUED_DATE", description: "issued_date không đúng định dạng Y-m-d H:i:s." },
  { code: 422, name: "INVALID_ORIGIN_INVOICE", description: "Hóa đơn gốc không hợp lệ: chưa phát hành, đã bị thay thế/xóa bỏ, hoặc là hóa đơn thay thế (loại này chỉ được thay thế tiếp)." },
  { code: 422, name: "REGISTRATION_NOT_TAX_APPROVED", description: "Doanh nghiệp chưa đăng ký HĐĐT hoặc CQT chưa duyệt tờ khai (chỉ check ở production)." },
  { code: 422, name: "INVOICE_DATE_BEFORE_TAX_APPROVAL", description: "Ngày hóa đơn trước ngày Cơ quan Thuế duyệt tờ khai (chỉ check ở production)." },
  { code: 409, name: "EINVOICE_DOCUMENT_EXISTED", description: "Hệ thống đã ghi nhận hoá đơn cho mã tham chiếu này. reference_code phải là duy nhất." },
  { code: 500, name: "EINVOICE_ACCOUNT_NOT_FOUND", description: "Không tìm thấy tài khoản xuất hóa đơn điện tử tương ứng provider_account_id." },
  { code: 500, name: "QUOTA_HAS_BEEN_USERD_UP", description: "Đã dùng hết hạn ngạch hóa đơn của gói dịch vụ. Cần gia hạn hạn ngạch rồi gọi lại - đừng retry tự động vì lỗi này không tự hết." },
  { code: 500, name: "REGISTRATION_NOT_COMPLETE", description: "Công ty chưa hoàn tất đăng ký gói hóa đơn điện tử." },
  { code: 500, name: "BILLING_IS_UNPAID", description: "Công ty chưa thanh toán cho gói hóa đơn điện tử." },
  { code: 500, name: "Internal Server Error", description: "Lỗi hệ thống khi xử lý yêu cầu." }
]}
/>

## Các loại điều chỉnh (adjustment\_type)

| adjustment\_type | Loại hóa đơn                 | Cách nhập dòng hàng                                                                                |
| ---------------- | ---------------------------- | -------------------------------------------------------------------------------------------------- |
| `increase`       | Hóa đơn điều chỉnh tăng      | Nhập phần chênh lệch TĂNG, giá trị dương                                                           |
| `decrease`       | Hóa đơn điều chỉnh giảm      | Nhập phần chênh lệch GIẢM - có thể truyền số âm hoặc số dương                                      |
| `combined`       | Hóa đơn điều chỉnh tăng/giảm | Mỗi dòng mang dấu riêng: dương = tăng, âm = giảm                                                   |
| `info`           | Hóa đơn điều chỉnh thông tin | Không mang giá trị tiền, tổng phải bằng 0. Nội dung điều chỉnh ghi vào dòng ghi chú (line\_type=4) |

## Thông tin người mua (buyer)

* Với `increase`/`decrease`/`combined`: `buyer` **không bắt buộc**. Nếu bỏ trống, hệ thống dùng lại **toàn bộ** thông tin người mua của hóa đơn gốc.
* Với `info`: `buyer` **bắt buộc** - đây chính là thông tin đúng cần điều chỉnh (ví dụ địa chỉ mới, mã số thuế đúng).

## Lưu ý sử dụng

* `payment_method` không bắt buộc - bỏ trống thì dùng lại phương thức thanh toán của hóa đơn gốc.
* `issued_date` không bắt buộc và luôn phải là **ngày hiện tại** (mặc định là thời điểm gọi API).
* Với `increase`/`decrease`/`combined`: cần ít nhất một dòng có giá trị khác 0. Các dòng hàng hóa toàn số 0 tự động bị bỏ qua.
* Với `decrease`: giá trị dòng hàng có thể truyền **số âm hoặc số dương** - hệ thống giữ nguyên đúng giá trị bạn gửi, loại điều chỉnh giảm đã được xác định bởi `adjustment_type=decrease`.
* Với `info`: cần ít nhất một dòng nội dung mô tả phần thông tin được điều chỉnh (thường là dòng ghi chú `line_type=4`), và tổng tiền phải bằng 0.
* `discount_tax` (chiết khấu %) không bao giờ âm - điều chỉnh phần chiết khấu thì truyền thẳng qua `discount_amount`, theo cùng quy ước dấu với các dòng hàng.
* Hóa đơn gốc là hóa đơn bán hàng (`template_code = 2`): không có thuế - `tax_rate`/`tax_amount` truyền lên bị bỏ qua, tiền thuế luôn bằng 0.
* `tax_amount` trên từng dòng: nếu truyền thì phải là **số nguyên** (giá trị có phần thập phân bị từ chối với lỗi `400`) và hệ thống dùng **đúng giá trị bạn gửi**, không tính lại từ `tax_rate`.
* `is_draft=true` tạo bản nháp - phát hành sau bằng **[API phát hành hóa đơn](/vi/einvoice-api/v1/phat-hanh-hoa-don-dien-tu)** với `reference_code` của bản nháp. `is_draft=false` xuất và phát hành luôn.
* Sau khi gửi yêu cầu thành công, dùng `tracking_url` để gọi **[API theo dõi trạng thái xuất hóa đơn](/vi/einvoice-api/v1/theo-doi-trang-thai-xuat-hoa-don)** - luôn poll tới khi có kết quả cuối (thành công hoặc thất bại).

## Code mẫu

<CodeSamples>
  <CodeSamplesList>
    <CodeSamplesTrigger value="shell_curl">
      cURL
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="php_curl">
      PHP
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="python_python3">
      Python
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="node_native">
      NodeJS
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="java_okhttp">
      Java
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="ruby_native">
      Ruby
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="go_native">
      Go
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="csharp_httpclient">
      .NET
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="swift_nsurlsession">
      Swift
    </CodeSamplesTrigger>

    <CodeSamplesTrigger value="kotlin_okhttp">
      Kotlin
    </CodeSamplesTrigger>

  </CodeSamplesList>

  <CodeSample value="shell_curl" lang="bash">
    ```bash
    curl --request POST \
      --url https://einvoice-api.sepay.vn/v1/invoices/adjust \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"provider_account_id":"0aea3134-da40-11f0-aef4-52c7e9b4f41b","original_invoice_reference_code":"084e179d-d95a-11f0-aef4-52c7e9b4f41b","adjustment_type":"increase","issued_date":"2026-07-30 00:00:00","reference_code":"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b","payment_method":"TM/CK","is_draft":false,"buyer":{"type":"personal","name":"Công ty TNHH ABC","legal_name":"CÔNG TY CỔ PHẦN ABC","tax_code":"0123456789","address":"123 Đường ABC, Quận 1, TP.HCM","email":"contact@abc.com","phone":"0901234567","buyer_code":"KH-001","national_id":"001234567890"},"items":[{"line_number":1,"line_type":1,"item_code":"SP001","item_name":"Điều chỉnh tăng mặt hàng Sản phẩm A","unit":"cái","quantity":1,"unit_price":500000,"tax_rate":10,"discount_tax":10,"discount_amount":100000,"before_discount_and_tax_amount":500000,"sub_amount":400000,"tax_amount":40000,"amount":440000}],"notes":"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026","total_amount":4950000}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/invoices/adjust",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}",
      CURLOPT_HTTPHEADER => [
        "Authorization: Bearer REPLACE_BEARER_TOKEN",
        "content-type: application/json"
      ],
    ]);
    
    $response = curl_exec($curl);
    $err = curl_error($curl);
    
    curl_close($curl);
    
    if ($err) {
      echo "cURL Error #:" . $err;
    } else {
      echo $response;
    }
    ```
  </CodeSample>

  <CodeSample value="python_python3" lang="python">
    ```python
    import http.client
    
    conn = http.client.HTTPSConnection("einvoice-api.sepay.vn")
    
    payload = "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/invoices/adjust", payload, headers)
    
    res = conn.getresponse()
    data = res.read()
    
    print(data.decode("utf-8"))
    ```
  </CodeSample>

  <CodeSample value="node_native" lang="javascript">
    ```javascript
    const http = require("https");
    
    const options = {
      "method": "POST",
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/invoices/adjust",
      "headers": {
        "Authorization": "Bearer REPLACE_BEARER_TOKEN",
        "content-type": "application/json"
      }
    };
    
    const req = http.request(options, function (res) {
      const chunks = [];
    
      res.on("data", function (chunk) {
        chunks.push(chunk);
      });
    
      res.on("end", function () {
        const body = Buffer.concat(chunks);
        console.log(body.toString());
      });
    });
    
    req.write(JSON.stringify({
      provider_account_id: '0aea3134-da40-11f0-aef4-52c7e9b4f41b',
      original_invoice_reference_code: '084e179d-d95a-11f0-aef4-52c7e9b4f41b',
      adjustment_type: 'increase',
      issued_date: '2026-07-30 00:00:00',
      reference_code: '5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b',
      payment_method: 'TM/CK',
      is_draft: false,
      buyer: {
        type: 'personal',
        name: 'Công ty TNHH ABC',
        legal_name: 'CÔNG TY CỔ PHẦN ABC',
        tax_code: '0123456789',
        address: '123 Đường ABC, Quận 1, TP.HCM',
        email: 'contact@abc.com',
        phone: '0901234567',
        buyer_code: 'KH-001',
        national_id: '001234567890'
      },
      items: [
        {
          line_number: 1,
          line_type: 1,
          item_code: 'SP001',
          item_name: 'Điều chỉnh tăng mặt hàng Sản phẩm A',
          unit: 'cái',
          quantity: 1,
          unit_price: 500000,
          tax_rate: 10,
          discount_tax: 10,
          discount_amount: 100000,
          before_discount_and_tax_amount: 500000,
          sub_amount: 400000,
          tax_amount: 40000,
          amount: 440000
        }
      ],
      notes: 'Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026',
      total_amount: 4950000
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}");
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/invoices/adjust")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build();
    
    Response response = client.newCall(request).execute();
    ```
  </CodeSample>

  <CodeSample value="ruby_native" lang="ruby">
    ```ruby
    require 'uri'
    require 'net/http'
    require 'openssl'
    
    url = URI("https://einvoice-api.sepay.vn/v1/invoices/adjust")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}"
    
    response = http.request(request)
    puts response.read_body
    ```
  </CodeSample>

  <CodeSample value="go_native" lang="go">
    ```go
    package main
    
    import (
    	"fmt"
    	"strings"
    	"net/http"
    	"io/ioutil"
    )
    
    func main() {
    
    	url := "https://einvoice-api.sepay.vn/v1/invoices/adjust"
    
    	payload := strings.NewReader("{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
    
    	req, _ := http.NewRequest("POST", url, payload)
    
    	req.Header.Add("Authorization", "Bearer REPLACE_BEARER_TOKEN")
    	req.Header.Add("content-type", "application/json")
    
    	res, _ := http.DefaultClient.Do(req)
    
    	defer res.Body.Close()
    	body, _ := ioutil.ReadAll(res.Body)
    
    	fmt.Println(res)
    	fmt.Println(string(body))
    
    }
    ```
  </CodeSample>

  <CodeSample value="csharp_httpclient" lang="csharp">
    ```csharp
    var client = new HttpClient();
    var request = new HttpRequestMessage
    {
        Method = HttpMethod.Post,
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/invoices/adjust"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
        {
            Headers =
            {
                ContentType = new MediaTypeHeaderValue("application/json")
            }
        }
    };
    using (var response = await client.SendAsync(request))
    {
        response.EnsureSuccessStatusCode();
        var body = await response.Content.ReadAsStringAsync();
        Console.WriteLine(body);
    }
    ```
  </CodeSample>

  <CodeSample value="swift_nsurlsession" lang="swift">
    ```swift
    import Foundation
    
    let headers = [
      "Authorization": "Bearer REPLACE_BEARER_TOKEN",
      "content-type": "application/json"
    ]
    let parameters = [
      "provider_account_id": "0aea3134-da40-11f0-aef4-52c7e9b4f41b",
      "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
      "adjustment_type": "increase",
      "issued_date": "2026-07-30 00:00:00",
      "reference_code": "5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b",
      "payment_method": "TM/CK",
      "is_draft": false,
      "buyer": [
        "type": "personal",
        "name": "Công ty TNHH ABC",
        "legal_name": "CÔNG TY CỔ PHẦN ABC",
        "tax_code": "0123456789",
        "address": "123 Đường ABC, Quận 1, TP.HCM",
        "email": "contact@abc.com",
        "phone": "0901234567",
        "buyer_code": "KH-001",
        "national_id": "001234567890"
      ],
      "items": [
        [
          "line_number": 1,
          "line_type": 1,
          "item_code": "SP001",
          "item_name": "Điều chỉnh tăng mặt hàng Sản phẩm A",
          "unit": "cái",
          "quantity": 1,
          "unit_price": 500000,
          "tax_rate": 10,
          "discount_tax": 10,
          "discount_amount": 100000,
          "before_discount_and_tax_amount": 500000,
          "sub_amount": 400000,
          "tax_amount": 40000,
          "amount": 440000
        ]
      ],
      "notes": "Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
      "total_amount": 4950000
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/invoices/adjust")! as URL,
                                            cachePolicy: .useProtocolCachePolicy,
                                        timeoutInterval: 10.0)
    request.httpMethod = "POST"
    request.allHTTPHeaderFields = headers
    request.httpBody = postData as Data
    
    let session = URLSession.shared
    let dataTask = session.dataTask(with: request as URLRequest, completionHandler: { (data, response, error) -> Void in
      if (error != nil) {
        print(error)
      } else {
        let httpResponse = response as? HTTPURLResponse
        print(httpResponse)
      }
    })
    
    dataTask.resume()
    ```
  </CodeSample>

  <CodeSample value="kotlin_okhttp" lang="kotlin">
    ```kotlin
    val client = OkHttpClient()
    
    val mediaType = MediaType.parse("application/json")
    val body = RequestBody.create(mediaType, "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"adjustment_type\":\"increase\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"5c2d9e4b-7d96-11f1-aef4-52c7e9b4f41b\",\"payment_method\":\"TM/CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Điều chỉnh tăng mặt hàng Sản phẩm A\",\"unit\":\"cái\",\"quantity\":1,\"unit_price\":500000,\"tax_rate\":10,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":500000,\"sub_amount\":400000,\"tax_amount\":40000,\"amount\":440000}],\"notes\":\"Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/invoices/adjust")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Tham số mẫu cho các loại điều chỉnh

Các JSON dưới đây chỉ là ví dụ minh họa cho request - số liệu và nội dung là giả định.

### Hóa đơn điều chỉnh tăng (increase)

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "adjustment_type": "increase",
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Điều chỉnh tăng mặt hàng Sản phẩm A",
            "unit": "cái",
            "quantity": 1,
            "unit_price": 500000,
            "tax_rate": 10
        }
    ],
    "notes": "Điều chỉnh tăng đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "is_draft": false
}
```
</Response>

### Hóa đơn điều chỉnh giảm (decrease)

Ví dụ dưới ghi số âm - với `decrease` bạn cũng có thể truyền số dương, hệ thống giữ nguyên giá trị bạn gửi.

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "adjustment_type": "decrease",
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Điều chỉnh giảm mặt hàng Sản phẩm A",
            "unit": "cái",
            "quantity": 1,
            "unit_price": -500000,
            "tax_rate": 10
        }
    ],
    "notes": "Điều chỉnh giảm đơn giá cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "is_draft": false
}
```
</Response>

### Hóa đơn điều chỉnh tăng/giảm (combined - dấu theo từng dòng)

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "adjustment_type": "combined",
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Điều chỉnh tăng mặt hàng Sản phẩm A",
            "unit": "cái",
            "quantity": 1,
            "unit_price": 300000,
            "tax_rate": 10
        },
        {
            "line_number": 2,
            "line_type": 1,
            "item_code": "SP002",
            "item_name": "Điều chỉnh giảm mặt hàng Sản phẩm B",
            "unit": "cái",
            "quantity": 1,
            "unit_price": -200000,
            "tax_rate": 10
        }
    ],
    "notes": "Điều chỉnh tăng/giảm cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "is_draft": false
}
```
</Response>

### Hóa đơn điều chỉnh thông tin (info - tổng bằng 0, buyer bắt buộc)

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "adjustment_type": "info",
    "buyer": {
        "type": "company",
        "name": "Công ty ABC",
        "legal_name": "CÔNG TY TNHH ABC",
        "tax_code": "0101234567",
        "address": "456 Đường B, Quận C, Hà Nội",
        "email": "buyer@example.com"
    },
    "items": [
        {
            "line_number": 1,
            "line_type": 4,
            "item_name": "Điều chỉnh địa chỉ người mua từ 123 Đường A thành 456 Đường B"
        }
    ],
    "notes": "Điều chỉnh thông tin cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "is_draft": false
}
```
</Response>

### Điều chỉnh chỉ số tiền dòng bằng thành tiền sau thuế

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "adjustment_type": "decrease",
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Điều chỉnh giảm mặt hàng Sản phẩm A",
            "unit": "cái",
            "tax_rate": 10,
            "amount": -550000
        }
    ],
    "notes": "Điều chỉnh giảm thành tiền cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "is_draft": false
}
```
</Response>

***

## Bước tiếp theo

Sau khi gửi yêu cầu thành công và nhận được `tracking_code`:

1. **[Theo dõi trạng thái xuất hóa đơn](/vi/einvoice-api/v1/theo-doi-trang-thai-xuat-hoa-don)** - Sử dụng `tracking_code` để kiểm tra kết quả xử lý (bắt buộc)

<Callout type="info" title="Sau khi xác nhận trạng thái thành công">
Nếu xuất 
nháp
 (
`is_draft=true`
): Tiếp tục 
Phát hành hóa đơn
 để phát hành chính thức
Nếu xuất 
chính thức
 (
`is_draft=false`
): Có thể 
Tải hóa đơn
 hoặc 
Xem chi tiết hóa đơn
 - hóa đơn điều chỉnh có 
`document_type`
 6/7/8/9 tùy loại, hóa đơn gốc chuyển thành 
`document_type=10`
 (bị điều chỉnh)
</Callout>