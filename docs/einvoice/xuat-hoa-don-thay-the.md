# API xuất hóa đơn thay thế

## Tạo hóa đơn thay thế cho hóa đơn điện tử đã phát hành qua SePay eInvoice API. API trả về tracking_code để theo dõi trạng thái xử lý.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

## Hóa đơn thay thế là gì?

**Hóa đơn thay thế** được lập khi hóa đơn điện tử đã phát hành và được cơ quan thuế (CQT) tiếp nhận có sai sót cần sửa. Khi hóa đơn thay thế được phát hành thành công, hóa đơn gốc được đánh dấu **bị thay thế** và không còn hiệu lực.

Khác với [hóa đơn điều chỉnh](/vi/einvoice-api/v1/xuat-hoa-don-dieu-chinh) chỉ ghi nhận phần chênh lệch, hóa đơn thay thế ghi **đầy đủ nội dung chính xác** của hóa đơn mới, bao gồm dòng hàng hóa, số lượng, đơn giá, thuế suất và thông tin người mua.

Ví dụ: Hóa đơn gốc có sai tên hàng và đơn giá ở hai dòng. Thay vì lập các hóa đơn điều chỉnh, bạn có thể lập một hóa đơn thay thế với toàn bộ nội dung đã được sửa chính xác.

### Điều kiện của hóa đơn gốc

Hóa đơn gốc được xác định bằng `original_invoice_reference_code` (chính là `reference_code` của hóa đơn đã tạo trước đó, cùng `provider_account_id`).

Hóa đơn gốc phải:

* Đã phát hành thành công và được cơ quan thuế tiếp nhận. Hóa đơn nháp hoặc đang xử lý không thể thay thế.
* Là hóa đơn mới hoặc hóa đơn thay thế. **Hóa đơn đã được điều chỉnh không thay thế được.**
* Chưa bị thay thế hoặc xóa bỏ.
* Chưa từng được thay thế. Mỗi hóa đơn gốc chỉ được thay thế một lần.

Một hóa đơn thay thế có thể tiếp tục được thay thế nếu phát hiện sai sót.

Lần sửa đầu tiên chọn hình thức cho tất cả các lần sau. Hóa đơn thay thế thì chỉ thay thế tiếp, không lập được hóa đơn điều chỉnh. Ngược lại, hóa đơn đã được [điều chỉnh](/vi/einvoice-api/v1/xuat-hoa-don-dieu-chinh) chỉ điều chỉnh tiếp, không lập được hóa đơn thay thế.

### Cách hoạt động trong SePay eInvoice API

API `POST v1/invoices/replace` nhận `original_invoice_reference_code` cùng toàn bộ nội dung của hóa đơn mới.

API trả về `tracking_code` ngay sau khi nhận yêu cầu. Việc ký số và gửi hóa đơn lên CQT được xử lý **bất đồng bộ**. Bạn dùng `tracking_code` để theo dõi trạng thái qua [API theo dõi trạng thái xuất hóa đơn](/vi/einvoice-api/v1/theo-doi-trang-thai-xuat-hoa-don) và lấy kết quả cuối cùng.

Mỗi hóa đơn gốc chỉ được thay thế một lần. Nếu hóa đơn thay thế tiếp tục có sai sót, bạn lập hóa đơn thay thế mới cho chính hóa đơn thay thế đó, không quay lại hóa đơn gốc ban đầu.

***

## API Endpoint

