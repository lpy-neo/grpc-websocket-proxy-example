package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/lpy-neo/grpc-websocket-proxy/examples/cmd/wsechoserver/helloserver"
	log "github.com/sirupsen/logrus"
)

type Server2 struct{}

func (s *Server2) Stream(_ *helloserver.Empty, stream helloserver.HelloService_StreamServer) error {
	start := time.Now()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&helloserver.HelloResponse{
			Message: "hello there2!" + fmt.Sprint(time.Now().Sub(start)),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server2) Hello(srv helloserver.HelloService_HelloServer) error {
	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		if err := srv.Send(&helloserver.HelloResponse{
			Message: req.Message + "2!",
		}); err != nil {
			return err
		}
	}
}

func (s *Server2) Heartbeats(srv helloserver.HelloService_HeartbeatsServer) error {
	go func() {
		for {
			_, err := srv.Recv()
			if err != nil {
				log.Println("Recv() err:", err)
				return
			}
			log.Println("got hb from client")
		}
	}()
	t := time.NewTicker(time.Second * 1)
	for {
		log.Println("sending hb")
		hb := &helloserver.Heartbeat{
			Status: helloserver.Heartbeat_OK,
		}
		b := new(bytes.Buffer)
		if err := (&jsonpb.Marshaler{}).Marshal(b, hb); err != nil {
			log.Println("marshal err:", err)
		}
		log.Println(string(b.Bytes()))
		if err := srv.Send(hb); err != nil {
			return err
		}
		<-t.C
	}
	return nil
}
