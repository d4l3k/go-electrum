package irc

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/thoj/go-ircevent"
)

var (
	UsernameTemplate = "go-electrum-%d"
	IRCServer        = "chat.freenode.net:8000"
	IRCChannel       = "#electrum"
)

// FindElectrumServersIRC finds nodes to connect to by connecting to the Freenode #electrum channel.
func FindElectrumServers() ([]string, error) {
	user := fmt.Sprintf(UsernameTemplate, rand.Int31())
	ircobj := irc.IRC(user, user)
	ircobj.UseTLS = false
	ircobj.VerboseCallbackHandler = true
	ircobj.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	if err := ircobj.Connect(IRCServer); err != nil {
		return nil, err
	}
	defer ircobj.Quit()
	go ircobj.Loop()
	ircobj.Join(IRCChannel)
	ircobj.SendRawf("NAMES %s", IRCChannel)
	var wg sync.WaitGroup
	wg.Add(1)
	// NAMES END
	ircobj.AddCallback("366", func(event *irc.Event) {
		// THIS IS RACY UGH. This irc library doesn't preserve order of callbacks.
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	})

	// NAMES
	ircobj.AddCallback("353", func(event *irc.Event) {
		if len(event.Arguments) < 4 {
			return
		}
		words := strings.Split(event.Arguments[3], " ")
		for _, nick := range words {
			if strings.HasPrefix(nick, "E_") {
				wg.Add(1)
				ircobj.Whois(nick)
			}
		}
	})

	var serversLock sync.Mutex
	var servers []string

	// WHOIS realname
	ircobj.AddCallback("311", func(event *irc.Event) {
		detailsLine := strings.TrimSpace(event.Arguments[len(event.Arguments)-1])
		serversLock.Lock()
		servers = append(servers, detailsLine)
		serversLock.Unlock()

		wg.Done()
	})
	wg.Wait()
	return servers, nil
}
