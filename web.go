package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

var pageTemplate *template.Template

func StartTinyURLServer(cfg *Config, db *DB) error {
	var err error
	pageTemplate, err = template.New("page").Parse(pageHTML)
	if err != nil {
		Warnf("Create page html/template is failed. Error: %v\n", err)
	}
	s := CreateTinyURLServer(cfg, db)

	// now, http only
	return http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), s)
}

func CreateTinyURLServer(cfg *Config, db *DB) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/page", pageHandleMiddle(cfg, db))
	server.HandleFunc("/", tinyURLHandleMiddle(cfg, db))
	return server
}

func pageHandleMiddle(cfg *Config, db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if pageTemplate == nil {
				msg := fmt.Sprintf("Page was not found.\n")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(msg))
				return
			}
			if err := pageTemplate.Execute(w, nil); err != nil {
				Errorf("Executing page template is failed.\n")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			Debugf("Request not allowed method '%s'\n", r.Method)
			msg := fmt.Sprintf("HTTP method '%s' is not allowed.\n", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(msg))
		}
	}
}

func tinyURLHandleMiddle(cfg *Config, db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getTinyURL(cfg, db, w, r)
		case "POST":
			postTinyURL(cfg, db, w, r)
		default:
			Debugf("Request not allowed method '%s'\n", r.Method)
			msg := fmt.Sprintf("HTTP method '%s' is not allowed.\n", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(msg))
		}
	}
}

func getTinyURL(cfg *Config, db *DB, w http.ResponseWriter, r *http.Request) {
	Debugf("Request redirect of tiny '%s' from %s\n", r.URL.Path, r.RemoteAddr)
	origin, err := db.GetOriginURL(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("'%s' is not found.\n", r.RequestURI)))
		return
	}
	w.Header().Set("Location", origin)
	w.WriteHeader(http.StatusMovedPermanently)
}

type TinyPost struct {
	Origin string `json:"Origin"`
	Tiny   string `json:"Tiny"`
	Error  string `json:"Error"`
}

func postTinyURL(cfg *Config, db *DB, w http.ResponseWriter, r *http.Request) {
	Debugf("New URL is posted from %s\n", r.RemoteAddr)

	body := make([]byte, int(r.ContentLength))
	_, err := r.Body.Read(body)
	if err != io.EOF {
		Warnf("Request's body might have been read all.\n")
	}

	data := TinyPost{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rBody, _ := json.Marshal(TinyPost{Error: "Internal server error.\n"})
		w.Write(rBody)
		Errorf("JsonUnmarshalError: %v\n", err)
		return
	}

	c := http.Client{Timeout: time.Second * 10}
	if resp, err := c.Get(data.Origin); err != nil || resp.StatusCode >= 300 || resp.StatusCode < 200 {
		if err != nil {
			Warnf("HEAD request for origin is failed. Error: %v\n", err)
		} else {
			Infof("Unexpected status code '%d' is returned by HEAD request for origin\n", resp.StatusCode)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		rBody, _ := json.Marshal(TinyPost{Error: "Content of requested URL is invalid.\n"})
		w.Write(rBody)
		return
	}

	tiny, err := db.AddTinyURL(data.Origin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rBody, _ := json.Marshal(TinyPost{Error: "Internal server error.\n"})
		w.Write(rBody)
		Errorf("AddTinyURLError: %v\n", err)
		return
	}

	rBody, _ := json.Marshal(TinyPost{
		Origin: data.Origin,
		Tiny:   cfg.Protocol + "://" + r.Host + "/" + tiny,
	})
	w.Write(rBody)
}

const pageHTML string = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>tiny-url</title>
    <meta name="viewport" content="width=device-width,initial-scale=1.0,minimum-scale=1.0" />
  <body>
    <div class="title">tiny-url</div>
    <div class="form">
      <input id="url" type="text" placeholder="input text you want to shorten and enter!">
      <div class="result">
        tiny -> <span>http://yahoo.com/hogehoge</span>
      </div>
    </div>
    <script>
      document.getElementById("url").addEventListener("change",()=> {
        document.querySelector(".result").style.display = "none";
        if(document.getElementById("url").value.match(/(http|https):\/\//) === null) {
          alert("Invalid format for URL.");
          return;
        }
        let xhr = new XMLHttpRequest();
        xhr.open("POST", location.href+"api"); 
        xhr.onload = () => {
		  console.log("HTTP status code: " + xhr.statusText)
	      if (xhr.status.toString().match(/2[0-9]{2}/) === null) {
            let error = JSON.parse(xhr.responseText).Error;
			if(error) alert(error); return;
		  	alert("Request is failed.");
			return;
		  }
          let tiny = JSON.parse(xhr.responseText).Tiny;
          document.querySelector(".result span").innerText = tiny;
          document.querySelector(".result").style.display = "block";
        }
        xhr.send(JSON.stringify({Origin: document.getElementById("url").value}));
      });
    </script>
  </body>
  <style>
    body,div {
      margin:0px;
      padding:0px;
    }
    body {
      color: rgb(68, 67, 67);
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }
    .title {
      margin-top: 80px;
      margin-bottom: 15px;
      font-size: 40px;
      text-align: center;
    }
    .form {
      text-align: center;
      width: 80vw;
      margin: 0px auto;
    }
    .form input[type="text"] {
      width: 100%;
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      font-size:16px;
      border-width: 2px;
      border-style: solid;
      border-radius: 4px;
      padding: 4px 10px;
      vertical-align: center;
    }
    .form input[type="text"]::placeholder {
      color:rgb(163, 162, 162);
    }
    .result {
      display: none;
      padding: 0px 13px;
      margin-top: 10px;
      font-size:16px;
      text-align: left;
      animation: slidein 1s;
    }
    .result span {
      user-select: all;
    }

    @media (min-width: 640px) {
      .title{
        margin-top: 100px;
        margin-bottom: 20px;
        font-size: 56px;
      }
      .form {
        width:70%;
      }
      .form input[type="text"] {
        font-size:20px;
      }
      .result {
        font-size:20px;
      }
    }
    @media (min-width: 960px) {
      .title{
        margin-top: 150px;
        margin-bottom: 25px;
        font-size: 64px;
      }
      .form {
        width:50%;
      }
      .form input[type="text"] {
        font-size:24px;
      }
      .result {
        font-size:24px;
      }
    }

    @keyframes slidein {
      0% {
        opacity: 0;
        transform: translateX(-64px);
      }
      100% {
        opacity: 1;
        transform: translateX(0px);
      }
    }
    </style>
</html>
`
