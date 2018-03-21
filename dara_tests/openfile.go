package main

import (
    "os"
    "log"
    "fmt"
)

func main() {
    f, err := os.Open("hello_world.txt")
    if err != nil {
        log.Fatal("Failed to open file")
    }
    b := make([]byte, 12)
    n, _ := f.Read(b)
    fmt.Println(string(b[:n]))
    b2 := make([]byte, 10)
    n, _ = f.ReadAt(b2, 6)
    fmt.Println(string(b2[:n]))
    f.Close()
}
