package main

import (
    "fmt"
    "net/http"
    "strings"
    "math/rand"
    "log"
)

const URLComp = ":/?-_.~&*@0123456789"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ServePage(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if len(r.Form) == 0 {
        path := r.URL.Path
        if path != "/" {
            HandleShortened(path, w, r)
        } else {
            CreateLinkPage(w, r)
        }
    } else {
        CreateLink(w, r)
    }
}

func NotFound(message string, w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Couldn't find what you're looking for ): - " + message)
}

func CreateLinkPage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, FormOutput())
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
    url := r.Form["url"][0]
    key := r.Form["key"][0]

    if url == "" || key == "" {
        NotFound("Not enough information supplied.", w, r)
        return
    }

    acceptable := len(key) == 14 && IsAlpha(key)

    if !acceptable {
        NotFound("Incorrect key.", w, r)
        return
    }

    if !IsAcceptableURL(url) {
        NotFound("URL is not acceptable.", w, r)
        return
    }

    if !IsKeyAccepted(key) {
        NotFound("Key is incorrect.", w, r)
        return
    }

    link, exists := URLExists(url)
    if exists {
        message := fmt.Sprintf("URL already created at %s/%s.", r.Host, link)
        NotFound(message, w, r)
        return
    }

    shorten := AddLink(key, url)

    output := fmt.Sprintf("Created your link at %s/%s.", r.Host, shorten)
    fmt.Fprintf(w, output)
}

func AddLink(key string, url string) string {
    shorten := GenRandomShorten(4)
    query := fmt.Sprintf("INSERT INTO links VALUES (\"%s\", \"%s\", \"%s\", 0)", key, url, shorten)
    Execute(query)
    return shorten
}

func GenRandomShorten(n int) string {
    b := ""
    for i := 0; i < n; i++ {
        b = b + string(letters[rand.Intn(len(letters))])
    }
    return b
}

func URLExists(urls string) (string, bool) {
    query := "SELECT shorten, url FROM links"
    result, err := Database().Query(query)
    if err != nil {
        log.Fatal(err)
        return "", false
    }
    for result.Next() {
        var url string
        var shorten string
        err = result.Scan(&url, &shorten)
        if err != nil {
            return "", false
        }
        if url == urls {
            return shorten, true
        }
    }
    return "", false
}

func IsKeyAccepted(key string) bool {
    accepted := false
    query := fmt.Sprintf("SELECT key FROM keys WHERE key=\"%s\"", key)
    result, err := Database().Query(query)
    if err != nil {
        log.Fatal(err)
        return false
    }
    for result.Next() {
        accepted = true
    }
    return accepted;
}

func IsAcceptableURL(url string) bool {
    for _, r := range url {
        contains := false
        for _, w := range URLComp {
            if w == r {
                contains = true
                break
            }
        }
        for _, w := range letters {
            if w == r {
                contains = true
                break
            }
        }
        if !contains {
            return contains
        }
    }
    return true
}

func HandleShortened(path string, w http.ResponseWriter, r *http.Request) {
    path = strings.Replace(path, "/", "", 1)

    acceptable := len(path) == 4 && IsAlpha(path)

    if acceptable {
        query := fmt.Sprintf("SELECT shorten, url FROM links") //WHERE doesn't seem to work properly - work around for now
        result, err := Database().Query(query)
        if err != nil {
            NotFound("Database error.", w, r)
            return
        }
        var url string
        var shorten string
        var correctURL string
        for result.Next() {
            err = result.Scan(&url, &shorten)
            if err != nil {
                NotFound("AHHHH. DB error.", w, r)
                return
            }
            if shorten == path {
                correctURL = url
            }
        }

        if correctURL == "" {
            NotFound("URL not found.", w, r)
            return
        }

        http.Redirect(w, r, correctURL, 301)
    } else {
        NotFound("Incorrect URL. Must be 4 characters.", w, r)
    }
}

func IsAlpha(path string) bool {
    for _, p := range path {
        if (p < 'a' || p > 'z') && (p < 'A' || p > 'Z')  {
            return false
        }
    }
    return true
}
