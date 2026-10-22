# API lấy chi tiết Tài khoản

## GET /bank-account/{id} — lấy chi tiết một tài khoản ngân hàng đã liên kết trong SePay Bank Hub, gồm số tài khoản, chủ, trạng thái và Virtual Account.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/bank-account/{bank_account_id}</Path>

  <Description>
    Chi tiết tài khoản ngân hàng
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Path Parameters

<Params>
  <PathParams>
    <Param name="bank_account_id" type="string" required="true">
      ID của tài khoản ngân hàng
    </Param>
  </PathParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Tài khoản phải thuộc về merchant đang xác thực
Response bao gồm cả danh sách tài khoản ảo (Virtual Account) nếu có
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 401, name: "Unauthorized", description: "Access token không hợp lệ hoặc đã hết hạn" },
  { code: 404, name: "Not Found", description: "Tài khoản không tồn tại hoặc không thuộc merchant này" }
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
        <Description>Thông tin chi tiết tài khoản</Description>
        <Fields>
          <Field name="xid" type="string (uuid)" required="false">
            ID tài khoản
          </Field>
          <Field name="company_xid" type="string (uuid)" required="false">
            ID công ty
          </Field>
          <Field name="brand_name" type="string" required="false">
            Tên ngân hàng
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
          <Field name="bank_api_connected" type="string" required="false">
            Trạng thái kết nối API ngân hàng
          </Field>
          <Field name="active" type="string" required="false">
            Trạng thái hoạt động
          </Field>
          <Field name="last_transaction" type="string" required="false">
            Thời điểm giao dịch cuối
          </Field>
          <Field name="identification_number" type="string" required="false">
            Số CMND/CCCD
          </Field>
          <Field name="phone_number" type="string" required="false">
            Số điện thoại
          </Field>
          <Field name="created_at" type="string" required="false">
            Thời điểm tạo
          </Field>
          <Field name="updated_at" type="string" required="false">
            Thời điểm cập nhật
          </Field>
          <Field name="va_list" type="array" required="false">
            <Description>Danh sách tài khoản ảo (Virtual Account)</Description>
            <ArrayItems>
              <Fields>
                <Field name="xid" type="string (uuid)" required="false">
                  ID tài khoản ảo
                </Field>
                <Field name="sub_account" type="string" required="false">
                  Số tài khoản ảo
                </Field>
                <Field name="sub_holder_name" type="string" required="false">
                  Tên chủ tài khoản ảo
                </Field>
                <Field name="active" type="string" required="false">
                  Trạng thái hoạt động
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
      --url https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b",
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
    
    conn.request("GET", "/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b", headers=headers)
    
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
      "path": "/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## Ví dụ sử dụng

<Node title="Lấy chi tiết tài khoản ngân hàng">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/bank-account/780717f8-eeb6-11f0-b16e-52c7e9b4f41b' \\
--header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Callout type="warn" title="Lưu ý quan trọng">
Virtual Account (VA)
: Tài khoản ảo là các tài khoản con được tạo từ tài khoản ngân hàng chính, giúp phân biệt nguồn tiền từ các khách hàng khác nhau
Identification Number
: Số CMND/CCCD phục vụ mục đích xác thực và tuân thủ quy định
Phone Number
: Số điện thoại dùng để nhận thông báo và xác thực giao dịch
Empty VA List
: Nếu 
`va_list`
 rỗng (
`[]`
), tài khoản chưa có tài khoản ảo nào được tạo
</Callout>