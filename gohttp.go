package goshare

import (
  "fmt"
  "net/http"
  "html/template"
  "runtime"
  "time"

  "github.com/jmhodges/levigo"

  "github.com/abhishekkr/goshare/leveldb"
)

func HomePage(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/index.html")
  t.Execute(w, nil)
}

func GetReadKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  val := abkleveldb.GetValues(req.Form["key"][0], db)
  w.Write([]byte(val))
}

func GetPushKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  status := abkleveldb.PushKeyVal(req.Form["key"][0], req.Form["val"][0], db)
  if status != true {
    http.Error(w, "FATAL Error", http.StatusInternalServerError)
  }
  w.Write([]byte("Success"))
}

func GoShareHTTP(leveldb *levigo.DB, httpport int) {
  db = leveldb
  runtime.GOMAXPROCS(runtime.NumCPU())

  http.HandleFunc("/", HomePage)
  http.HandleFunc("/get", GetReadKey)
  http.HandleFunc("/put", GetPushKey)

  srv := &http.Server{
    Addr:        fmt.Sprintf(":%d", httpport),
    Handler:     http.DefaultServeMux,
    ReadTimeout: time.Duration(5) * time.Second,
  }

  fmt.Printf("access your goshare at http://<IP>:%d\n", httpport)
  srv.ListenAndServe()
}
