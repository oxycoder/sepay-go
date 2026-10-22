# API tạo Link Token

## POST /link-token/create — tạo Link Token một lần để mở WebView SePay Bank Hub cho luồng liên kết hoặc hủy liên kết tài khoản với redirect URI tùy biến.

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
  <Method>POST</Method>

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/link-token/create</Path>

  <Description>
    Tạo Link Token
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Request Body

<Params>
  <RequestBody>
    <Fields>
      <Field name="company_xid" type="string (uuid)" required="true">
        Company XID - công ty thực hiện liên kết/hủy liên kết
      </Field>
      <Field name="purpose" type="enum: LINK_BANK_ACCOUNT, UNLINK_BANK_ACCOUNT, REACTIVATE_BANK_ACCOUNT" required="true">
        Mục đích tạo link token:
• LINK_BANK_ACCOUNT: liên kết tài khoản mới
• UNLINK_BANK_ACCOUNT: hủy liên kết (vô hiệu hóa) tài khoản
• REACTIVATE_BANK_ACCOUNT: khôi phục tài khoản đã bị vô hiệu hóa trong Bank Hub
      </Field>
      <Field name="completion_redirect_uri" type="string (uri)" required="false">
        URL chuyển hướng sau khi hoàn thành flow
      </Field>
      <Field name="is_mobile_app" type="enum: 0, 1" required="false">
        Flow được mở từ mobile app (1) hay web (0)
      </Field>
      <Field name="language" type="enum: vi, en" required="false">
        Ngôn ngữ hiển thị giao diện
      </Field>
      <Field name="bank_account_xid" type="string (uuid)" required="false">
        Bank Account XID. Bắt buộc khi purpose = UNLINK_BANK_ACCOUNT hoặc REACTIVATE_BANK_ACCOUNT
      </Field>
    </Fields>
  </RequestBody>
</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Link token có thời gian hiệu lực giới hạn (xem trường 
`expires_at`
)
Sử dụng 
`hosted_link_url`
 để mở WebView cho người dùng
Với purpose 
`UNLINK_BANK_ACCOUNT`
, bắt buộc phải có 
`bank_account_xid`
Với purpose 
`REACTIVATE_BANK_ACCOUNT`
, bắt buộc phải có 
`bank_account_xid`
; tài khoản phải đang thuộc diện khôi phục (đã bị ngừng kích hoạt qua luồng hủy liên kết trước đó)
Với purpose 
`LINK_BANK_ACCOUNT`
, không cần truyền 
`bank_account_xid`
`bank_account_xid`
 có thể lấy từ danh sách ngân hàng đã liên kết 
API danh sách ngân hàng
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Validation Error", description: "Dữ liệu đầu vào không hợp lệ (company_xid không tồn tại, purpose sai format, thiếu bank_account_xid khi UNLINK hoặc REACTIVATE, tài khoản không thuộc diện khôi phục)" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 404, name: "Not Found", description: "Công ty hoặc tài khoản ngân hàng không tồn tại" }
]}
/>

## API Response

<Responses>
  <Response status="201">
    <Description>
      Tạo link token thành công
    </Description>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="201">
    <Fields>
      <Field name="xid" type="string (uuid)" required="false">
        ID duy nhất của link token
      </Field>
      <Field name="hosted_link_url" type="string (uri)" required="false">
        URL WebView hoàn chỉnh để mở trong iframe hoặc SDK
      </Field>
      <Field name="link_token" type="string (uuid)" required="false">
        Token duy nhất để xác thực
      </Field>
      <Field name="expires_at" type="string" required="false">
        Thời điểm hết hạn (yyyy-mm-dd HH:mm:ss)
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

<Responses>
  <Response status="400">
    <Description>
      Validation error
    </Description>
  </Response>

  <Response status="401">
    <Description>
      Access token không hợp lệ hoặc đã hết hạn
    </Description>
  </Response>

  <Response status="404">
    <Description>
      Công ty hoặc tài khoản ngân hàng không tồn tại
    </Description>
  </Response>

