# API tra cứu lịch sử gửi email hóa đơn

## Tra cứu lịch sử gửi email hóa đơn điện tử qua SePay eInvoice API với phân trang. Cho biết hóa đơn đã gửi tới email nào, trạng thái và thời điểm gửi.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

## API Endpoint

<Endpoint>
  <Method>GET</Method>

  <Path>https://einvoice-api.sepay.vn/v1/email-logs</Path>

  <Description>
    Lịch sử gửi email hóa đơn
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## API Request

<Params>
  <QueryParams>
    <Param name="page" type="integer" required="false">
      Trang hiện tại (mặc định 1)
    </Param>
    <Param name="per_page" type="integer" required="false">
      Số bản ghi mỗi trang (mặc định 20, tối đa 100)
    </Param>
    <Param name="invoice_number" type="string" required="false">
      Lọc theo số hóa đơn (khớp một phần)
    </Param>
    <Param name="status" type="enum: pending, sending, sent, failed, retrying, permanent_failed" required="false">
      Lọc theo trạng thái gửi email. Nếu không truyền sẽ trả về tất cả.
- `pending`: Chờ gửi
- `sending`: Đang gửi
- `sent`: Đã gửi thành công
- `failed`: Gửi thất bại (còn có thể thử lại)
- `retrying`: Đang thử gửi lại
- `permanent_failed`: Thất bại vĩnh viễn (đã hết lượt thử lại)

    </Param>
    <Param name="to_email" type="string" required="false">
      Lọc theo địa chỉ email nhận (khớp một phần)
    </Param>
    <Param name="reference_code" type="string" required="false">
      Lọc theo mã tham chiếu hóa đơn (khớp một phần)
    </Param>
  </QueryParams>

</Params>

<Callout type="info" title="Phân trang">
Sử dụng 
`page`
 để chỉ định trang cần lấy (mặc định là 1).
Sử dụng 
`per_page`
 để giới hạn số bản ghi mỗi trang (mặc định là 20, tối đa 100).
Response trả về 
`has_more: true`
 nếu còn dữ liệu ở trang tiếp theo.
</Callout>

## API Response

<Responses>
  <Response status="200">
    <Description>
      Lịch sử gửi email hóa đơn
    </Description>

    <Example>
      {
        "data": {
          "paging": {
            "per_page": 20,
            "total": 3,
            "has_more": false,
            "current_page": 1,
            "page_count": 1
          },
          "items": [
            {
              "invoice_series": "C26TSE",
              "invoice_number": "3872",
              "full_invoice_number": "C26TSE3872",
              "reference_code": "3a012136-6577-11f1-9f9b-52c7e9b4f41b",
              "to_email": "buyer@example.com",
              "status": "sent",
              "sent_at": "2026-06-11 14:26:40",
              "created_at": "2026-06-11 14:26:35"
            }
          ]
        }
      }
    </Example>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="200">
    <Fields>
      <Field name="data" type="object" required="false">
        <Fields>
          <Field name="paging" type="object" required="false">
            <Fields>
              <Field name="per_page" type="integer" required="false">
                Số bản ghi mỗi trang
              </Field>
              <Field name="total" type="integer" required="false">
                Tổng số bản ghi
              </Field>
              <Field name="has_more" type="boolean" required="false">
                Còn dữ liệu ở trang tiếp theo hay không
              </Field>
              <Field name="current_page" type="integer" required="false">
                Trang hiện tại
              </Field>
              <Field name="page_count" type="integer" required="false">
                Tổng số trang
              </Field>
            </Fields>
          </Field>
          <Field name="items" type="array" required="false">
            <ArrayItems>
              <Fields>
                <Field name="invoice_series" type="string" required="false">
                  Ký hiệu hóa đơn (ví dụ "C26TSE"). `null` khi hóa đơn chưa có ký hiệu.
                </Field>
                <Field name="invoice_number" type="string" required="false">
                  Số hóa đơn
                </Field>
                <Field name="full_invoice_number" type="string" required="false">
                  Số hóa đơn đầy đủ (nối ký hiệu và số), khớp giá trị hiển thị trên web
                </Field>
                <Field name="reference_code" type="string" required="false">
                  Mã tham chiếu hóa đơn
                </Field>
                <Field name="to_email" type="string" required="false">
                  Địa chỉ email nhận hóa đơn
                </Field>
                <Field name="status" type="enum: pending, sending, sent, failed, retrying, permanent_failed" required="false">
                  Trạng thái gửi email. Giá trị enum:
