package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"net"
	"time"
	"tp.app/src/client"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	conn := startclient()

	g.SetManagerFunc(client.Layout)

	// Bind enter key to input to send new messages.
	err = g.SetKeybinding("jugador", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, iv *gocui.View) error {
		// Read buffer from the beginning.
		iv.Rewind()

		// Send message if text was entered.
		if len(iv.Buffer()) >= 2 {
			msg := client.AddingValue(iv.Buffer())
			conn.Write([]byte(msg))
			// Reset input.
			iv.Clear()

			// Reset cursor.
			err := iv.SetCursor(1, 1)
			if err != nil {
				log.Println("Failed to set cursor:", err)
			}
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	if err := client.InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}
}

func startclient() *net.TCPConn {
	time.Sleep(time.Second) // so that server has time to start
	servAddr := "127.0.0.1:8081"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	//conn.SetNoDelay(false)
	conn.SetWriteBuffer(10000)
	//msg := "abc\n"
	start := time.Now()
	//for i := 0; i < 1000000; i++ {
	//	conn.Write([]byte(msg))
	//	//bufio.NewReader(conn).ReadString('\n')
	//	//fmt.Print("Message from server: ", response)
	//}
	fmt.Println("took:", time.Since(start))
	return conn
}