<Endpoint>
  <Method>POST</Method>

  <Path>https://einvoice-api.sepay.vn/v1/invoices/replace</Path>

  <Description>
    Xuất hóa đơn thay thế
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
        Mã tham chiếu (reference_code) của hóa đơn gốc đã phát hành cần thay thế
      </Field>
      <Field name="issued_date" type="string" required="false">
        Ngày lập hóa đơn (YYYY-MM-DD HH:mm:ss). Không bắt buộc - mặc định là thời điểm hiện tại. Nếu truyền thì phải là ngày hiện tại.
      </Field>
      <Field name="reference_code" type="string" required="false">
        Mã tham chiếu của hóa đơn thay thế, phải là duy nhất. Nếu không truyền, hệ thống tự sinh UUID.
      </Field>
      <Field name="currency" type="enum: VND, USD, CAD" required="false">
        Đơn vị tiền tệ của hóa đơn thay thế. **Không bắt buộc** - bỏ trống thì dùng lại đơn vị tiền tệ của hóa đơn gốc.

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
      <Field name="buyer" type="object" required="true">
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
        <Description>Danh sách hàng hóa/dịch vụ</Description>
        <ArrayItems>
          <Fields>
            <Field name="line_number" type="integer" required="true">
              Số thứ tự dòng
            </Field>
            <Field name="line_type" type="enum: 1, 2, 3, 4" required="true">
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
              Số lượng
            </Field>
            <Field name="unit_price" type="number" required="false">
              Đơn giá
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
            <Field name="tax_amount" type="integer" required="false">
              Tiền thuế của dòng. **Không bắt buộc.** Nếu truyền thì phải là **số nguyên** (không có phần thập phân).

- **Không truyền:** hệ thống tự tính từ thành tiền trước thuế và thuế suất, làm tròn về số nguyên.
- **Có truyền:** bạn tự tính và **tự làm tròn** về số nguyên; hệ thống dùng **đúng giá trị bạn gửi**, **KHÔNG tính lại / đối chiếu** với `tax_rate`. Gửi giá trị có phần thập phân sẽ bị từ chối với lỗi `400`.

Lưu ý: hóa đơn bán hàng (template_code = 2) - hệ thống tự bỏ qua `tax_rate` và `tax_amount` truyền lên.

            </Field>
            <Field name="discount_tax" type="number" required="false">
              Phần trăm chiết khấu trên sản phẩm (%)
            </Field>
            <Field name="discount_amount" type="number" required="false">
              Số tiền chiết khấu trên sản phẩm
            </Field>
            <Field name="before_discount_and_tax_amount" type="number" required="false">
              Số tiền trước chiết khấu và thuế (dùng cho line_type=3)
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
        "buyer": {
          "type": "company",
          "name": "Công ty ABC",
          "legal_name": "CÔNG TY TNHH ABC",
          "tax_code": "0101234567",
          "address": "123 Đường A, Quận B, Hà Nội",
          "email": "buyer@example.com"
        },
        "items": [
          {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Sản phẩm A",
            "unit": "cái",
            "quantity": 2,
            "unit_price": 4500000,
            "tax_rate": 10
          }
        ],
        "notes": "Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
        "payment_method": "CK",
        "is_draft": false
      }
    </Example>
  </RequestBody>
</Params>

## API Response

