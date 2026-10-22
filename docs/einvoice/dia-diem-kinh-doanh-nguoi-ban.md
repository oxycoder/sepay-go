# API địa điểm kinh doanh người bán

## Quản lý danh mục địa điểm kinh doanh người bán (mã số, tên, địa chỉ) qua SePay eInvoice API. Tạo, sửa, xóa địa điểm và lấy ID để truyền mã số và tên địa điểm lên hóa đơn điện tử.

---

**API Overview:**

API tạo và quản lý hóa đơn điện tử theo quy định của Tổng cục Thuế Việt Nam.

**Base URLs:**
- Production: `https://einvoice-api.sepay.vn`
- Sandbox: `https://einvoice-api-sandbox.sepay.vn`


---

<Callout type="info" title="Luồng sử dụng">
Gọi 
API tạo địa điểm
 để thêm địa điểm, nhận 
`xid`
.
Truyền 
`xid`
 đó vào trường 
`seller_store_xid`
 khi 
Xuất hóa đơn điện tử
.
Có thể 
lọc hóa đơn theo địa điểm
 và xem địa điểm trong 
chi tiết hóa đơn
.
</Callout>

## API danh sách địa điểm

Lấy danh sách địa điểm kinh doanh của công ty.

<Endpoint>
  <Method>GET</Method>

  <Path>https://einvoice-api.sepay.vn/v1/business-locations</Path>

  <Description>
    Danh sách địa điểm kinh doanh
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Responses>
  <Response status="200">
    <Description>
      Danh sách địa điểm kinh doanh
    </Description>

    <Example>
      {
        "success": true,
        "data": [
          {
            "xid": "8f14e45f-ceea-4b6d-9c1a-2b3c4d5e6f70",
            "code": "CN01",
            "name": "Chi nhánh 1",
            "address": "12 Đường Ánh Sao, Quận 9, TP.HCM",
            "is_used": true,
            "is_active": true
          },
          {
            "xid": "3c59dc04-8e88-4506-b12a-9f3d1e7c5a21",
            "code": "HQ",
            "name": "Trụ sở chính",
            "address": "1 Đường Ban Mai, Quận 9, TP.HCM",
            "is_used": false,
            "is_active": true
          }
        ]
      }
    </Example>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="200">
    <Fields>
      <Field name="success" type="boolean" required="false">
      </Field>
      <Field name="data" type="array" required="false">
        <ArrayItems>
          <Fields>
            <Field name="xid" type="string (uuid)" required="false">
              Mã định danh địa điểm. Dùng giá trị này làm `seller_store_xid` khi tạo hóa đơn.
            </Field>
            <Field name="code" type="string" required="false">
              Mã số địa điểm kinh doanh (tối đa 50 ký tự). Là định danh cố định, không sửa được sau khi tạo.
            </Field>
            <Field name="name" type="string" required="false">
              Tên địa điểm kinh doanh theo giấy chứng nhận đăng ký kinh doanh (tối đa 400 ký tự)
            </Field>
            <Field name="address" type="string" required="false">
              Địa chỉ địa điểm kinh doanh (tối đa 255 ký tự)
            </Field>
            <Field name="is_used" type="boolean" required="false">
              Cho biết địa điểm đã được dùng để xuất hóa đơn hay chưa. Địa điểm đã dùng (`true`) không xóa được, chỉ có thể hủy kích hoạt bằng `is_active`, nhưng vẫn sửa được tên và địa chỉ bình thường.

            </Field>
            <Field name="is_active" type="boolean" required="false">
              Trạng thái kích hoạt của địa điểm. Địa điểm bị hủy kích hoạt (`false`) không dùng để xuất hóa đơn mới được nữa, nhưng các hóa đơn đã gắn với nó trước đó vẫn tra cứu và lọc được bình thường.

            </Field>
          </Fields>
        </ArrayItems>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

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
      --url https://einvoice-api.sepay.vn/v1/business-locations \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/business-locations",
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
    
    conn = http.client.HTTPSConnection("einvoice-api.sepay.vn")
    
    headers = { 'Authorization': "Bearer REPLACE_BEARER_TOKEN" }
    
    conn.request("GET", "/v1/business-locations", headers=headers)
    
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
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/business-locations",
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
      .url("https://einvoice-api.sepay.vn/v1/business-locations")
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
    
    url = URI("https://einvoice-api.sepay.vn/v1/business-locations")
    
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
    
    	url := "https://einvoice-api.sepay.vn/v1/business-locations"
    
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
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/business-locations"),
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
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/business-locations")! as URL,
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
      .url("https://einvoice-api.sepay.vn/v1/business-locations")
      .get()
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## API tạo địa điểm

