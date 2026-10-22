# API lấy danh sách Công ty

## GET /company — lấy danh sách công ty trong SePay Bank Hub có phân trang và lọc theo trạng thái để quản lý cấu trúc merchant nhiều công ty.

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

  <Path>https://bankhub-api-sandbox.sepay.vn/v1/company</Path>

  <Description>
    Danh sách Công ty
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

## Query Parameters

<Params>
  <QueryParams>
    <Param name="page" type="integer" required="false">
      Số trang cần lấy
    </Param>
    <Param name="per_page" type="integer" required="false">
      Số item mỗi trang
    </Param>
    <Param name="q" type="string" required="false">
      Từ khóa tìm kiếm (tìm trong full_name)
    </Param>
    <Param name="status" type="enum: Pending, Active, Suspended, Terminated, Cancelled, Fraud" required="false">
      Lọc theo trạng thái
    </Param>
    <Param name="sort[created_at]" type="enum: asc, desc" required="false">
      Sắp xếp theo created_at
    </Param>
  </QueryParams>

</Params>

<Callout type="info" title="Lưu ý">
API này yêu cầu Bearer Token trong header Authorization
Kết quả được sắp xếp theo thời gian tạo mới nhất
Trường 
`active`
 cho biết công ty có đang hoạt động hay không (1 = Active, 0 = Pending)
</Callout>

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
      <Field name="data" type="array" required="false">
        <Description>Danh sách công ty</Description>
        <ArrayItems>
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
              1 = Active, 0 = Pending
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
      <Field name="meta" type="object" required="false">
        <Description>Thông tin phân trang</Description>
        <Fields>
          <Field name="per_page" type="integer" required="false">
            Số item mỗi trang
          </Field>
          <Field name="total" type="integer" required="false">
            Tổng số item
          </Field>
          <Field name="has_more" type="boolean" required="false">
            Còn trang tiếp theo không
          </Field>
          <Field name="current_page" type="integer" required="false">
            Trang hiện tại
          </Field>
          <Field name="page_count" type="integer" required="false">
            Tổng số trang
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
      --url 'https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE' \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE",
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
    
    conn.request("GET", "/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE", headers=headers)
    
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
      "path": "/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE",
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE")
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
    
    url = URI("https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE")
    
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
    
    	url := "https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE"
    
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
        RequestUri = new Uri("https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE")! as URL,
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
      .url("https://bankhub-api-sandbox.sepay.vn/v1/company?page=SOME_INTEGER_VALUE&per_page=SOME_INTEGER_VALUE&q=SOME_STRING_VALUE&status=SOME_STRING_VALUE&sort%5Bcreated_at%5D=SOME_STRING_VALUE")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

***

## Ví dụ sử dụng

<Node title="Lấy trang đầu tiên với 10 công ty">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/company?per_page=10&page=1' \\\n  --header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Node title="Phân trang - Lấy trang 2">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/company?page=2' \\\n  --header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Node title="Lọc chỉ công ty đang hoạt động">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/company?status=Active' \\\n  --header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Node title="Lọc chỉ công ty với trạng thái đang chờ">
```js
curl --location 'https://bankhub-api-sandbox.sepay.vn/v1/company?status=Pending' \\\n  --header 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```
</Node>

<Callout type="warn" title="Lưu ý về Performance">
Nên sử dụng 
`per_page`
 phù hợp để tối ưu hiệu suất
Không nên lấy quá nhiều công ty trong một lần request
Sử dụng 
`has_more`
 để kiểm tra còn dữ liệu hay không thay vì dựa vào 
`page_count`
</Callout>