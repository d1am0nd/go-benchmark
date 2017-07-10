package main

import (
    "os"
    "fmt"
    "math/rand"
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
    server := Server{r: NewRouter()}

    fmt.Println("Listening on port " + conf.ServerPort)

    http.ListenAndServe(conf.ServerPort, &server)
}

/* Server stuff */
type Server struct {
    r *httprouter.Router
}

// Global setting of headers
func (s *Server) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // http://localhost:8080
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    s.r.ServeHTTP(w, r)
}

// Routes
func NewRouter() *httprouter.Router {
    r := httprouter.New()
    r.GET("/", Home)
    r.GET("/api/first_db_result", FirstDbResult)
    r.GET("/api/all_db_results", AllDbResults)
    r.GET("/api/string_result", StringResult)
    r.GET("/api/qsort_result", QSortResult)
    return r
}

func Home(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Write([]byte("Hello"))
}

/* Controllers */
// Returns json of first db row
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
// Returns json of all db rows
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
// Returns hand written string
func StringResult(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    w.Write([]byte("Lorem ipsum"))
}

func QSortResult(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    sorted := qsort(GetNumbersFromFile())
    rjson, err := json.Marshal(sorted)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(rjson)
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

// Puts SQL DSN together
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

// Inits config from file
func GetConf() Conf {
    file, err := os.Open("./config.json")
    if err != nil {
        panic(err)
    }
    defer file.Close()

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

// Returns first row in `tests` table
func FindFirstTest() (TestModel, error) {
    tm := TestModel{}
    err := SQL.Get(&tm, "SELECT * FROM tests LIMIT 1")
    return tm, err
}

// Returns all rows in `tests` table
func GetAllTests() ([]TestModel, error) {
    tms := []TestModel{}
    err := SQL.Select(&tms, "SELECT * FROM tests")
    return tms, err
}


/** Quicksort stuff */
func GetNumbersFromFile() []int {
    var numbers []int

    file, err := os.Open("./numbers.json")
    if err != nil {
        return []int{}
    }
    defer file.Close()

    decoder := json.NewDecoder(file)

    err = decoder.Decode(&numbers)
    if err != nil {
        return []int{}
    }

    return numbers
}


// https://stackoverflow.com/questions/15802890/idiomatic-quicksort-in-go
func qsort(a []int) []int {
  if len(a) < 2 { return a }

  left, right := 0, len(a) - 1

  // Pick a pivot
  pivotIndex := rand.Int() % len(a)

  // Move the pivot to the right
  a[pivotIndex], a[right] = a[right], a[pivotIndex]

  // Pile elements smaller than the pivot on the left
  for i := range a {
    if a[i] < a[right] {
      a[i], a[left] = a[left], a[i]
      left++
    }
  }

  // Place the pivot after the last smaller element
  a[left], a[right] = a[right], a[left]

  // Go down the rabbit hole
  qsort(a[:left])
  qsort(a[left + 1:])


  return a
}