<Responses>
  <Response status="202">
    <Description>
      Yêu cầu xuất hóa đơn thay thế đã được tiếp nhận
    </Description>

    <Example>
      {
        "success": true,
        "data": {
          "tracking_code": "3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b",
          "tracking_url": "https://einvoice-api.sepay.vn/v1/invoices/replace/check/3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b",
          "message": "Đã tạo yêu cầu xuất hóa đơn thay thế"
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
  { code: 400, name: "INVALID_IS_DRAFT_TYPE", description: "is_draft phải là boolean thật (true/false), không phải chuỗi." },
  { code: 400, name: "ISSUED_DATE_INVALID", description: "issued_date không phải là ngày hiện tại." },
  { code: 400, name: "INVALID_ISSUED_DATE", description: "issued_date không đúng định dạng Y-m-d H:i:s." },
  { code: 422, name: "INVALID_ORIGIN_INVOICE", description: "Hóa đơn gốc không hợp lệ: chưa phát hành, đã bị thay thế/xóa bỏ, đã được điều chỉnh (loại này chỉ được điều chỉnh tiếp), hoặc đã có hóa đơn thay thế khác được phát hành." },
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

## Lưu ý sử dụng

* `items` dùng chung bộ tham số với **[API xuất hóa đơn điện tử](/vi/einvoice-api/v1/xuat-hoa-don-dien-tu)** - ghi lại **toàn bộ nội dung đúng** của hóa đơn, KHÔNG ghi phần chênh lệch. Tổng tiền hóa đơn không được âm.
* Dòng hàng hóa/dịch vụ và hàng khuyến mại (`line_type` = 1, 2) bắt buộc có `item_code`, `unit` và `quantity` lớn hơn 0.
* Hóa đơn gốc là hóa đơn GTGT: `tax_rate` **bắt buộc trên từng dòng**, trừ dòng ghi chú (`line_type = 4`) - thiếu sẽ bị lỗi 400 với thông báo `"thuế suất là bắt buộc đối với hóa đơn giá trị gia tăng (template_code = 1)"`.
* Hóa đơn gốc là hóa đơn bán hàng (`template_code = 2`): không có thuế - `tax_rate`/`tax_amount` truyền lên bị bỏ qua, tiền thuế luôn bằng 0.
* `tax_amount` trên từng dòng không bắt buộc; nếu truyền thì phải là **số nguyên** (giá trị có phần thập phân bị từ chối với lỗi `400`) và hệ thống dùng **đúng giá trị bạn gửi**, không tính lại từ `tax_rate` - giống **[API xuất hóa đơn điện tử](/vi/einvoice-api/v1/xuat-hoa-don-dien-tu)**.
* `buyer` là bắt buộc - hóa đơn thay thế lập lại toàn bộ nội dung, bao gồm cả thông tin người mua.
* `payment_method` không bắt buộc - bỏ trống thì dùng lại phương thức thanh toán của hóa đơn gốc.
* `currency` không bắt buộc (`VND`/`USD`/`CAD`) - bỏ trống thì dùng lại đơn vị tiền tệ của hóa đơn gốc.
* `issued_date` không bắt buộc và luôn phải là **ngày hiện tại** (mặc định là thời điểm gọi API).
* `is_draft=true` tạo bản nháp hóa đơn thay thế - phát hành sau bằng **[API phát hành hóa đơn](/vi/einvoice-api/v1/phat-hanh-hoa-don-dien-tu)** với `reference_code` của bản nháp. `is_draft=false` xuất và phát hành luôn.
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
      --url https://einvoice-api.sepay.vn/v1/invoices/replace \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"provider_account_id":"0aea3134-da40-11f0-aef4-52c7e9b4f41b","original_invoice_reference_code":"084e179d-d95a-11f0-aef4-52c7e9b4f41b","issued_date":"2026-07-30 00:00:00","reference_code":"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b","currency":"VND","payment_method":"CK","is_draft":false,"buyer":{"type":"personal","name":"Công ty TNHH ABC","legal_name":"CÔNG TY CỔ PHẦN ABC","tax_code":"0123456789","address":"123 Đường ABC, Quận 1, TP.HCM","email":"contact@abc.com","phone":"0901234567","buyer_code":"KH-001","national_id":"001234567890"},"items":[{"line_number":1,"line_type":1,"item_code":"SP001","item_name":"Sản phẩm A","unit":"cái","quantity":10,"unit_price":100000,"tax_rate":10,"tax_amount":900000,"discount_tax":10,"discount_amount":100000,"before_discount_and_tax_amount":4500000}],"notes":"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026","total_amount":4950000}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/invoices/replace",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}",
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
    
    payload = "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/invoices/replace", payload, headers)
    
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
      "path": "/v1/invoices/replace",
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
      issued_date: '2026-07-30 00:00:00',
      reference_code: '3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b',
      currency: 'VND',
      payment_method: 'CK',
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
          item_name: 'Sản phẩm A',
          unit: 'cái',
          quantity: 10,
          unit_price: 100000,
          tax_rate: 10,
          tax_amount: 900000,
          discount_tax: 10,
          discount_amount: 100000,
          before_discount_and_tax_amount: 4500000
        }
      ],
      notes: 'Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026',
      total_amount: 4950000
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}");
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/invoices/replace")
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
    
    url = URI("https://einvoice-api.sepay.vn/v1/invoices/replace")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}"
    
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
    
    	url := "https://einvoice-api.sepay.vn/v1/invoices/replace"
    
    	payload := strings.NewReader("{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
    
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
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/invoices/replace"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
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
      "issued_date": "2026-07-30 00:00:00",
      "reference_code": "3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b",
      "currency": "VND",
      "payment_method": "CK",
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
          "item_name": "Sản phẩm A",
          "unit": "cái",
          "quantity": 10,
          "unit_price": 100000,
          "tax_rate": 10,
          "tax_amount": 900000,
          "discount_tax": 10,
          "discount_amount": 100000,
          "before_discount_and_tax_amount": 4500000
        ]
      ],
      "notes": "Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
      "total_amount": 4950000
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/invoices/replace")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"provider_account_id\":\"0aea3134-da40-11f0-aef4-52c7e9b4f41b\",\"original_invoice_reference_code\":\"084e179d-d95a-11f0-aef4-52c7e9b4f41b\",\"issued_date\":\"2026-07-30 00:00:00\",\"reference_code\":\"3f8b1c2a-7d95-11f1-aef4-52c7e9b4f41b\",\"currency\":\"VND\",\"payment_method\":\"CK\",\"is_draft\":false,\"buyer\":{\"type\":\"personal\",\"name\":\"Công ty TNHH ABC\",\"legal_name\":\"CÔNG TY CỔ PHẦN ABC\",\"tax_code\":\"0123456789\",\"address\":\"123 Đường ABC, Quận 1, TP.HCM\",\"email\":\"contact@abc.com\",\"phone\":\"0901234567\",\"buyer_code\":\"KH-001\",\"national_id\":\"001234567890\"},\"items\":[{\"line_number\":1,\"line_type\":1,\"item_code\":\"SP001\",\"item_name\":\"Sản phẩm A\",\"unit\":\"cái\",\"quantity\":10,\"unit_price\":100000,\"tax_rate\":10,\"tax_amount\":900000,\"discount_tax\":10,\"discount_amount\":100000,\"before_discount_and_tax_amount\":4500000}],\"notes\":\"Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026\",\"total_amount\":4950000}")
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/invoices/replace")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Tham số mẫu

