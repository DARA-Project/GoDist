package main

import "net"
import "fmt"
import "os"
import "bufio"

func main() {
   conn, _ := net.Dial("tcp", "127.0.0.1:8081")
   for {
     // read in input from stdin
     reader := bufio.NewReader(os.Stdin)
     text, _ := reader.ReadString('\n')
     // send to socket
     fmt.Fprintf(conn, text)
     // listen for reply
     bufio.NewReader(conn).ReadString('\n')
   }
}
