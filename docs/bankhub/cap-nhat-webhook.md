# API tạo hoặc cập nhật Webhook

## PUT /webhooks/config — một lệnh duy nhất để tạo hoặc cập nhật webhook Bank Hub (URL, secret, sự kiện đăng ký), gọi nhiều lần vẫn an toàn.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/webhook</Path>

  <Description>
    Tạo hoặc cập nhật Webhook
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Request Body

<Params>
  <RequestBody>
    <Fields>
      <Field name="webhook_url" type="string (uri)" required="false">
        URL webhook endpoint (HTTPS required in production)
      </Field>
      <Field name="auth_type" type="enum: SECRET_KEY" required="false">
        Loại xác thực
      </Field>
      <Field name="secret_key" type="string" required="false">
        Secret key (sẽ gửi trong header X-Secret-Key)
      </Field>
      <Field name="active" type="enum: 0, 1" required="false">
        1 = bật, 0 = tắt
      </Field>
      <Field name="allow_events" type="array" required="false">
        <Description>["*"] = tất cả events, [] = không events nào</Description>
      </Field>
    </Fields>
  </RequestBody>
</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Pattern upsert: tạo mới nếu chưa có, cập nhật nếu đã tồn tại
Chỉ cho phép một webhook trên mỗi merchant
Secret key đầy đủ được trả về trong response để tham khảo
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Bad Request", description: "URL không hợp lệ, IP nội bộ bị chặn (SSRF protection), hoặc tên sự kiện không hợp lệ" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" }
]}
/>

## API Response

<Responses>
  <Response status="200">
    <Description>
      Upsert thành công
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
        <Description>Cấu hình webhook (secret key đầy đủ - chỉ trong POST response)</Description>
        <Fields>
          <Field name="webhook_url" type="string (uri)" required="false">
            URL endpoint webhook
          </Field>
          <Field name="auth_type" type="enum: SECRET_KEY" required="false">
            Loại xác thực
          </Field>
          <Field name="secret_key" type="string" required="false">
            Secret key đầy đủ (không che)
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
    curl --request POST \
      --url https://bankhub-api-sandbox.sepay.vn/v1/webhook \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"webhook_url":"http://example.com","auth_type":"SECRET_KEY","secret_key":"string","active":0,"allow_events":["*"]}'
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
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}",
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
    
    payload = "{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/webhook", payload, headers)
    
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
      "path": "/v1/webhook",
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
      webhook_url: 'http://example.com',
      auth_type: 'SECRET_KEY',
      secret_key: 'string',
      active: 0,
      allow_events: ['*']
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/webhook"
    
    	payload := strings.NewReader("{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/webhook"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}")
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
      "webhook_url": "http://example.com",
      "auth_type": "SECRET_KEY",
      "secret_key": "string",
      "active": 0,
      "allow_events": ["*"]
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/webhook")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"webhook_url\":\"http://example.com\",\"auth_type\":\"SECRET_KEY\",\"secret_key\":\"string\",\"active\":0,\"allow_events\":[\"*\"]}")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/webhook")
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

<Node title="Đăng ký tất cả sự kiện">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/webhook' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "webhook_url": "https://merchant.com/webhook",
  "auth_type": "SECRET_KEY",
  "secret_key": "your-secret-key-here",
  "active": 1,
  "allow_events": ["*"]
}'
```
</Node>

<Node title="Tắt tất cả sự kiện">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/webhook' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "webhook_url": "https://merchant.com/webhook",
  "allow_events": []
}'
```
</Node>

<Node title="Đăng ký sự kiện cụ thể">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/webhook' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "webhook_url": "https://merchant.com/webhook",
  "auth_type": "SECRET_KEY",
  "secret_key": "your-secret-key-here",
  "active": 1,
  "allow_events": [
    "BANK_ACCOUNT_LINKED",
    "BANK_ACCOUNT_UNLINKED",
    "LINK_TOKEN_CREATED",
    "LINK_SESSION_COMPLETED"
  ]
}'
```
</Node>

***

<Callout type="warn" title="Validation Rules">
Webhook URL:
Phải là HTTPS trong production
Tối đa 2000 ký tự
Không được sử dụng IP nội bộ (SSRF protection)
Dải IP bị chặn: 127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
Secret Key:
Tối đa 500 ký tự
Khuyến nghị: string ngẫu nhiên, độ dài 32-64 ký tự
Allow Events:
Chỉ chấp nhận tên sự kiện hợp lệ từ whitelist (13 events)
`["*"]`
: Wildcard để đăng ký tất cả sự kiện
`[]`
 hoặc 
`null`
: Không gửi bất kỳ sự kiện nào
</Callout>

<Callout type="info" title="13 Webhook Events">
Link Token Events (3):
`LINK_TOKEN_CREATED`
, 
`LINK_TOKEN_REVOKED`
, 
`LINK_TOKEN_EXPIRED`
Link Session Events (7):
`LINK_SESSION_INITIALIZED`
, 
`LINK_SESSION_REINITIALIZED`
, 
`LINK_SESSION_STATE_CHANGED`
`LINK_SESSION_COMPLETED`
, 
`LINK_SESSION_FAILED`
, 
`LINK_SESSION_ABANDONED`
, 
`LINK_SESSION_EXPIRED`
Bank Account Events (3):
`BANK_ACCOUNT_LINKED`
, 
`BANK_ACCOUNT_UNLINKED`
, 
`BANK_ACCOUNT_INACTIVATED`
</Callout>

<Callout type="info" title="Upsert Behavior">
Tạo mới
: Nếu merchant chưa có webhook, tạo webhook mới
Cập nhật
: Nếu merchant đã có webhook, cập nhật cấu hình hiện tại
Một webhook per merchant
: Chỉ được có một webhook config per merchant
Partial update
: Chỉ gửi các field cần update
</Callout>

<Callout type="warn" title="Lưu ý quan trọng">
SSRF Protection
: IP nội bộ và localhost bị chặn để bảo vệ khỏi SSRF attacks
Secret Key Display
: Secret key đầy đủ chỉ hiển thị trong POST response, GET response sẽ che (mask)
Event Validation
: Tên sự kiện không hợp lệ sẽ bị reject với error 400
Active Flag
: Nếu 
`active = 0`
, webhook sẽ không gửi bất kỳ event nào dù có config 
`allow_events`
Immediate Effect
: Thay đổi cấu hình áp dụng ngay lập tức cho tất cả events mới
</Callout>

<Callout type="info" title="Webhook Authentication">
Khi 
`auth_type = "SECRET_KEY"`
, mỗi webhook request sẽ có header:

`X-Secret-Key: your-secret-key-here`
Verify webhook trong endpoint của bạn:
</Callout>

```javascript
app.post('/webhook', (req, res) => {
  const receivedSecret = req.headers['x-secret-key'];
  const expectedSecret = 'your-secret-key-here';

  if (receivedSecret !== expectedSecret) {
    return res.status(401).send('Unauthorized');
  }

  // Process webhook event
  const event = req.body;
  console.log('Event:', event.event);

  res.status(200).send('OK');
});
```