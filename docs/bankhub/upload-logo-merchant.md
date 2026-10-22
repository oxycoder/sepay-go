# API upload logo Merchant

## POST /merchants/logo — upload logo merchant vào SePay Bank Hub qua multipart/form-data. Logo hiển thị ngay trên header mọi hosted link đang chạy.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo</Path>

  <Description>
    Upload Logo Merchant
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Request Body (multipart/form-data)

<Params>
</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Content-Type phải là 
`multipart/form-data`
Logo cũ (nếu có) sẽ tự động bị xóa khỏi S3
Logo mới được upload lên S3 với tên file: 
`{merchant_id}-{timestamp}.png`
Logo hiển thị ngay lập tức trên tất cả hosted link
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Bad Request", description: "File không hợp lệ (định dạng sai, kích thước quá lớn, MIME type không khớp)" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 413, name: "Payload Too Large", description: "File vượt quá 1MB" }
]}
/>

## API Response

<Responses>
  <Response status="200">
    <Description>
      Upload thành công
    </Description>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="200">
    <Fields>
      <Field name="data" type="object" required="false">
        <Description>Dữ liệu logo</Description>
        <Fields>
          <Field name="logo_url" type="string (uri)" required="false">
            URL logo mới sau khi upload
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: multipart/form-data' \
      --form logo=string
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n",
      CURLOPT_HTTPHEADER => [
        "Authorization: Bearer REPLACE_BEARER_TOKEN",
        "content-type: multipart/form-data"
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
    
    payload = "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "multipart/form-data"
        }
    
    conn.request("POST", "/v1/merchant-config/logo", payload, headers)
    
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
      "path": "/v1/merchant-config/logo",
      "headers": {
        "Authorization": "Bearer REPLACE_BEARER_TOKEN",
        "content-type": "multipart/form-data"
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
    
    req.write("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n");
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("multipart/form-data; boundary=---011000010111000001101001");
    RequestBody body = RequestBody.create(mediaType, "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "multipart/form-data")
      .build();
    
    Response response = client.newCall(request).execute();
    ```
  </CodeSample>

  <CodeSample value="ruby_native" lang="ruby">
    ```ruby
    require 'uri'
    require 'net/http'
    require 'openssl'
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'multipart/form-data'
    request.body = "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo"
    
    	payload := strings.NewReader("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n")
    
    	req, _ := http.NewRequest("POST", url, payload)
    
    	req.Header.Add("Authorization", "Bearer REPLACE_BEARER_TOKEN")
    	req.Header.Add("content-type", "multipart/form-data")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new MultipartFormDataContent
        {
            new StringContent("string")
            {
                Headers =
                {
                    ContentDisposition = new ContentDispositionHeaderValue("form-data")
                    {
                        Name = "logo",
                    }
                }
            },
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
    
    let headers = [
      "Authorization": "Bearer REPLACE_BEARER_TOKEN",
      "content-type": "multipart/form-data"
    ]
    let parameters = [
      [
        "name": "logo",
        "value": "string"
      ]
    ]
    
    let boundary = "---011000010111000001101001"
    
    var body = ""
    var error: NSError? = nil
    for param in parameters {
      let paramName = param["name"]!
      body += "--\(boundary)\r\n"
      body += "Content-Disposition:form-data; name=\"\(paramName)\""
      if let filename = param["fileName"] {
        let contentType = param["content-type"]!
        let fileContent = String(contentsOfFile: filename, encoding: String.Encoding.utf8)
        if (error != nil) {
          print(error)
        }
        body += "; filename=\"\(filename)\"\r\n"
        body += "Content-Type: \(contentType)\r\n\r\n"
        body += fileContent
      } else if let paramValue = param["value"] {
        body += "\r\n\r\n\(paramValue)"
      }
    }
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo")! as URL,
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
    
    val mediaType = MediaType.parse("multipart/form-data; boundary=---011000010111000001101001")
    val body = RequestBody.create(mediaType, "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"logo\"\r\n\r\nstring\r\n-----011000010111000001101001--\r\n")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "multipart/form-data")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

### Ví dụ sử dụng

<Node title="Upload logo từ file">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/merchant-config/logo' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--form 'logo=@"/path/to/logo.png"'
```
</Node>

***

<Callout type="warn" title="Validation Rules">
File Requirements:
Format
: Chỉ chấp nhận JPG hoặc PNG
Size
: Tối đa 1MB (1,048,576 bytes)
MIME Type
: Phải khớp với extension (image/jpeg hoặc image/png)
Extension
: File extension được kiểm tra (.jpg, .jpeg, .png)
</Callout>

<Callout type="warn" title="Lưu ý quan trọng">
Logo cũ tự động bị xóa, không thể rollback
Logo mới hiển thị ngay lập tức trên tất cả hosted link
Timestamp trong URL giúp bypass browser cache
Nếu file > 1MB, API trả về 400
Chỉ chấp nhận JPG/PNG, không chấp nhận GIF, WebP, SVG
</Callout>