Thêm một địa điểm mới vào danh mục. Địa điểm mới ở trạng thái đang bật (`is_active=true`).

<Endpoint>
  <Method>POST</Method>

  <Path>https://einvoice-api.sepay.vn/v1/business-locations/create</Path>

  <Description>
    Tạo địa điểm kinh doanh
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Params>
  <RequestBody>
    <Fields>
      <Field name="code" type="string" required="true">
        Mã số địa điểm kinh doanh, tối đa 50 ký tự. Là định danh cố định - đặt khi tạo và không sửa được. Duy nhất trong công ty, trùng trả lỗi `SELLER_STORE_DUPLICATE_CODE`.
      </Field>
      <Field name="name" type="string" required="true">
        Tên địa điểm kinh doanh theo giấy chứng nhận đăng ký kinh doanh, tối đa 400 ký tự. Duy nhất trong công ty, trùng trả lỗi `SELLER_STORE_DUPLICATE_NAME`.
      </Field>
      <Field name="address" type="string" required="true">
        Địa chỉ địa điểm kinh doanh, tối đa 255 ký tự. Duy nhất trong công ty, trùng trả lỗi `SELLER_STORE_DUPLICATE_ADDRESS`.
      </Field>
    </Fields>

    <Example>
      {
        "code": "CN02",
        "name": "Chi nhánh 2",
        "address": "34 Đường Hoa Nắng, Quận 7, TP.HCM"
      }
    </Example>
  </RequestBody>
</Params>

<Responses>
  <Response status="201">
    <Description>
      Đã tạo địa điểm kinh doanh
    </Description>

    <Example>
      {
        "success": true,
        "data": {
          "xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40",
          "code": "CN02",
          "name": "Chi nhánh 2",
          "address": "34 Đường Hoa Nắng, Quận 7, TP.HCM",
          "is_used": false,
          "is_active": true
        }
      }
    </Example>
  </Response>

</Responses>

<ResponseDescriptionFields>
  <ResponseSchema status="201">
    <Fields>
      <Field name="success" type="boolean" required="false">
      </Field>
      <Field name="data" type="object" required="false">
        <Fields>
          <Field name="xid" type="string (uuid)" required="false">
            Mã định danh địa điểm. Dùng giá trị này làm `seller_store_xid` khi tạo hóa đơn.
          </Field>
          <Field name="code" type="string" required="false">
            Mã số địa điểm kinh doanh (tối đa 50 ký tự). Là định danh cố định, không sửa được sau khi tạo.
          </Field>
          <Field name="name" type="string" required="false">
            Tên địa điểm kinh doanh theo giấy chứng nhận đăng ký kinh doanh (tối đa 400 ký tự)
          </Field>
          <Field name="address" type="string" required="false">
            Địa chỉ địa điểm kinh doanh (tối đa 255 ký tự)
          </Field>
          <Field name="is_used" type="boolean" required="false">
            Cho biết địa điểm đã được dùng để xuất hóa đơn hay chưa. Địa điểm đã dùng (`true`) không xóa được, chỉ có thể hủy kích hoạt bằng `is_active`, nhưng vẫn sửa được tên và địa chỉ bình thường.

          </Field>
          <Field name="is_active" type="boolean" required="false">
            Trạng thái kích hoạt của địa điểm. Địa điểm bị hủy kích hoạt (`false`) không dùng để xuất hóa đơn mới được nữa, nhưng các hóa đơn đã gắn với nó trước đó vẫn tra cứu và lọc được bình thường.

          </Field>
        </Fields>
      </Field>
    </Fields>
  </ResponseSchema>

