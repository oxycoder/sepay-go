# API giả lập giao dịch

## POST /transaction/create — tạo giao dịch giả lập trong SePay Bank Hub sandbox để kiểm tra luồng xử lý và webhook IPN mà không cần giao dịch thật.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/transaction/create</Path>

  <Description>
    Giả lập giao dịch (Sandbox)
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Request Body

<Params>
  <RequestBody>
    <Fields>
      <Field name="bank_account_xid" type="string (uuid)" required="true">
        ID tài khoản ngân hàng
      </Field>
      <Field name="transfer_type" type="enum: credit, debit" required="true">
        Loại giao dịch (credit = tiền vào, debit = tiền ra)
      </Field>
      <Field name="amount" type="number" required="true">
        Số tiền giao dịch (VND)
      </Field>
      <Field name="transaction_content" type="string" required="true">
        Nội dung giao dịch
      </Field>
    </Fields>
  </RequestBody>
</Params>

<Callout type="info" title="Lưu ý">
Chỉ sử dụng trong môi trường Sandbox
Giao dịch giả lập sẽ được dùng để test webhook và luồng xử lý giao dịch
Trường 
`transfer_type`
 dùng để phân biệt tiền vào (
`credit`
) và tiền ra (
`debit`
)
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 400, name: "Validation Error", description: "Dữ liệu đầu vào không hợp lệ" },
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 403, name: "Forbidden", description: "API này chỉ khả dụng trên môi trường Sandbox" }
]}
/>

## API Response

<Responses>
  <Response status="201">
    <Description>
      Tạo giao dịch giả lập thành công
    </Description>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="201">
    <Fields>
      <Field name="code" type="integer" required="false">
        Mã trạng thái HTTP
      </Field>
      <Field name="message" type="string" required="false">
        Thông báo kết quả
      </Field>
      <Field name="data" type="object" required="false">
        <Description>Dữ liệu giao dịch giả lập</Description>
        <Fields>
          <Field name="transaction_id" type="string (uuid)" required="false">
            ID giao dịch
          </Field>
          <Field name="reference_number" type="string" required="false">
            Mã tham chiếu
          </Field>
          <Field name="transaction_date" type="string" required="false">
            Ngày giờ giao dịch
          </Field>
          <Field name="amount" type="number" required="false">
            Số tiền
          </Field>
          <Field name="transfer_type" type="string" required="false">
            Loại giao dịch
          </Field>
          <Field name="accumulated" type="number" required="false">
            Số dư tích lũy sau giao dịch
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/transaction/create \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"bank_account_xid":"52930e2b-c38b-4cd1-b2a2-9ddcdc603104","transfer_type":"credit","amount":0,"transaction_content":"string"}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/transaction/create",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}",
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
    
    payload = "{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/transaction/create", payload, headers)
    
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
      "path": "/v1/transaction/create",
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
      bank_account_xid: '52930e2b-c38b-4cd1-b2a2-9ddcdc603104',
      transfer_type: 'credit',
      amount: 0,
      transaction_content: 'string'
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}");
    Request request = new Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/transaction/create")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/transaction/create")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}"
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/transaction/create"
    
    	payload := strings.NewReader("{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}")
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/transaction/create"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}")
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
      "bank_account_xid": "52930e2b-c38b-4cd1-b2a2-9ddcdc603104",
      "transfer_type": "credit",
      "amount": 0,
      "transaction_content": "string"
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/transaction/create")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"bank_account_xid\":\"52930e2b-c38b-4cd1-b2a2-9ddcdc603104\",\"transfer_type\":\"credit\",\"amount\":0,\"transaction_content\":\"string\"}")
    val request = Request.Builder()
      .url("https://bankhub-api-sandbox.sepay.vn/v1/transaction/create")
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

<Node title="Tạo giao dịch tiền vào">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/transaction/create' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--header 'Content-Type: application/json' \\
--data '{
  "bank_account_xid": "fd770652-f11d-11f0-b21a-a6006ab65aca",
  "transfer_type": "credit",
  "amount": 500000,
  "transaction_content": "Thanh toan DH123456"
}'
```
</Node>

<Node title="Tạo giao dịch tiền ra">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/transaction/create' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \\
--header 'Content-Type: application/json' \\
--data '{
  "bank_account_xid": "fd770652-f11d-11f0-b21a-a6006ab65aca",
  "transfer_type": "debit",
  "amount": 100000,
  "transaction_content": "Chuyen khoan"
}'
```
</Node>

***

<Callout type="info" title="Webhook">
Sau khi tạo giao dịch giả lập, webhook giao dịch sẽ được gửi đến endpoint của bạn nếu đã cấu hình
Xem hướng dẫn cấu hình webhook 
tại đây
</Callout>

<Callout type="warn" title="Chỉ dùng trên Sandbox">
API này 
chỉ khả dụng trên môi trường Sandbox
. Nếu gọi trên Production sẽ nhận được lỗi 403 Forbidden.
</Callout>