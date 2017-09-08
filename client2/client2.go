package main

import (
	"log"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/chat/messages"
	"github.com/AsynkronIT/protoactor-go/remote"
	"time"
	"strconv"
	"fmt"
)

type Client struct {
	name string
}

func (c *Client) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.Connected:
		log.Println(msg.Message)
	case *messages.SayResponse:
		log.Printf("%v receive msg %s send info %v", c.name,msg.UserName, msg.Message)
	case *messages.NickResponse:
		log.Printf("%v is now known as %v", msg.OldUserName, msg.NewUserName)
	case *actor.Terminated:
		log.Printf("Note:%v send terminated:", msg.Who)

	}
}
func main() {
	remote.Start("127.0.0.1:12345")

	//server := actor.NewPID("127.0.0.1:8080", "chatserver")
	server, err := remote.SpawnNamed("127.0.0.1:8080", "chatserver", "room", time.Second*3)
	if err != nil {
		println(err)
		return
	}
	//spawn our chat client inline
	props1 := actor.FromInstance(&Client{"c"})
	props2 := actor.FromInstance(&Client{"d"})

	client1 := actor.Spawn(props1)
	client2 := actor.Spawn(props2)
	fmt.Println(client1.Id)
	fmt.Println(client2.Id)

	go foo("c", "message ", client2, server)

	server.Tell(&messages.Connect{
		Sender: client1,
	})
	nick := "d"
	cons := console.NewConsole(func(text string) {
		server.Request(&messages.SayRequest{
			UserName: nick,
			Message:  text,
		}, client1)
	})
	time.Sleep(time.Second*5)
	//write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		server.Request(&messages.NickRequest{
			OldUserName: nick,
			NewUserName: newNick,
		}, client1)
		nick = newNick
	})
	cons.Run()
}
func foo(senderName, text string, client *actor.PID, server *actor.PID) {
	server.Tell(&messages.Connect{
		Sender: client,
	})
	time.Sleep(time.Second*2)
	for i := 0; i < 5; i++ {
		server.Request(&messages.SayRequest{
			UserName: senderName,
			Message:  strconv.Itoa(i)+":"+text,
		}, client)
		time.Sleep(time.Second*2)
	}
	client.Stop()
}
