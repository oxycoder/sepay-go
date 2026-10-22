# API lấy chi tiết Link Token

## POST /link-token/get — lấy chi tiết Link Token SePay Bank Hub gồm trạng thái, lịch sử phiên và sự kiện để theo dõi và audit luồng liên kết.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/link-token/get</Path>

  <Description>
    Thông tin Link Token
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Request Body

<Params>
  <RequestBody>
    <Fields>
      <Field name="link_token" type="string (uuid)" required="true">
        Token string của link token (UUID format)
      </Field>
    </Fields>
  </RequestBody>
</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Trả về thông tin chi tiết bao gồm company, sessions và events
Sessions chứa lịch sử các lần người dùng truy cập link token
Events chứa các sự kiện lifecycle của link token
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 404, name: "Not Found", description: "Link token không tồn tại hoặc không thuộc merchant này" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 400, name: "Validation Error", description: "link_token không hợp lệ (format sai hoặc thiếu)" }
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
      <Field name="xid" type="string (uuid)" required="false">
        ID duy nhất của link token
      </Field>
      <Field name="purpose" type="enum: LINK_BANK_ACCOUNT, UNLINK_BANK_ACCOUNT, REACTIVATE_BANK_ACCOUNT" required="false">
        Mục đích của link token
      </Field>
      <Field name="status" type="enum: Active, Revoked, Expired" required="false">
        Trạng thái hiện tại
      </Field>
      <Field name="completion_redirect_uri" type="string" required="false">
        URL chuyển hướng sau khi hoàn thành
      </Field>
      <Field name="is_mobile_app" type="integer" required="false">
        Có phải từ mobile app không
      </Field>
      <Field name="language" type="string" required="false">
        Ngôn ngữ giao diện
      </Field>
      <Field name="session_count" type="integer" required="false">
        Số lần người dùng đã truy cập link token này
      </Field>
      <Field name="created_at" type="string" required="false">
        Thời điểm tạo
      </Field>
      <Field name="updated_at" type="string" required="false">
        Thời điểm cập nhật cuối
      </Field>
      <Field name="expires_at" type="string" required="false">
        Thời điểm hết hạn
      </Field>
      <Field name="company" type="object" required="false">
        <Description>Thông tin công ty</Description>
        <Fields>
          <Field name="full_name" type="string" required="false">
            Tên đầy đủ công ty
          </Field>
          <Field name="status" type="string" required="false">
            Trạng thái công ty
          </Field>
          <Field name="xid" type="string (uuid)" required="false">
            ID công ty
          </Field>
        </Fields>
      </Field>
      <Field name="sessions" type="array" required="false">
        <Description>Danh sách sessions</Description>
        <ArrayItems>
          <Fields>
            <Field name="xid" type="string (uuid)" required="false">
              ID session
            </Field>
            <Field name="context" type="object" required="false">
              Context data của session
            </Field>
            <Field name="state" type="enum: PENDING, IN_PROGRESS, COMPLETED, FAILED" required="false">
              Trạng thái session
            </Field>
            <Field name="status" type="string" required="false">
              Status chi tiết
            </Field>
            <Field name="result" type="object" required="false">
              Kết quả session
            </Field>
            <Field name="finished_at" type="string" required="false">
              Thời điểm hoàn thành
            </Field>
            <Field name="expires_at" type="string" required="false">
              Thời điểm hết hạn
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
      <Field name="events" type="array" required="false">
        <Description>Danh sách events</Description>
        <ArrayItems>
          <Fields>
            <Field name="xid" type="string (uuid)" required="false">
              ID event
            </Field>
            <Field name="name" type="string" required="false">
              Tên event (LINK_TOKEN_CREATED, SESSION_STARTED, etc.)
            </Field>
            <Field name="metadata" type="object" required="false">
              Metadata của event
            </Field>
            <Field name="created_at" type="string" required="false">
              Thời điểm tạo
            </Field>
          </Fields>
        </ArrayItems>
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/link-token/get \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"link_token":"3034ccc8-9cad-4b39-adb3-738c554a7b77"}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/link-token/get",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}",
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
    
    payload = "{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/link-token/get", payload, headers)
    
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
      "path": "/v1/link-token/get",
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
    
    req.write(JSON.stringify({link_token: '3034ccc8-9cad-4b39-adb3-738c554a7b77'}));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/link-token/get")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/link-token/get")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/link-token/get"
    
    	payload := strings.NewReader("{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/link-token/get"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}")
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
    let parameters = ["link_token": "3034ccc8-9cad-4b39-adb3-738c554a7b77"] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/link-token/get")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"link_token\":\"3034ccc8-9cad-4b39-adb3-738c554a7b77\"}")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/link-token/get")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>