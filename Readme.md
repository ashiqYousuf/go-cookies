# Cookies in Go

##### This repo will show you how to use cookies in your Go web application to persist data between HTTP requests for a specific client.

##### Cookies in Go are represented by the http.Cookie struct type.

```go
type Cookie struct {
    Name  string // Name of the cookies
    Value string // data that you want to persist

    // The Domain and Path attributes define the scope of the cookie
    Path       string    
    Domain     string    
    Expires    time.Time 
    RawExpires string   
    MaxAge   int 
    // Secure enables encrypted transmission & so on...
    Secure   bool
    HttpOnly bool
    SameSite SameSite
    Raw      string
    Unparsed []string
}

```

###### Cookies can be written in a HTTP Response using:
```http.SetCookie()```
###### And cookies can be read from http Request using:
```*Request.Cookie()```

We should not trust cookie data as they are stored on the client & its straightforward for a user to manipulate them. We must first verify the cookie hasn't been manipulated, for that we can generate a HMAC signature of the cookie name and value, and then prepend this signature to the cookie value before sending it to the client as:-
```cookie.Value = "{HMAC signature}{original value}"```


When we receive the cookie back from the client, we can recalculate the HMAC signature from the cookie name and original value, and check that the recalculated HMAC signature matches the signature at the start of the received cookie. If they match, it confirms the integrity of the cookie.

But if you do want to prevent the client from being able to read the cookie data, then we need to encrypt the data before writing it.


Storing complicated data types in cookies, then use encoding/gob package.
