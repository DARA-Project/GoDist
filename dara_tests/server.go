package main

import "net"
import "fmt"
import "bufio"

func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  fmt.Println("Listening now...")
  // accept connection on port
  conn, _ := ln.Accept()
  fmt.Println("Acception connection")
  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    // output message received
    // sample process for string received
    // send new string back to client
    conn.Write([]byte(message + "\n"))
  }
}
