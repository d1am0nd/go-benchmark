package main

import (
    "os"
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/jmoiron/sqlx"
    "github.com/julienschmidt/httprouter"
    _ "github.com/go-sql-driver/mysql"
)

var SQL  *sqlx.DB

func main() {
    fmt.Println("Hello world, I'm alive")

    conf := GetConf()
    SQL = sqlx.MustConnect("mysql", conf.DSN())
    router := NewRouter()

    fmt.Println("Listening on port " + conf.ServerPort)

    http.ListenAndServe(conf.ServerPort, router)
}

/* Config stuff */
type Conf struct {
    ServerPort string `json:"server_port"`
    DbHostname string `json:"db_hostname"`
    DbUsername string `json:"db_username"`
    DbPassword string `json:"db_password"`
    DbPort int `json:"db_port"`
    DbName string `json:"db_name"`
}

func (c Conf) DSN() string {
    return c.DbUsername +
        ":" +
        c.DbPassword +
        "@tcp(" +
        ":" +
        fmt.Sprintf("%d", c.DbPort) +
        ")/" +
        c.DbName
}

func GetConf() Conf {
    file, err := os.Open("./config.json")
    if err != nil {
        panic(err)
    }

    decoder := json.NewDecoder(file)
    configuration := Conf{}

    err = decoder.Decode(&configuration)
    if err != nil {
      panic(err)
    }

    return configuration
}

/* Mysql stuff */
type TestModel struct {
    Id uint32 `json:"id" db:"id"`
    String string `json:"string" db:"string"`
    Content string `json:"content" db:"content"`
    CreatedAt string `json:"created_at" db:"created_at"`
    UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func FindFirstTest() (TestModel, error) {
    tm := TestModel{}
    err := SQL.Get(&tm, "SELECT * FROM tests LIMIT 1")
    return tm, err
}

func GetAllTests() ([]TestModel, error) {
    tms := []TestModel{}
    err := SQL.Select(&tms, "SELECT * FROM tests")
    return tms, err
}

/* Router stuff */
func NewRouter() *httprouter.Router {
    r := httprouter.New()
    r.GET("/", Home)
    r.GET("/api/first_db_result", FirstDbResult)
    r.GET("/api/all_db_results", AllDbResults)
    r.GET("/api/string_result", StringResult)
    return r
}

func Home(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Write([]byte("Hello"))
}

func FirstDbResult(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    res, err := FindFirstTest()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    rjson, err := json.Marshal(res)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(rjson)
}

func AllDbResults(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    res, err := GetAllTests()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    rjson, err := json.Marshal(res)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(rjson)
}

func StringResult(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Write([]byte("Hello"))
}
