# Gophernoon - Building A Link Shortner Service

## Project Requirements

* Accept a full URL and store an alias
    ``` sh
    curl -X POST \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "url=https://example.com&alias=alias" \
        localhost:80/create
    ```
    * Ensure you're responding with a `201 Created` HTTP status code
    * Consider using a map to store `alias:url`
        * https://gobyexample.com/maps
* Validate the input to ensure the URL is valid. If validation fails ensure you
    are responding with a `400 Bad Request` HTTP status code.
    * https://gobyexample.com/url-parsing
* Expose a web server which you can access a `/{alias}` on
    * if an `alias:url` pair is found redirect using a `301 Moved Permanently` 
        HTTP status code
    * If an `alias:url` pair is not found respond with a `404 Not Found` HTTP 
        status code
    * https://gobyexample.com/http-servers
* When writing a HTTP Status Code you can use the Go `net/http` package 
    constants rather than using magic strings etc.
    * https://go.dev/src/net/http/status.go

## Optional Tasks

### Make Your Application Configurable

* Use flags and environment variables to change the HTTP port that is exposed
    * https://gobyexample.com/command-line-flags
    * https://gobyexample.com/environment-variables

### Containerise Your Application

* Use docker, or your favourite container platform, to make your application 
    more portable.
    * https://docs.docker.com/language/golang/build-images/

### Implement Alias Generation

* When a request comes through without an alias generate one automatically
    ``` sh
    curl -X POST \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "url=https://oid.skybet.net" \
        localhost:80
    ```

### Implement HTTP Testing

* Use the `httptest` library to write an integration test
    * https://pkg.go.dev/net/http/httptest

### Add PPROF for profiling

* Expose the PPROF frontend through HTTP and take a look at profiling your
    application
    * https://pkg.go.dev/net/http/pprof

### Implement Persistent Storage

* Ensure url:alias pairs are stored somewhere and can be loaded in on 
    application startup.
    * File Storage: https://gobyexample.com/writing-files
    * Redis Storage: https://github.com/redis/go-redis
    * Remember if you use a database or other application for storage you 
        should also containerise it. If you're using docker consider using
        `docker-compose` to handle multiple images
        * https://docs.docker.com/compose/

### Implement Content Management

* List all the `alias:url` pairs in the system
    ``` sh
    curl localhost:80
    ```
    * How should this be delivered? Plaintext? JSON? CSV? XML?
        * https://gobyexample.com/json
        * https://pkg.go.dev/encoding/csv
        * https://gobyexample.com/xml
    * Is the response sorted, if so in which order?
        * https://pkg.go.dev/sort
        * https://gobyexample.com/sorting
* Accept a `DELETE` request to remove an `alias:url` pair
    ``` sh
    curl -X DELETE \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "alias=oid" \
        localhost:80
    ```

### Implement Usage Metrics

* How many times has each `alias:url` pair been accessed?
    * Consider making a struct to store more information
    ``` go
    type Alias struct {
        Name string
        URL string
        Usage int
    }
    ```
* Consider using Prometheus, a common monitoring tool used in industry. 
    * https://prometheus.io/docs/guides/go-application/
    
### Implement a Front End

* Create a HTML page with a form to allow easier creation of alias'.
``` html
<form action="http://localhost:80" method="POST">
    <label for="url">URL</label>
    <input type="text" id="url" name="url">
    <label for="alias">Alias</label>
    <input type="text" id="alias" name="alias">
    <input type="submit" value="Submit">
</form>
```
* Will this be delivered through the Go application or via a web server like
    Nginx?
    * https://www.nginx.com/
* If delivering through your application consider using templates and the 
    embedded filesystem to improve extensibility and portability.
    * https://gobyexample.com/text-templates
    * https://gobyexample.com/embed-directive