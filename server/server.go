package main

import (
	"log"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/chat/messages"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/emirpasic/gods/sets/hashset"
)
const server_addr="127.0.0.1:8080"
func notifyAll(clients *hashset.Set, message interface{}) {
	log.Printf("send msg %v %d client\n",message,clients.Size())
	for _, tmp := range clients.Values() {
		client := tmp.(*actor.PID)
		client.Tell(message)
	}
}

func main() {
	remote.Start(server_addr)
	clients := hashset.New()
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *actor.Started:
			println("room is started!")
		case *actor.Stopping:
			println("room is stopping!")

		case *messages.Connect:
			log.Printf("Client [%v,%p] connected", msg.Sender,msg.Sender)
			clients.Add(msg.Sender)
			msg.Sender.Tell(&messages.Connected{Message: "Welcome!"})
			context.Watch(msg.Sender)
		case *actor.Terminated:
			log.Printf("remove before:clients[%+v],who[%p,%v],connected size is %v ", clients,msg.Who,msg.Who,clients.Size())
			clients.Remove(msg.Who)
			log.Printf("remove after :clients[%+v],who[%p],connected size is %v ", clients,msg.Who,clients.Size())

		case *messages.SayRequest:
			notifyAll(clients, &messages.SayResponse{
				UserName: msg.UserName,
				Message:  msg.Message,
			})
		case *messages.NickRequest:
			notifyAll(clients, &messages.NickResponse{
				OldUserName: msg.OldUserName,
				NewUserName: msg.NewUserName,
			})
		}
	})
	remote.Register("room",props)
	console.ReadLine()
}
