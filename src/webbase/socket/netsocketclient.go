package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	//if len(os.Args) != 2 {
	//	fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
	//	os.Exit(1)
	//}
	params :=[]string{"127.0.0.1:8281"};
	service := params[0]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	CheckError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	CheckError(err)
	result, err := ioutil.ReadAll(conn)
	fmt.Println(string(result[:]))
	CheckError(err)
	fmt.Println(string(result))
	fmt.Println("tst")
	os.Exit(0)
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}