</Responses>

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
      --url https://bankhub-api-sandbox.sepay.vn/v1/link-token/create \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"company_xid":"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2","purpose":"LINK_BANK_ACCOUNT","completion_redirect_uri":"http://example.com","is_mobile_app":0,"language":"vi","bank_account_xid":"52930e2b-c38b-4cd1-b2a2-9ddcdc603104"}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/link-token/create",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}",
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
    
    conn = http.client.HTTPSConnection("bankhub-api-sandbox.sepay.vn")
    
    payload = "{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/link-token/create", payload, headers)
    
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
      "hostname": "bankhub-api-sandbox.sepay.vn",
      "port": null,
      "path": "/v1/link-token/create",
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
      company_xid: 'd0424f2d-8a2c-4d55-aeed-a5ad9955f6a2',
      purpose: 'LINK_BANK_ACCOUNT',
      completion_redirect_uri: 'http://example.com',
      is_mobile_app: 0,
      language: 'vi',
      bank_account_xid: '52930e2b-c38b-4cd1-b2a2-9ddcdc603104'
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/link-token/create")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/link-token/create")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/link-token/create"
    
    	payload := strings.NewReader("{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/link-token/create"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}")
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
      "company_xid": "d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2",
      "purpose": "LINK_BANK_ACCOUNT",
      "completion_redirect_uri": "http://example.com",
      "is_mobile_app": 0,
      "language": "vi",
      "bank_account_xid": "52930e2b-c38b-4cd1-b2a2-9ddcdc603104"
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/link-token/create")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"company_xid\":\"d0424f2d-8a2c-4d55-aeed-a5ad9955f6a2\",\"purpose\":\"LINK_BANK_ACCOUNT\",\"completion_redirect_uri\":\"http://example.com\",\"is_mobile_app\":0,\"language\":\"vi\",\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\"}")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/link-token/create")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Ví dụ sử dụng

<Node title="Tạo link token để liên kết tài khoản mới">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/link-token/create' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "company_xid": "d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
  "purpose": "LINK_BANK_ACCOUNT",
  "completion_redirect_uri": "https://yourapp.com/success"
}'
```
</Node>

<Node title="Tạo link token để hủy liên kết tài khoản">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/link-token/create' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "company_xid": "d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
  "purpose": "UNLINK_BANK_ACCOUNT",
  "bank_account_xid": "f414b73c-eebe-11f0-b16e-52c7e9b4f41b",
  "completion_redirect_uri": "https://yourapp.com/unlinked"
}'
```
</Node>

<Node title="Tạo link token để khôi phục tài khoản ngừng kích hoạt">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/link-token/create' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "company_xid": "d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
  "purpose": "REACTIVATE_BANK_ACCOUNT",
  "bank_account_xid": "f414b73c-eebe-11f0-b16e-52c7e9b4f41b",
  "completion_redirect_uri": "https://yourapp.com/reactivated"
}'
```
</Node>

***

<Callout type="warning" title="Lưu ý quan trọng">
Expiration
: Link token có thời gian hiệu lực giới hạn. Kiểm tra 
`expires_at`
 và tạo token mới nếu hết hạn
One-time Use
: Link token chỉ sử dụng một lần. Sau khi hoàn thành (thành công hoặc thất bại), cần tạo token mới
Redirect URI
: URL 
`completion_redirect_uri`
 phải là HTTPS trong production
Validation
: Với purpose 
`UNLINK_BANK_ACCOUNT`
 hoặc 
`REACTIVATE_BANK_ACCOUNT`
, nếu thiếu 
`bank_account_xid`
 sẽ trả về lỗi 400
Diện khôi phục
: Purpose 
`REACTIVATE_BANK_ACCOUNT`
 chỉ áp dụng cho tài khoản đang ngừng kích hoạt (bị hủy liên kết kiểu ngừng kích hoạt trước đó); tài khoản không thuộc diện này sẽ trả về lỗi 400
</Callout>

<Callout type="info" title="Gợi ý sử dụng">
Backend tạo link token qua API này
Backend trả 
`hosted_link_url`
 cho frontend
Frontend mở WebView (iframe hoặc SDK) với URL
Người dùng hoàn tất liên kết/hủy liên kết trong WebView
Sau khi xong, người dùng được redirect về 
`completion_redirect_uri`
 (nếu có)
Backend nhận webhook notification về kết quả
</Callout>

***

## Bước tiếp theo

Sau khi tạo link token thành công:

1. **[SDK JavaScript](/vi/bankhub/sdk)** - Sử dụng SDK để mở Hosted Link và quản lý sự kiện
2. **[Webhook Events](/vi/bankhub/webhook-event)** - Nhận thông báo khi liên kết thành công
3. **[Thông tin Link Token](/vi/bankhub/api/link-token/thong-tin-link-token)** - Kiểm tra trạng thái link token
4. **[Danh sách tài khoản ngân hàng](/vi/bankhub/api/tai-khoan-ngan-hang/danh-sach-tai-khoan)** - Xem các tài khoản đã liên kết