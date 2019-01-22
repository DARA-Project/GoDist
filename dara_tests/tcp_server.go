package main

import (
    "net"
    "log"
    "fmt"
)

func main() {
    ln, err := net.Listen("tcp", "127.0.0.1:9000")
    if err != nil {
        log.Fatal(err)
    }

    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal(err)
        }

        var bs = make([]byte, 1024)
        n, err := conn.Read(bs)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Bytes read: ", n)
    }

}