</ResponseDescriptionFields>

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
      --url https://einvoice-api.sepay.vn/v1/business-locations/create \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"code":"CN02","name":"Chi nhánh 2","address":"34 Đường Hoa Nắng, Quận 7, TP.HCM"}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/business-locations/create",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}",
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
    
    conn = http.client.HTTPSConnection("einvoice-api.sepay.vn")
    
    payload = "{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/business-locations/create", payload, headers)
    
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
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/business-locations/create",
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
      code: 'CN02',
      name: 'Chi nhánh 2',
      address: '34 Đường Hoa Nắng, Quận 7, TP.HCM'
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}");
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/create")
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
    
    url = URI("https://einvoice-api.sepay.vn/v1/business-locations/create")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}"
    
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
    
    	url := "https://einvoice-api.sepay.vn/v1/business-locations/create"
    
    	payload := strings.NewReader("{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}")
    
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
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/business-locations/create"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}")
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
      "code": "CN02",
      "name": "Chi nhánh 2",
      "address": "34 Đường Hoa Nắng, Quận 7, TP.HCM"
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/business-locations/create")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"code\":\"CN02\",\"name\":\"Chi nhánh 2\",\"address\":\"34 Đường Hoa Nắng, Quận 7, TP.HCM\"}")
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/create")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## API cập nhật địa điểm

Sửa tên, địa chỉ hoặc kích hoạt/hủy kích hoạt một địa điểm. Địa điểm đã dùng để xuất hóa đơn vẫn sửa được bình thường, hóa đơn đã xuất không bị ảnh hưởng.

<Endpoint>
  <Method>POST</Method>

  <Path>https://einvoice-api.sepay.vn/v1/business-locations/update</Path>

  <Description>
    Cập nhật địa điểm kinh doanh
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Params>
  <RequestBody>
    <Fields>
      <Field name="xid" type="string (uuid)" required="true">
        Mã định danh (`xid`, dạng UUID) của địa điểm cần cập nhật, lấy từ API danh sách địa điểm.
      </Field>
      <Field name="name" type="string" required="false">
        Tên mới, tối đa 400 ký tự, duy nhất trong công ty. Không bắt buộc - chỉ gửi khi muốn đổi tên, không gửi giữ nguyên tên cũ. Gửi nhưng để rỗng sẽ bị từ chối.
      </Field>
      <Field name="address" type="string" required="false">
        Địa chỉ mới, tối đa 255 ký tự, duy nhất trong công ty. Không bắt buộc - chỉ gửi khi muốn đổi địa chỉ, không gửi giữ nguyên địa chỉ cũ. Gửi nhưng để rỗng sẽ bị từ chối.
      </Field>
      <Field name="is_active" type="boolean" required="false">
        Kích hoạt (`true`) hoặc hủy kích hoạt (`false`) địa điểm. Địa điểm bị hủy kích hoạt không dùng để xuất hóa đơn mới nhưng hóa đơn cũ đã gắn vẫn tra cứu được.
      </Field>
    </Fields>

    <Example>
      {
        "xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40",
        "name": "Chi nhánh 2 - Quận 3",
        "address": "58 Đường Trăng Thanh, Quận 3, TP.HCM"
      }
    </Example>
  </RequestBody>
</Params>

<Responses>
  <Response status="200">
    <Description>
      Đã cập nhật địa điểm kinh doanh
    </Description>

    <Example>
      {
        "success": true,
        "data": {
          "xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40",
          "code": "CN02",
          "name": "Chi nhánh 2 - Quận 3",
          "address": "58 Đường Trăng Thanh, Quận 3, TP.HCM",
          "is_used": false,
          "is_active": true
        }
      }
    </Example>
  </Response>

