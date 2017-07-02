package main

import (
    "os"
    "fmt"
    "encoding/json"
)

func main() {
    fmt.Println("Hello world, I'm alive")

    conf := GetConf()

    fmt.Println(conf)
}

/* Config stuff */
type Conf struct {
    ServerPort string `json: "server_port"`
    DbHostname string `json: "db_hostname"`
    DbUsername string `json: "db_username"`
    DbPassword string `json: "db_password"`
    DbPort int `json: "db_port"`
    DbName string `json: "db_name"`
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


