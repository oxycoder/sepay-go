# API lấy cấu hình Webhook

## GET /webhooks/config — lấy cấu hình webhook hiện tại của merchant SePay Bank Hub gồm URL đã đăng ký, trạng thái và danh sách sự kiện đang nhận.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/webhook</Path>

  <Description>
    Lấy cấu hình Webhook
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Trả về cấu hình webhook hiện tại của merchant
Nếu chưa cấu hình webhook, trả về 
`data: null`
Secret key được che (hiển thị 4 ký tự cuối)
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
      <Field name="status" type="integer" required="false">
        Mã trạng thái HTTP
      </Field>
      <Field name="error" type="string" required="false">
        Mã lỗi (nếu có)
      </Field>
      <Field name="messages" type="object" required="false">
        <Description>Thông báo kết quả</Description>
        <Fields>
          <Field name="success" type="array" required="false">
            Danh sách thông báo thành công
          </Field>
          <Field name="error" type="array" required="false">
            Danh sách thông báo lỗi
          </Field>
        </Fields>
      </Field>
      <Field name="data" type="object" required="false">
        <Description>Cấu hình webhook (secret key được che)</Description>
        <Fields>
          <Field name="webhook_url" type="string (uri)" required="false">
            URL endpoint webhook
          </Field>
          <Field name="auth_type" type="enum: SECRET_KEY" required="false">
            Loại xác thực
          </Field>
          <Field name="secret_key" type="string" required="false">
            Secret key (được che, chỉ hiển thị 4 ký tự cuối)
          </Field>
          <Field name="request_content_type" type="string" required="false">
            Content type của request
          </Field>
          <Field name="active" type="enum: 0, 1" required="false">
            1 = hoạt động, 0 = tắt
          </Field>
          <Field name="allow_events" type="array" required="false">
            <Description>Danh sách events được đăng ký</Description>
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/webhook \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/webhook",
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
    
    conn.request("GET", "/v1/webhook", headers=headers)
    
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
      "path": "/v1/webhook",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/webhook"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/webhook"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/webhook")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Ví dụ sử dụng

<Node title="Lấy cấu hình webhook">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/webhook' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

***

<Callout type="info" title="Allow Events Configuration">
Cách hoạt động:
`null`
 hoặc 
`[]`
: Không gửi bất kỳ sự kiện nào
`["*"]`
: Gửi tất cả 13 loại sự kiện
`["BANK_ACCOUNT_LINKED", "BANK_ACCOUNT_UNLINKED"]`
: Chỉ gửi các sự kiện cụ thể
13 loại sự kiện:
Link Token: 
`LINK_TOKEN_CREATED`
, 
`LINK_TOKEN_REVOKED`
, 
`LINK_TOKEN_EXPIRED`
Link Session: 
`LINK_SESSION_INITIALIZED`
, 
`LINK_SESSION_REINITIALIZED`
, 
`LINK_SESSION_STATE_CHANGED`
, 
`LINK_SESSION_COMPLETED`
, 
`LINK_SESSION_FAILED`
, 
`LINK_SESSION_ABANDONED`
, 
`LINK_SESSION_EXPIRED`
Bank Account: 
`BANK_ACCOUNT_LINKED`
, 
`BANK_ACCOUNT_UNLINKED`
, 
`BANK_ACCOUNT_INACTIVATED`
</Callout>

<Callout type="info" title="Secret Key Display">
Trong GET response, secret key được che để bảo mật
Chỉ hiển thị 4 ký tự cuối: 
`***5678`
Secret key đầy đủ chỉ trả về khi tạo/cập nhật qua POST
</Callout>

<Callout type="warn" title="Lưu ý quan trọng">
Nếu 
`data`
 là 
`null`
, merchant chưa cấu hình webhook
Secret key được mask trong GET response (hiển thị 
`***xxxx`
)
`allow_events = null`
 hoặc 
`[]`
 nghĩa là webhook sẽ không gửi bất kỳ event nào
Active = 0 nghĩa là webhook bị tắt, không gửi event dù có cấu hình 
`allow_events`
</Callout>