</Responses>

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
      --url https://einvoice-api.sepay.vn/v1/business-locations/update \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"xid":"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40","name":"Chi nhánh 2 - Quận 3","address":"58 Đường Trăng Thanh, Quận 3, TP.HCM","is_active":false}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/business-locations/update",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}",
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
    
    conn = http.client.HTTPSConnection("einvoice-api.sepay.vn")
    
    payload = "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/business-locations/update", payload, headers)
    
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
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/business-locations/update",
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
      xid: 'b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40',
      name: 'Chi nhánh 2 - Quận 3',
      address: '58 Đường Trăng Thanh, Quận 3, TP.HCM',
      is_active: false
    }));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}");
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/update")
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
    
    url = URI("https://einvoice-api.sepay.vn/v1/business-locations/update")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}"
    
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
    
    	url := "https://einvoice-api.sepay.vn/v1/business-locations/update"
    
    	payload := strings.NewReader("{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}")
    
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
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/business-locations/update"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}")
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
      "xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40",
      "name": "Chi nhánh 2 - Quận 3",
      "address": "58 Đường Trăng Thanh, Quận 3, TP.HCM",
      "is_active": false
    ] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/business-locations/update")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\",\"name\":\"Chi nhánh 2 - Quận 3\",\"address\":\"58 Đường Trăng Thanh, Quận 3, TP.HCM\",\"is_active\":false}")
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/update")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## API xóa địa điểm

Xóa một địa điểm khỏi danh mục. Chỉ xóa được địa điểm **chưa từng dùng** để xuất hóa đơn.

<Endpoint>
  <Method>POST</Method>

  <Path>https://einvoice-api.sepay.vn/v1/business-locations/delete</Path>

  <Description>
    Xóa địa điểm kinh doanh
  </Description>

  <Authentication>
    bearerAuth
  </Authentication>
</Endpoint>

<Params>
  <RequestBody>
    <Fields>
      <Field name="xid" type="string (uuid)" required="true">
        Mã định danh (`xid`, dạng UUID) của địa điểm cần xóa, lấy từ API danh sách địa điểm. Chỉ xóa được địa điểm chưa dùng để xuất hóa đơn - địa điểm đã dùng hãy hủy kích hoạt (`is_active=false`) thay vì xóa.
      </Field>
    </Fields>

    <Example>
      {
        "xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40"
      }
    </Example>
  </RequestBody>
</Params>

<Responses>
  <Response status="200">
    <Description>
      Đã xóa địa điểm kinh doanh
    </Description>

    <Example>
      {
        "success": true,
        "message": "Đã xóa địa điểm kinh doanh người bán"
      }
    </Example>
  </Response>

