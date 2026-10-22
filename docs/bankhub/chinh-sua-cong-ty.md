# API chỉnh sửa Công ty

## POST /company/edit/{id} — chỉnh sửa thông tin công ty SePay Bank Hub gồm tên, trạng thái và giới hạn giao dịch qua partial update đã xác thực.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/company/edit/{company_id}</Path>

  <Description>
    Chỉnh sửa Công ty
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Path Parameters

<Params>
  <PathParams>
    <Param name="company_id" type="string" required="true">
      ID của công ty cần chỉnh sửa
    </Param>
  </PathParams>

  <RequestBody>
    <Fields>
      <Field name="full_name" type="string" required="false">
        Tên đầy đủ của công ty
      </Field>
      <Field name="short_name" type="string" required="false">
        Tên viết tắt
      </Field>
      <Field name="status" type="enum: Pending, Active, Suspended, Terminated, Cancelled, Fraud" required="false">
        Trạng thái công ty
      </Field>
      <Field name="payment_code" type="enum: on, off" required="false">
        Bật/tắt nhận diện mã thanh toán
      </Field>
      <Field name="payment_code_prefix" type="string" required="false">
        Tiền tố mã thanh toán (2-5 ký tự chữ cái)
      </Field>
      <Field name="payment_code_suffix_from" type="integer" required="false">
        Độ dài tối thiểu hậu tố
      </Field>
      <Field name="payment_code_suffix_to" type="integer" required="false">
        Độ dài tối đa hậu tố
      </Field>
      <Field name="payment_code_suffix_character_type" type="enum: NumberOnly, NumberAndLetter" required="false">
        Loại ký tự hậu tố
      </Field>
      <Field name="transaction_amount" type="object" required="false">
        Số lượng giao dịch
      </Field>
    </Fields>
  </RequestBody>
</Params>

## Request Body

<Params>
  <PathParams>
    <Param name="company_id" type="string" required="true">
      ID của công ty cần chỉnh sửa
    </Param>
  </PathParams>

  <RequestBody>
    <Fields>
      <Field name="full_name" type="string" required="false">
        Tên đầy đủ của công ty
      </Field>
      <Field name="short_name" type="string" required="false">
        Tên viết tắt
      </Field>
      <Field name="status" type="enum: Pending, Active, Suspended, Terminated, Cancelled, Fraud" required="false">
        Trạng thái công ty
      </Field>
      <Field name="payment_code" type="enum: on, off" required="false">
        Bật/tắt nhận diện mã thanh toán
      </Field>
      <Field name="payment_code_prefix" type="string" required="false">
        Tiền tố mã thanh toán (2-5 ký tự chữ cái)
      </Field>
      <Field name="payment_code_suffix_from" type="integer" required="false">
        Độ dài tối thiểu hậu tố
      </Field>
      <Field name="payment_code_suffix_to" type="integer" required="false">
        Độ dài tối đa hậu tố
      </Field>
      <Field name="payment_code_suffix_character_type" type="enum: NumberOnly, NumberAndLetter" required="false">
        Loại ký tự hậu tố
      </Field>
      <Field name="transaction_amount" type="object" required="false">
        Số lượng giao dịch
      </Field>
    </Fields>
  </RequestBody>
</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Hỗ trợ partial update - chỉ gửi các trường cần thay đổi
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
{ code: 404, name: "Not Found", description: "Công ty không tồn tại hoặc không thuộc merchant này" },
{ code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
{ code: 400, name: "Validation Error", description: "Dữ liệu đầu vào không hợp lệ" }
]}
/>

## API Response

<Responses>
  <Response status="200">
    <Description>
      Cập nhật thành công
    </Description>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="200">
    <Fields>
      <Field name="code" type="integer" required="false">
      </Field>
      <Field name="message" type="string" required="false">
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
      Công ty không tồn tại
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"full_name":"string","short_name":"string","status":"Pending","payment_code":"on","payment_code_prefix":"string","payment_code_suffix_from":1,"payment_code_suffix_to":1,"payment_code_suffix_character_type":"NumberOnly","transaction_amount":null}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}",
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
    
    payload = "{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b", payload, headers)
    
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
      "path": "/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
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
      full_name: 'string',
      short_name: 'string',
      status: 'Pending',
      payment_code: 'on',
      payment_code_prefix: 'string',
      payment_code_suffix_from: 1,
      payment_code_suffix_to: 1,
      payment_code_suffix_character_type: 'NumberOnly',
      transaction_amount: null
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b"
    
    	payload := strings.NewReader("{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}")
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
      "full_name": "string",
      "short_name": "string",
      "status": "Pending",
      "payment_code": "on",
      "payment_code_prefix": "string",
      "payment_code_suffix_from": 1,
      "payment_code_suffix_to": 1,
      "payment_code_suffix_character_type": "NumberOnly",
      "transaction_amount": 
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"full_name\":\"string\",\"short_name\":\"string\",\"status\":\"Pending\",\"payment_code\":\"on\",\"payment_code_prefix\":\"string\",\"payment_code_suffix_from\":1,\"payment_code_suffix_to\":1,\"payment_code_suffix_character_type\":\"NumberOnly\",\"transaction_amount\":null}")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
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

<Node title="Cập nhật tên công ty">
```js
curl --location --request POST 'https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "full_name": "Công ty TNHH ABC Updated"
}'
```
</Node>

<Node title="Thay đổi trạng thái công ty thành Pending">
```js
curl --location --request POST 'https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "status": "Pending"
}'
```
</Node>

<Node title="Cập nhật giới hạn giao dịch">
```js
curl --location --request POST 'https://bankhub-api-sandbox.sepay.vn/v1/company/edit/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--data '{
  "transaction_amount": "10000000"
}'
```
</Node>

***

<Callout type="warn" title="Lưu ý quan trọng">
Partial Update
: API này hỗ trợ partial update, chỉ gửi các trường cần thay đổi
Validation
: Các trường gửi lên vẫn phải tuân thủ quy tắc validation (tên tối thiểu 3 ký tự, status đúng format, v.v.)
404 Error
: Kiểm tra company_id trước khi cập nhật để tránh lỗi 404
Transaction Amount
: Có thể set là "Unlimited" hoặc một số cụ thể (string format)
</Callout>