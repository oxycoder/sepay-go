# API thống kê giao dịch Merchant

## GET /merchants/counter — số liệu giao dịch theo ngày của merchant Bank Hub, hỗ trợ lọc theo khoảng thời gian và xem tổng hợp hoạt động giao dịch.

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

## API Endpoint

<Endpoint>
  <Method>GET</Method>

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter</Path>

  <Description>
    Bộ đếm giao dịch Merchant
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Nếu không truyền 
`start_date`
 và 
`end_date`
, mặc định trả về thống kê của ngày hiện tại
`start_date`
 phải nhỏ hơn hoặc bằng 
`end_date`
Trường 
`transaction`
 trong mỗi ngày chỉ đếm giao dịch tiền vào (không bao gồm tiền ra)
</Callout>

## Query Parameters

<Params>
  <QueryParams>
    <Param name="start_date" type="string" required="false">
      Ngày bắt đầu (YYYY-MM-DD). Mặc định là ngày hiện tại.
    </Param>
    <Param name="end_date" type="string" required="false">
      Ngày kết thúc (YYYY-MM-DD). Mặc định là ngày hiện tại.
    </Param>
  </QueryParams>

</Params>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Bad Request", description: "Dữ liệu đầu vào không hợp lệ — định dạng ngày sai, ngày không tồn tại trong lịch, hoặc start_date lớn hơn end_date" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 500, name: "Internal Server Error", description: "Lỗi máy chủ nội bộ" }
]}
/>

<Callout type="warn" title="Validation ngày">
Các lỗi 400 phổ biến khi truyền tham số ngày:
`start_date`
 hoặc 
`end_date`
 sai định dạng (phải là 
`YYYY-MM-DD`
)
Ngày không tồn tại trong lịch (ví dụ: 
`2026-02-30`
)
`start_date`
 lớn hơn 
`end_date`
</Callout>

## API Response

<Responses>
  <Response status="200">
    <Description>
      Thành công
    </Description>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="200">
    <Fields>
      <Field name="data" type="object" required="false">
        <Description>Dữ liệu thống kê</Description>
        <Fields>
          <Field name="dates" type="array" required="false">
            <Description>Danh sách bộ đếm theo từng ngày</Description>
            <ArrayItems>
              <Fields>
                <Field name="company_id" type="integer" required="false">
                  Mã công ty
                </Field>
                <Field name="date" type="string (date)" required="false">
                  Ngày (YYYY-MM-DD)
                </Field>
                <Field name="transaction" type="integer" required="false">
                  Tổng số giao dịch (chỉ đếm tiền vào, bằng transaction_in)
                </Field>
                <Field name="transaction_in" type="integer" required="false">
                  Số giao dịch tiền vào
                </Field>
                <Field name="transaction_out" type="integer" required="false">
                  Số giao dịch tiền ra
                </Field>
              </Fields>
            </ArrayItems>
          </Field>
          <Field name="total" type="object" required="false">
            <Description>Tổng cộng tất cả các ngày</Description>
            <Fields>
              <Field name="transaction" type="integer" required="false">
                Tổng transaction của tất cả các ngày (bằng transaction_in)
              </Field>
              <Field name="transaction_in" type="integer" required="false">
                Tổng transaction_in của tất cả các ngày
              </Field>
              <Field name="transaction_out" type="integer" required="false">
                Tổng transaction_out của tất cả các ngày
              </Field>
            </Fields>
          </Field>
        </Fields>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

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
      --url 'https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03' \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03",
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
    
    conn = http.client.HTTPSConnection("bankhub-api-sandbox.sepay.vn")
    
    headers = { 'Authorization': "Bearer REPLACE_BEARER_TOKEN" }
    
    conn.request("GET", "/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03", headers=headers)
    
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
      "hostname": "bankhub-api-sandbox.sepay.vn",
      "port": null,
      "path": "/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/counter?start_date=2026-04-01&end_date=2026-04-03")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Bước tiếp theo

Sau khi lấy được thống kê bộ đếm giao dịch của merchant:

1. **[API Lấy cấu hình Merchant](/vi/bankhub/api/api-merchant/thong-tin-cau-hinh)** - Xem cấu hình hiện tại của merchant bao gồm whitelist keywords và logo
2. **[API Chỉnh sửa cấu hình Merchant](/vi/bankhub/api/api-merchant/chinh-sua-thong-tin-cau-hinh)** - Cập nhật cấu hình merchant

<Callout type="info" title="Gợi ý">
Dùng bộ đếm này để theo dõi lưu lượng giao dịch theo ngày và phát hiện bất thường trong hoạt động thanh toán của merchant.
</Callout>