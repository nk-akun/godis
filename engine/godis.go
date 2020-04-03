package godis

import (
	"bytes"
	"net"
)

// Client stores client info
type Client struct {
	CmdBuf string
	Argc   int
	Argv   []*Object
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
		conn.Close()
		return err
	}
	c.CmdBuf = string(buf)
	return nil
}

// TransClientContent convert the content from client into parameters
func (c *Client) TransClientContent() error {
	decoder := NewDecoder(bytes.NewReader([]byte(c.CmdBuf)))
	log.Info(c.CmdBuf)
	bulks, err := decoder.DecodeMultiBulks()
	if err != nil {
		log.Errorf("translate command error:%v", err)
		return err
	}

	c.Argc = len(bulks)
	c.Argv = make([]*Object, c.Argc)

	for i, bulk := range bulks {
		c.Argv[i] = NewObject(OBJString, string(bulk.Value))
	}
	return nil
}