- `pending` - Chờ gửi
- `sending` - Đang gửi
- `sent` - Đã gửi thành công
- `failed` - Gửi thất bại (còn có thể thử lại)
- `retrying` - Đang thử gửi lại
- `permanent_failed` - Thất bại vĩnh viễn (đã hết lượt thử lại)

                </Field>
                <Field name="sent_at" type="string" required="false">
                  Thời điểm gửi thành công (yyyy-MM-dd HH:mm:ss). `null` khi chưa gửi thành công.
                </Field>
                <Field name="created_at" type="string" required="false">
                  Thời điểm tạo bản ghi gửi email (yyyy-MM-dd HH:mm:ss)
                </Field>
              </Fields>
            </ArrayItems>
          </Field>
        </Fields>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

<Callout type="info" title="Về trường status">
`status`
 cho biết trạng thái gửi email: 
`pending`
 (chờ gửi), 
`sending`
 (đang gửi), 
`sent`
 (đã gửi thành công), 
`failed`
 (thất bại, còn có thể thử lại), 
`retrying`
 (đang thử lại), 
`permanent_failed`
 (thất bại vĩnh viễn, đã hết lượt thử lại). Khi email chưa gửi thành công, 
`sent_at`
 là 
`null`
.
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 401, name: "Unauthorized", description: "Thiếu hoặc sai Bearer token. Vui lòng kiểm tra lại token trong header Authorization." }
]}
/>

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
    curl --request GET \
      --url 'https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b' \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "GET",
      CURLOPT_HTTPHEADER => [
        "Authorization: Bearer REPLACE_BEARER_TOKEN"
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
    
    headers = { 'Authorization': "Bearer REPLACE_BEARER_TOKEN" }
    
    conn.request("GET", "/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b", headers=headers)
    
    res = conn.getresponse()
    data = res.read()
    
    print(data.decode("utf-8"))
    ```
  </CodeSample>

  <CodeSample value="node_native" lang="javascript">
    ```javascript
    const http = require("https");
    
    const options = {
      "method": "GET",
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b",
      "headers": {
        "Authorization": "Bearer REPLACE_BEARER_TOKEN"
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
    
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build();
    
    Response response = client.newCall(request).execute();
    ```
  </CodeSample>

  <CodeSample value="ruby_native" lang="ruby">
    ```ruby
    require 'uri'
    require 'net/http'
    require 'openssl'
    
    url = URI("https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Get.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    
    response = http.request(request)
    puts response.read_body
    ```
  </CodeSample>

  <CodeSample value="go_native" lang="go">
    ```go
    package main
    
    import (
    	"fmt"
    	"net/http"
    	"io/ioutil"
    )
    
    func main() {
    
    	url := "https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b"
    
    	req, _ := http.NewRequest("GET", url, nil)
    
    	req.Header.Add("Authorization", "Bearer REPLACE_BEARER_TOKEN")
    
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
        Method = HttpMethod.Get,
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
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
    
    let headers = ["Authorization": "Bearer REPLACE_BEARER_TOKEN"]
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b")! as URL,
                                            cachePolicy: .useProtocolCachePolicy,
                                        timeoutInterval: 10.0)
    request.httpMethod = "GET"
    request.allHTTPHeaderFields = headers
    
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
    
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/email-logs?page=1&per_page=20&invoice_number=3872&status=sent&to_email=buyer%40example.com&reference_code=3a012136-6577-11f1-9f9b-52c7e9b4f41b")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## Bước tiếp theo

Từ mã tham chiếu (`reference_code`) trong kết quả, bạn có thể:

1. **[Chi tiết hóa đơn](/vi/einvoice-api/v1/chi-tiet-hoa-don)** - Xem thông tin chi tiết của hóa đơn tương ứng
2. **[Tải hóa đơn](/vi/einvoice-api/v1/tai-hoa-don)** - Tải file PDF hoặc XML của hóa đơn đã phát hành