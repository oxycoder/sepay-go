# API lấy chi tiết giao dịch

## GET /transaction/{id} — lấy chi tiết một giao dịch theo transaction_id qua SePay Bank Hub. Trả về trạng thái, số tiền, loại giao dịch và metadata.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/transaction/{transaction_id}</Path>

  <Description>
    Chi tiết giao dịch
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Path Parameters

<Params>
  <PathParams>
    <Param name="transaction_id" type="string" required="true">
      Mã giao dịch duy nhất
    </Param>
  </PathParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Transaction ID phải tồn tại trong hệ thống
Trả về thông tin chi tiết đầy đủ của một giao dịch
Bao gồm cả thông tin tài khoản ảo (VA) nếu có
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 404, name: "Not Found", description: "Giao dịch không tồn tại hoặc không thuộc merchant này" }
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
      <Field name="data" type="object" required="false">
        <Fields>
          <Field name="transaction_id" type="string" required="false">
            Mã giao dịch duy nhất
          </Field>
          <Field name="transaction_date" type="string" required="false">
            Ngày giờ giao dịch
          </Field>
          <Field name="bank_account_xid" type="string (uuid)" required="false">
            ID tài khoản ngân hàng
          </Field>
          <Field name="company_xid" type="string (uuid)" required="false">
            ID công ty
          </Field>
          <Field name="brand_name" type="string" required="false">
            Tên ngân hàng
          </Field>
          <Field name="account_number" type="string" required="false">
            Số tài khoản
          </Field>
          <Field name="va" type="string" required="false">
            Số tài khoản ảo (nếu có)
          </Field>
          <Field name="va_xid" type="string (uuid)" required="false">
            ID tài khoản ảo (nếu có)
          </Field>
          <Field name="reference_number" type="string" required="false">
            Mã tham chiếu từ ngân hàng
          </Field>
          <Field name="transfer_type" type="enum: credit, debit" required="false">
            Loại giao dịch (credit = tiền vào, debit = tiền ra)
          </Field>
          <Field name="amount" type="number" required="false">
            Số tiền giao dịch (VND)
          </Field>
          <Field name="transaction_content" type="string" required="false">
            Nội dung chuyển khoản
          </Field>
          <Field name="payment_code" type="string" required="false">
            Mã thanh toán được trích xuất
          </Field>
          <Field name="webhook_success" type="enum: 0, 1" required="false">
            Trạng thái webhook (1 = thành công, 0 = thất bại)
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca",
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
    
    conn.request("GET", "/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca", headers=headers)
    
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
      "path": "/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/transaction/a47b6ab2-8f8b-11f0-b21a-a6006ab65aca")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Ví dụ sử dụng

<Node title="Lấy chi tiết giao dịch">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/transaction/123456' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

***

<Callout type="warning" title="Lưu ý">
Transaction ID
: Phải là ID hợp lệ tồn tại trong hệ thống
404 Error
: Nếu transaction_id không tồn tại, API sẽ trả về lỗi 404
VA Fields
: Trường 
`va`
 và 
`va_xid`
 chỉ có giá trị với giao dịch qua tài khoản ảo
Payment Code
: Được tự động trích xuất từ 
`transaction_content`
 dựa trên cấu hình merchant
</Callout>