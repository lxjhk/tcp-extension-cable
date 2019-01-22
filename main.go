package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	listeningAddrPrefix := "0.0.0.0:"

	springboardMode := false

	modePtr := flag.String("mode", "", "value can be 'springboard' or 'local-client'")
	listeningAddrPtr := flag.String("lp", "", "specify the listening port on the server")
	destAddrPtr := flag.String("da", "", "dest address with address in the form of x.x.x.x:x")

	flag.Parse()

	if *modePtr == "springboard" {
		springboardMode = true
		log.Println("springboard mode selected")
	} else if *modePtr == "local-client" {
		springboardMode = false
		log.Println("local-client mode selected")
	} else if *modePtr == "" {
		flag.PrintDefaults()
		os.Exit(0)
	} else {
		log.Fatalln("Mode is not set correctly")
	}

	listeningAddr, err := net.ResolveTCPAddr("tcp", listeningAddrPrefix+*listeningAddrPtr)
	log.Println("listening for client connection on " + listeningAddrPrefix + *listeningAddrPtr)
	checkError(err)

	log.Println("destination addr is " + *destAddrPtr)

	listener, err := net.ListenTCP("tcp", listeningAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		log.Println("New connection from " + conn.RemoteAddr().String())

		go func() {
			r := 5

			var connRemote net.Conn

			for err != nil && r > 0 || connRemote == nil {
				connRemote, err = net.DialTimeout("tcp", *destAddrPtr, 3*time.Second)
				r--
			}

			if err != nil {
				log.Fatal("Connection has been attempted 5 times. Failing..", err.Error())
				connRemote.Close()
				conn.Close()
			}

			if springboardMode {
				go ChanPipeEncrypt(conn, connRemote)
				go ChanPipeDecrypt(connRemote, conn)
			} else {
				go ChanPipeDecrypt(conn, connRemote)
				go ChanPipeEncrypt(connRemote, conn)
			}
		}()

	}

}

func ChanPipeEncrypt(conn1 net.Conn, conn2 net.Conn) {
	b := make([]byte, 1024)
	for {
		n, err := conn1.Read(b)
		if n > 0 {
			res := make([]byte, n)
			copy(res, b[:n])
			Encrypt(res)
			conn2.Write(res)
		}
		if err != nil {
			log.Println("Connection with " + conn1.RemoteAddr().String() + " closed")
			conn1.Close()
			log.Println("Closing Connection with " + conn2.RemoteAddr().String())
			conn2.Close()
			break
		}
	}
}

func ChanPipeDecrypt(conn1 net.Conn, conn2 net.Conn) {
	b := make([]byte, 1024)
	for {
		n, err := conn1.Read(b)
		if n > 0 {
			res := make([]byte, n)
			copy(res, b[:n])
			Decrypt(res)
			conn2.Write(res)
		}
		if err != nil {
			log.Println("Connection with " + conn1.RemoteAddr().String() + " closed")
			conn1.Close()
			log.Println("Closing Connection with " + conn2.RemoteAddr().String())
			conn2.Close()
			break
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func flipBits(bArray []byte) {
	for idx, element := range bArray {
		bArray[idx] = ^element
	}
}

func Encrypt(bArray []byte) {
	flipBits(bArray)
}

func Decrypt(bArray []byte) {
	flipBits(bArray)
}
