package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
)

var (
	root string
	addr string
)

func init() {
	flag.StringVar(&root, "root", "~/", "root of filesystem to serve over 9p")
	flag.StringVar(&addr, "addr", ":5640", "bind addr for 9p server, prefix with unix: for unix socket")
}

func main() {
	ctx := context.Background()
	log.SetFlags(0)
	flag.Parse()

	proto := "tcp"
	if strings.HasPrefix(addr, "unix:") {
		proto = "unix"
		addr = addr[5:]
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatalln("error listening:", err)
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln("error accepting:", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			ctx := context.WithValue(ctx, "conn", conn)
			log.Println("connected", conn.RemoteAddr())
			session, err := newLocalSession(ctx, root)
			if err != nil {
				log.Println("error creating session")
				return
			}

			if err := p9p.ServeConn(ctx, conn, p9p.Dispatch(session)); err != nil {
				log.Printf("serving conn: %v", err)
			}
		}(c)
	}
}

// newLocalSession returns a session to serve the local filesystem, restricted
// to the provided root.
func newLocalSession(ctx context.Context, root string) (p9p.Session, error) {
	// silly, just connect to ufs for now! replace this with real code later!
	log.Println("dialing", ":5640", "for", ctx.Value("conn"))
	conn, err := net.Dial("tcp", ":5640")
	if err != nil {
		return nil, err
	}

	session, err := p9p.NewSession(ctx, conn)
	if err != nil {
		return nil, err
	}

	return session, nil
}
