package blacklisthandler

import (
	"flag"
	"log"

	"github.com/devopsfaith/bloomfilter/rpc/client"
)

var server *string // server connection flag
var key *string    // key flag for jti


type blacklist interface {
	Add(id string)
}

type Conn struct {
	client *client.Bloomfilter
}

func newBlacklistHandlerConn(c *client.Bloomfilter) *Conn {
	return &Conn{
		client: c,
	}
}

func (cn *Conn) Add (id string)  {

	subject := *key + "-" + id
	cn.client.Add([]byte(subject))
}

// connect to the bloom filter server plugin in krakend
func Connect(uri string) (*Conn, error) {

	setFlags(uri)

	bloomFilterClient, err := client.New(*server)
	if err != nil {
		log.Println("unable to create the rpc client:", err.Error())
		return nil, err
	}
	log.Println("connected to bloom filter")

	conn := newBlacklistHandlerConn(bloomFilterClient)

	return conn, nil
}


func setFlags(uri string) {

	keyName := "jti"
	server = flag.String("server", uri, "ip:port of the remote bloomfilter")
	key = flag.String("key", keyName, "the name of the claim to inspect for revocations")
	flag.Parse()
}
