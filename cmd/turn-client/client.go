package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ernado/stun"
	"github.com/ernado/turn"
	"bytes"
	"encoding/base64"
)

var credentials = stun.NewLongTermIntegrity(
	"ernado",
	"a1.cydev.ru",
	"turn",
)

func main() {
	m := new(stun.Message)
	m.NewTransactionID()
	var (
		xma  stun.XORMappedAddress
		s    stun.Software
		code stun.ErrorCodeAttribute
	)
	raddr, err := net.ResolveUDPAddr("udp", "a1.cydev.ru:3479")
	if err != nil {
		log.Fatal(err)
	}
	laddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	// Type is allocate request.
	m.Type = stun.MessageType{
		Class:  stun.ClassRequest,
		Method: stun.MethodAllocate,
	}
	//turn.RequestedTransport{Protocol: turn.ProtoUDP}.AddTo(m)
	m.Add(stun.AttrRequestedTransport, []byte{17, 0, 0, 0})
	stun.Username("ernado").AddTo(m)
	stun.Realm("a1.cydev.ru").AddTo(m)
	m.WriteHeader()
	fmt.Println(credentials)
	if err := credentials.AddTo(m); err != nil {
		log.Fatal(err)
	}
	m.WriteHeader()
	decoded := new(stun.Message)
	decoded.Raw = make([]byte, 1024)
	if _, err := decoded.ReadFrom(bytes.NewReader(m.Raw)); err != nil {
		log.Fatal(err)
	}
	if err := credentials.Check(decoded); err != nil {
		log.Fatal(err)
	}
	fmt.Println("sending", m)
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Fatal(err)
	}
	//if _, err := m.WriteTo(conn); err != nil {
	//	log.Fatal(err)
	//}
	b, err := base64.StdEncoding.DecodeString("AAMACCESpEJYN2tEQUxIcGpWSUcAGQAEEQAAAA==")
	if err != nil {
		log.Fatal(err)
	}
	m.Reset()
	if _, err := m.Write(b); err != nil {
		log.Fatal(err)
	}
	v, err := m.Get(stun.AttrRequestedTransport)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
	fmt.Println("proto", v)
	credentials.AddTo(m)
	m.WriteHeader()
	if _, err := conn.Write(m.Raw); err != nil {
		log.Fatal(err)
	}
	conn.SetDeadline(time.Now().Add(time.Second * 2))
	got := new(stun.Message)
	got.Raw = make([]byte, 1024)
	if _, err := got.ReadFrom(conn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("got", got)
	s.GetFrom(got)
	fmt.Println("SOFTWARE:", string(s.Raw))
	for _, a := range got.Attributes {
		fmt.Println(a)
	}
	if got.Type.Class == stun.ClassErrorResponse {
		code.GetFrom(got)
		log.Println("error:", code)
		return
	}
	if err := xma.GetFrom(got); err != nil {
		log.Fatal(err)
	}
	fmt.Println("XOR RELAYED ADDR:", xma)
	m.Reset()
	m.NewTransactionID()
	m.Type = stun.MessageType{
		Class:  stun.ClassRequest,
		Method: stun.MethodRefresh,
	}
	if err := (turn.Lifetime{}).AddTo(m); err != nil {
		log.Fatal(err)
	}
	m.WriteHeader()
	if err := credentials.AddTo(m); err != nil {
		log.Fatal(err)
	}
	m.WriteLength()
	if _, err := m.WriteTo(conn); err != nil {
		log.Fatal(err)
	}
	got.Reset()
	conn.SetDeadline(time.Now().Add(time.Second * 2))
	if _, err := got.ReadFrom(conn); err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
	for _, a := range got.Attributes {
		fmt.Println(a)
	}
}
