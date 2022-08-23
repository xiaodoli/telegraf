package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

type AuthToken struct {
    AccessToken string `json:"access_token"`
    Scope       string `json:"scope"`
    ExpiresIn   int32  `json:"expires_in"`
    TokenType   string `json:"token_type"`
}

func main() {
    method := "POST"
    url := os.Getenv("oauth_url")
    client_id := os.Getenv("client_id")
    client_secret := os.Getenv("client_secret")
    oauth_audience := os.Getenv("oauth_audience")
    oauth_grant_type := os.Getenv("oauth_grant_type")    
    oauth_content_type := os.Getenv("oauth_content_type")
    output_file := "/tmp/telegraf/access_token"
    if len(client_id) == 0 {
        log.Printf("invalid client_id, %d\n", len(client_id))
        return
    }

    if len(client_secret) == 0 {
        log.Println("invalid client_secret")
        return
    }

    if len(output_file) == 0 {
        log.Println("invalid output file")
        return
    }

    if (len(oauth_content_type) > 0 && strings.EqualFold(oauth_content_type, "application/json")) {
        payload := strings.NewReader(`{
            "client_id":"` + client_id + `",
        "client_secret":"` + client_secret + `",
        "audience":"` + oauth_audience + `",
        "grant_type":"` + oauth_grant_type + `"
        }`)
    } else {
        payload := strings.NewReader("client_id=" + client_id + "&client_secret=" + client_secret + "&grant_type= " + oauth_grant_type + "&audiance=" + oauth_audience + "")
    }

    client := &http.Client{}
    req, err := http.NewRequest(method, url, payload)

    if err != nil {
        fmt.Println(err)
        return
    }
    
    req.Header.Add("Content-Type", oauth_content_type)

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(body))

    var authToken AuthToken
    json.Unmarshal([]byte(string(body)), &authToken)
    // fmt.Printf("access_token: %s, expires_in: %d, scope: %s\n", authToken.AccessToken, authToken.ExpiresIn, authToken.Scope)

    f, err := os.Create(output_file)
    if err != nil {
        log.Fatal(err)
    }
    f.WriteString(authToken.AccessToken)
}
