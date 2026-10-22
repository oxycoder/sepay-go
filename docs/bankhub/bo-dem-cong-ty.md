# API thống kê giao dịch Công ty

## GET /company/counter/{id} — số liệu giao dịch theo ngày của một công ty Bank Hub, lọc theo khoảng thời gian, tổng hợp giao dịch vào và ra.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/company/counter/{company_id}</Path>

  <Description>
    Đếm giao dịch theo công ty
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Path Parameters

<Params>
  <PathParams>
    <Param name="company_id" type="string" required="true">
      ID của công ty
    </Param>
  </PathParams>

  <QueryParams>
    <Param name="date" type="string" required="false">
      Lọc theo ngày cụ thể (YYYY-MM-DD)
    </Param>
  </QueryParams>

</Params>

## Query Parameters

<Params>
  <PathParams>
    <Param name="company_id" type="string" required="true">
      ID của công ty
    </Param>
  </PathParams>

  <QueryParams>
    <Param name="date" type="string" required="false">
      Lọc theo ngày cụ thể (YYYY-MM-DD)
    </Param>
  </QueryParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Công ty phải thuộc về merchant đang xác thực
Kết quả trả về thống kê theo ngày và tổng hợp
Các trường counter bao gồm: transaction (tổng), transaction_in (vào), transaction_out (ra)
</Callout>

## Xử lý lỗi

<ErrorCodes
  hiddenHead={true}
  rows={[
{ code: 404, name: "Not Found", description: "Công ty không tồn tại hoặc không thuộc merchant này" },
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
        <Description>Dữ liệu thống kê</Description>
        <Fields>
          <Field name="dates" type="array" required="false">
            <Description>Thống kê theo ngày</Description>
            <ArrayItems>
              <Fields>
                <Field name="date" type="string (date)" required="false">
                  Ngày
                </Field>
                <Field name="transaction" type="integer" required="false">
                  Tổng số giao dịch
                </Field>
                <Field name="transaction_in" type="integer" required="false">
                  Số giao dịch tiền vào
                </Field>
                <Field name="transaction_out" type="integer" required="false">
                  Số giao dịch tiền ra
                </Field>
              </Fields>
            </ArrayItems>
          </Field>
          <Field name="total" type="object" required="false">
            <Description>Tổng cộng</Description>
            <Fields>
              <Field name="transaction" type="integer" required="false">
                Tổng số giao dịch
              </Field>
              <Field name="transaction_in" type="integer" required="false">
                Tổng giao dịch tiền vào
              </Field>
              <Field name="transaction_out" type="integer" required="false">
                Tổng giao dịch tiền ra
              </Field>
            </Fields>
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
      --url 'https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE' \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE",
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
    
    conn.request("GET", "/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE", headers=headers)
    
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
      "path": "/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company/counter/d3dafd01-e06b-11f0-b29e-52c7e9b4f41b?date=SOME_STRING_VALUE")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

<Callout type="warn" title="Lưu ý quan trọng">
Transaction In/Out
: 
`transaction_in`
 là giao dịch tiền vào tài khoản, 
`transaction_out`
 là giao dịch tiền ra
Filtering
: Sử dụng query parameter 
`date`
 để lọc theo ngày cụ thể, nếu không truyền sẽ lấy tất cả
Total
: Trường 
`total`
 luôn chứa tổng hợp của tất cả các ngày trong 
`dates`
Empty Dates
: Nếu chưa có giao dịch, mảng 
`dates`
 sẽ rỗng và 
`total`
 sẽ có giá trị 0
</Callout>