package main

import (
  "fmt"
  "net/http"
  "html/template"
	"io"
  "os"
  "path/filepath"
  "encoding/base64"
  "bytes"
  "encoding/json"
  "time"
)

type FormData struct {
  Loaded bool
  Model []string
  Result string
}

type RespponsePostData struct {
  Detail int64
  Session_id int64
}

type RespponseGetData struct {
  Detail int64     `json:"detail"`
  Session_id int64 `json:"session_id"`
  Result string    `json:"result"`
}

var api_url = "http://78.142.222.199:8000/api/speech_api/"
var Data = FormData{false, []string{"ru simple", "ru hard", "en simple", "en hard"}, "" }

// Display the named template
func display(w http.ResponseWriter, page string, is_load bool, result string) {
  t, err := template.ParseFiles("templates/index.html")

  if err != nil {
    fmt.Println(w, err.Error())
  }

  Data.Loaded = is_load
  Data.Result = result
  t.Execute(w, Data)
}

func post_data(url string, encoded_data string, ext string, model string) (status int64, session_id int64) {
  postBody, _ := json.Marshal(map[string]string{
    "encoded_data":encoded_data,
    "ext":ext,
    "model":model,
    "vocab":"",
  })

  responseBody := bytes.NewBuffer(postBody)

  resp, err :=  http.Post(url, "application/json", responseBody)
  if err != nil {
    fmt.Println("An Error occured %v", err)
  }
  defer resp.Body.Close()

  post_data := RespponsePostData{}
  json.NewDecoder(resp.Body).Decode(&post_data)
  return post_data.Detail, post_data.Session_id
}

func encode_data(file io.Reader) (encoded_data string) {
  buf := new(bytes.Buffer)
  buf.ReadFrom(file)
  encoded_data = base64.StdEncoding.EncodeToString([]byte(buf.String()))
  return
}

func save_file(w http.ResponseWriter, file io.Reader, filename string) (file_path string) {
  // Create file
  file_path = filepath.Join("uploads", filename)
	dst, err := os.Create(file_path)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  return
}

func get_data_from_session(session int64) (result_txt string) {
  url_session := fmt.Sprintf("%s%d", api_url, session)
  fmt.Println(url_session)
  for {
    time.Sleep(time.Second)
    resp, err := http.Get(url_session)
    if err != nil {
      fmt.Println(err)
      return
    }
    defer resp.Body.Close()

    get_data := RespponseGetData{}
    json.NewDecoder(resp.Body).Decode(&get_data)

    if get_data.Detail == 200{
      return get_data.Result
    }
  }
}

func upload_file(w http.ResponseWriter, r *http.Request) {
  // -------------------------------------------------------- //
  uploaded_file, handler_file, err := r.FormFile("upload_file")
  if err != nil {
      panic(err)
  }
  defer uploaded_file.Close()
  // -------------------------------------------------------- //

  //var file_path = save_file(w, uploaded_file, handler_file.Filename)
  var encoded_file = encode_data(uploaded_file)
  model := r.FormValue("language_model")

  fmt.Println("model = ", model, "file = ", handler_file.Filename)
  ext := filepath.Ext(filepath.Join("uploads", handler_file.Filename))

  status, session := post_data(api_url, encoded_file, ext, model)
  fmt.Printf("status = %d session = %d ", status, session)

  result_txt := get_data_from_session(session)
  fmt.Println(result_txt)

  display(w, "index", true, result_txt)
}

func view_handler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  	case "GET":
  		display(w, "index", false, "")
  	case "POST":
  		upload_file(w, r)
	}
}

func main() {
	http.HandleFunc("/", view_handler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8000", nil)
}