</Responses>

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
      --url https://einvoice-api.sepay.vn/v1/business-locations/delete \
      --header 'Authorization: Bearer REPLACE_BEARER_TOKEN' \
      --header 'content-type: application/json' \
      --data '{"xid":"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40"}'
    ```
  </CodeSample>

  <CodeSample value="php_curl" lang="php">
    ```php
    <?php
    
    $curl = curl_init();
    
    curl_setopt_array($curl, [
      CURLOPT_URL => "https://einvoice-api.sepay.vn/v1/business-locations/delete",
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}",
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
    
    conn = http.client.HTTPSConnection("einvoice-api.sepay.vn")
    
    payload = "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}"
    
    headers = {
        'Authorization': "Bearer REPLACE_BEARER_TOKEN",
        'content-type': "application/json"
        }
    
    conn.request("POST", "/v1/business-locations/delete", payload, headers)
    
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
      "hostname": "einvoice-api.sepay.vn",
      "port": null,
      "path": "/v1/business-locations/delete",
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
    
    req.write(JSON.stringify({xid: 'b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40'}));
    req.end();
    ```
  </CodeSample>

  <CodeSample value="java_okhttp" lang="java">
    ```java
    OkHttpClient client = new OkHttpClient();
    
    MediaType mediaType = MediaType.parse("application/json");
    RequestBody body = RequestBody.create(mediaType, "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}");
    Request request = new Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/delete")
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
    
    url = URI("https://einvoice-api.sepay.vn/v1/business-locations/delete")
    
    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
    
    request = Net::HTTP::Post.new(url)
    request["Authorization"] = 'Bearer REPLACE_BEARER_TOKEN'
    request["content-type"] = 'application/json'
    request.body = "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}"
    
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
    
    	url := "https://einvoice-api.sepay.vn/v1/business-locations/delete"
    
    	payload := strings.NewReader("{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}")
    
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
        RequestUri = new Uri("https://einvoice-api.sepay.vn/v1/business-locations/delete"),
        Headers =
        {
            { "Authorization", "Bearer REPLACE_BEARER_TOKEN" },
        },
        Content = new StringContent("{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}")
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
    let parameters = ["xid": "b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40"] as [String : Any]
    
    let postData = JSONSerialization.data(withJSONObject: parameters, options: [])
    
    let request = NSMutableURLRequest(url: NSURL(string: "https://einvoice-api.sepay.vn/v1/business-locations/delete")! as URL,
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
    val body = RequestBody.create(mediaType, "{\"xid\":\"b6d767d2-f8ed-4c1b-9a2e-7f5c3d1a8b40\"}")
    val request = Request.Builder()
      .url("https://einvoice-api.sepay.vn/v1/business-locations/delete")
      .post(body)
      .addHeader("Authorization", "Bearer REPLACE_BEARER_TOKEN")
      .addHeader("content-type", "application/json")
      .build()
    
    val response = client.newCall(request).execute()
    ```
  </CodeSample>

</CodeSamples>

## Xử lý lỗi

Các endpoint địa điểm trả lỗi theo cấu trúc `{ "success": false, "error": { "code": "...", "message": "..." } }`.

<ErrorCodes
  hiddenHead={true}
  rows={[
  { code: 401, name: "UNAUTHORIZED", description: "Thiếu hoặc sai Bearer token." },
  { code: 422, name: "VALIDATION_ERROR", description: "Thiếu trường bắt buộc, vượt quá độ dài, hoặc không có thông tin nào để cập nhật." },
  { code: 422, name: "SELLER_STORE_DUPLICATE_CODE", description: "Mã số địa điểm đã tồn tại trong công ty (chỉ khi tạo)." },
  { code: 422, name: "SELLER_STORE_DUPLICATE_NAME", description: "Tên địa điểm đã tồn tại trong công ty." },
  { code: 422, name: "SELLER_STORE_DUPLICATE_ADDRESS", description: "Địa chỉ đã tồn tại trong danh mục công ty." },
  { code: 404, name: "SELLER_STORE_NOT_FOUND", description: "Địa điểm không tồn tại hoặc không thuộc tài khoản của bạn (khi cập nhật/xóa)." },
  { code: 409, name: "SELLER_STORE_IN_USE", description: "Địa điểm đã được dùng để xuất hóa đơn nên không thể xóa - hãy hủy kích hoạt (is_active=false) thay vì xóa." },
  { code: 405, name: "METHOD_NOT_ALLOWED", description: "Gọi sai HTTP method của endpoint." },
  { code: 500, name: "INTERNAL_ERROR", description: "Đã có lỗi xảy ra, vui lòng liên hệ SePay để được hỗ trợ." }
]}
/>

<Callout type="info" title="Về is_used và is_active">
`is_used=true`
: địa điểm đã được dùng để xuất ít nhất một hóa đơn. Địa điểm này không xóa được nhưng vẫn sửa tên, địa chỉ và kích hoạt/hủy kích hoạt (
`is_active`
) bình thường.
`is_active=false`
: địa điểm đã bị hủy kích hoạt, không chọn được khi xuất hóa đơn mới (
API xuất hóa đơn
 trả 
`400 SELLER_STORE_INACTIVE`
), nhưng các hóa đơn đã gắn trước đó vẫn lọc và tra cứu được.
</Callout>

## Bước tiếp theo

Sau khi có `xid` của địa điểm, bạn có thể:

1. **[Xuất hóa đơn điện tử](/vi/einvoice-api/v1/xuat-hoa-don-dien-tu)** - Gắn `seller_store_xid` để truyền mã số và tên địa điểm lên hóa đơn.
2. **[Danh sách hóa đơn](/vi/einvoice-api/v1/danh-sach-hoa-don)** - Lọc hóa đơn theo `seller_store_xid`.