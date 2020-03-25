package godis

import (
	"bytes"
	"net"
)

// Client stores client info
type Client struct {
	CmdBuf string
}

// Server stores server info
type Server struct {
}

// CreateClient create a client
func (s *Server) CreateClient() *Client {
	return &Client{
		CmdBuf: "",
	}
}

// ReadClientContent read command from client
func (c *Client) ReadClientContent(conn net.Conn) error {

	buf := make([]byte, 10240)
	_, err := conn.Read(buf)
	if err != nil {
		log.Errorf("read client %+v cmd err:%+v", conn, err)
		return err
	}
	c.CmdBuf = string(buf)
	return nil
}

// TransClientContent convert the content from client into parameters
func (c *Client) TransClientContent() error {
	decoder := NewDecoder(bytes.NewReader([]byte(c.CmdBuf)))
	bulks, err := decoder.DecodeMultiBulks()
}
