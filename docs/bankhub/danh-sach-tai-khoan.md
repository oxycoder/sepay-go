# API lấy danh sách Tài khoản

## GET /bank-account — lấy danh sách mọi tài khoản ngân hàng đã liên kết qua SePay Bank Hub với phân trang và lọc theo công ty, ngân hàng và từ khóa.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/bank-account</Path>

  <Description>
    Danh sách tài khoản ngân hàng
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Query Parameters

<Params>
  <QueryParams>
    <Param name="page" type="integer" required="false">
      Số trang cần lấy
    </Param>
    <Param name="per_page" type="integer" required="false">
      Số item mỗi trang
    </Param>
    <Param name="company_xid" type="string" required="false">
      Lọc theo Company XID
    </Param>
    <Param name="bank_id" type="string" required="false">
      Lọc theo mã ngân hàng
    </Param>
    <Param name="q" type="string" required="false">
      Từ khóa tìm kiếm
    </Param>
  </QueryParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Kết quả được sắp xếp theo thời gian tạo mới nhất
Sử dụng 
`company_xid`
 để lọc tài khoản theo công ty cụ thể
Phân trang với 
`page`
 và 
`per_page`
 (mặc định 20 items/trang)
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" }
]}
/>

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
      <Field name="data" type="array" required="false">
        <Description>Danh sách tài khoản ngân hàng</Description>
        <ArrayItems>
          <Fields>
            <Field name="xid" type="string (uuid)" required="false">
              ID tài khoản ngân hàng
            </Field>
            <Field name="company_xid" type="string (uuid)" required="false">
              ID công ty sở hữu
            </Field>
            <Field name="brand_name" type="string" required="false">
              Tên ngân hàng (OCB, Vietcombank, etc.)
            </Field>
            <Field name="bank_bin" type="string" required="false">
              Mã BIN ngân hàng theo chuẩn NAPAS (ví dụ 970416). Null nếu ngân hàng chưa có mã BIN
            </Field>
            <Field name="bank_code" type="string" required="false">
              Mã ngân hàng (ví dụ ACB). Null nếu ngân hàng chưa có mã
            </Field>
            <Field name="account_holder_name" type="string" required="false">
              Tên chủ tài khoản
            </Field>
            <Field name="account_number" type="string" required="false">
              Số tài khoản
            </Field>
            <Field name="accumulated" type="string" required="false">
              Số dư tích lũy
            </Field>
            <Field name="account_type" type="enum: individual, enterprise" required="false">
              Loại tài khoản (cá nhân/doanh nghiệp)
            </Field>
            <Field name="bank_api_connected" type="string" required="false">
              1 = Đã kết nối, 0 = Chưa kết nối
            </Field>
            <Field name="active" type="string" required="false">
              1 = Active, 0 = Inactive
            </Field>
            <Field name="last_transaction" type="string" required="false">
              Thời điểm giao dịch cuối
            </Field>
            <Field name="created_at" type="string" required="false">
              Thời điểm tạo
            </Field>
            <Field name="updated_at" type="string" required="false">
              Thời điểm cập nhật
            </Field>
          </Fields>
        </ArrayItems>
      </Field>
      <Field name="meta" type="object" required="false">
        <Description>Thông tin phân trang</Description>
        <Fields>
          <Field name="per_page" type="integer" required="false">
            Số item mỗi trang
          </Field>
          <Field name="total" type="integer" required="false">
            Tổng số item
          </Field>
          <Field name="has_more" type="boolean" required="false">
            Còn trang tiếp theo không
          </Field>
          <Field name="current_page" type="integer" required="false">
            Trang hiện tại
          </Field>
          <Field name="page_count" type="integer" required="false">
            Tổng số trang
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
      --url 'https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE' \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE",
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
    
    conn.request("GET", "/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE", headers=headers)
    
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
      "path": "/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&company_xid=SOME_STRING_VALUE&bank_id=SOME_STRING_VALUE&q=SOME_STRING_VALUE")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Ví dụ sử dụng

<Node title="Lấy trang đầu tiên với 10 tài khoản">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/bank-account?per_page=10&page=1' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Node title="Lọc theo công ty">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/bank-account?company_xid=d3dafd01-e06b-11f0-b29e-52c7e9b4f41b' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Node title="Phân trang - Lấy trang 2">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/bank-account?page=2' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

***

<Callout type="warn" title="Lưu ý về Performance">
Nên sử dụng 
`per_page`
 phù hợp để tối ưu hiệu suất
Với số lượng tài khoản lớn, nên phân trang thay vì lấy tất cả
Sử dụng 
`company_xid`
 để lọc theo công ty giúp giảm tải dữ liệu
</Callout>