package main

import (
    "net"
    "log"
    "fmt"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:9000")

    if err != nil {
        log.Fatal(err)
    }

    msg := "hello"

    m, err :=  conn.Write([]byte(msg))

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Wrote " , m)
}
