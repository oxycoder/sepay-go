# API lấy chi tiết Công ty

## GET /company/{id} — lấy hồ sơ chi tiết một công ty SePay Bank Hub, gồm tên, trạng thái, cấu hình và giới hạn giao dịch theo từng công ty.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/company/{company_id}</Path>

  <Description>
    Chi tiết Công ty
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Path Parameters

<Params>
  <PathParams>
    <Param name="company_id" type="string" required="true">
      ID của công ty (xid format)
    </Param>
  </PathParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Trường 
`configuration`
 chứa các cấu hình đặc biệt của công ty (nếu có)
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
{ code: 404, name: "Not Found", description: "Công ty (tổ chức) không tồn tại hoặc không thuộc merchant này" },
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
      <Field name="data" type="object" required="false">
        <Description>Thông tin chi tiết công ty</Description>
        <Fields>
          <Field name="xid" type="string (uuid)" required="false">
            ID công ty
          </Field>
          <Field name="full_name" type="string" required="false">
            Tên đầy đủ
          </Field>
          <Field name="status" type="string" required="false">
            Trạng thái
          </Field>
          <Field name="active" type="string" required="false">
            Trạng thái hoạt động
          </Field>
          <Field name="configuration" type="object" required="false">
            <Description>Cấu hình công ty</Description>
            <Fields>
              <Field name="payment_code" type="string" required="false">
                Cấu hình mã thanh toán
              </Field>
              <Field name="payment_code_prefix" type="string" required="false">
                Tiền tố mã thanh toán
              </Field>
              <Field name="payment_code_suffix_from" type="integer" required="false">
                Độ dài tối thiểu hậu tố
              </Field>
              <Field name="payment_code_suffix_to" type="integer" required="false">
                Độ dài tối đa hậu tố
              </Field>
              <Field name="payment_code_suffix_character_type" type="string" required="false">
                Loại ký tự hậu tố
              </Field>
              <Field name="data_storage_time" type="string" required="false">
                Thời gian lưu trữ dữ liệu
              </Field>
              <Field name="transaction_amount" type="string" required="false">
                Giới hạn số tiền giao dịch
              </Field>
            </Fields>
          </Field>
          <Field name="created_at" type="string" required="false">
            Thời điểm tạo
          </Field>
          <Field name="updated_at" type="string" required="false">
            Thời điểm cập nhật
          </Field>
        </Fields>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

<Responses>
  <Response status="401">
    <Description>
      Access token không hợp lệ hoặc đã hết hạn
    </Description>
  </Response>

  <Response status="404">
    <Description>
      Công ty không tồn tại hoặc không thuộc merchant này
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
    curl --request GET \
      --url https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
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
    
    conn.request("GET", "/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b", headers=headers)
    
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
      "path": "/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

<Callout type="warn" title="Xử lý lỗi 404">
Lỗi 404 có thể xảy ra khi:
Company ID không tồn tại trong hệ thống
Company ID không thuộc về merchant đang xác thực
Company ID bị nhập sai format
Luôn kiểm tra status code trước khi xử lý response
</Callout>