### Hóa đơn thay thế cơ bản

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "buyer": {
        "type": "company",
        "name": "Công ty ABC",
        "legal_name": "CÔNG TY TNHH ABC",
        "tax_code": "0101234567",
        "address": "123 Đường A, Quận B, Hà Nội",
        "email": "buyer@example.com"
    },
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Sản phẩm A",
            "unit": "cái",
            "quantity": 2,
            "unit_price": 4500000,
            "tax_rate": 10
        }
    ],
    "notes": "Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "payment_method": "CK",
    "is_draft": false
}
```
</Response>

### Hóa đơn thay thế có chiết khấu và dòng chiết khấu thương mại

<Response title="json">
```json
{
    "provider_account_id": "{{your-provider-account-id}}",
    "original_invoice_reference_code": "084e179d-d95a-11f0-aef4-52c7e9b4f41b",
    "reference_code": "don-thay-the-001",
    "buyer": {
        "type": "personal",
        "name": "Nguyễn Văn A",
        "address": "Số 88 Đường Ánh Sao, Phường Bình An, Quận 9, TP Hồ Chí Minh",
        "email": "buyeremaildemo@gmail.com"
    },
    "items": [
        {
            "line_number": 1,
            "line_type": 1,
            "item_code": "SP001",
            "item_name": "Sản phẩm A",
            "unit": "cái",
            "quantity": 10,
            "unit_price": 100000,
            "discount_tax": 10,
            "tax_rate": 10
        },
        {
            "line_number": 2,
            "line_type": 3,
            "item_name": "Chiết khấu thương mại",
            "before_discount_and_tax_amount": 50000,
            "tax_rate": 10
        }
    ],
    "notes": "Thay thế cho hóa đơn ký hiệu C26TSE số 589 ngày 26/01/2026",
    "payment_method": "TM/CK",
    "is_draft": true
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
 - hóa đơn thay thế có 
`document_type=4`
, hóa đơn gốc chuyển thành 
`document_type=5`
 (bị thay thế)
</Callout>