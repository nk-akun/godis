package godis

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

type cmdFunc func(c *Client, s *Server)

// Client stores client info
type Client struct {
	Query   *Sdshdr
	Command *GodisCommand
	Argc    int
	Argv    []*Object
	Db      *GodisDB
	Buf     *Sdshdr
}

// GodisDB ...
type GodisDB struct {
	Dt      *Dict // stores keys
	Expires *Dict // for timeout keys
	ID      int   // DB id
}

// Server stores server info
type Server struct {
	Db       []*GodisDB
	DbNum    int
	Port     int
	Clients  int32
	Pid      int
	Commands *Dict
	Dirty    int64
}

// GodisCommand ...
type GodisCommand struct {
	Name *Sdshdr
	Proc cmdFunc
}

// CreateClient create a client
func (s *Server) CreateClient() *Client {
	return &Client{}
}

// InitServer ...
func InitServer() *Server {
	s := new(Server)
	s.DbNum = 8
	s.Db = make([]*GodisDB, s.DbNum)
	for i := 0; i < s.DbNum; i++ {
		s.Db[i] = InitDB(i)
	}
	df := &DictFunc{
		calHash:    CalHashCommon,
		keyCompare: CompareValueCommon,
	}
	s.Commands = NewDict(df)
	return s
}

// InitDB ...
func InitDB(id int) *GodisDB {
	db := new(GodisDB)
	df := &DictFunc{
		calHash:    CalHashCommon,
		keyCompare: CompareValueCommon,
	}
	db.Dt = NewDict(df)
	db.Expires = NewDict(df)
	db.ID = id
	return db
}

// ProcessCommand ...
func (s *Server) ProcessCommand(c *Client) {
	name, ok := c.Argv[0].Ptr.(string)
	if !ok {
		log.Errorf("cmd error:%v", ok)
		os.Exit(1)
	}
	cmd := s.LookUpCommand(name)
	if cmd == nil {
		addReplyError(c, fmt.Sprintf("error: Unknown command %s", name))
	} else {
		c.Command = cmd
		process(c, s)
	}
}

func process(c *Client, s *Server) {
	c.Command.Proc(c, s)
}

// LookUpCommand return the cmd if name is a cmd
func (s *Server) LookUpCommand(name string) *GodisCommand {
	if v := s.Commands.Get(NewObject(OBJString, name)); v != nil {
		cmd, ok := v.Ptr.(*GodisCommand)
		if !ok {
			return nil
		}
		return cmd
	}
	return nil
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
	c.Query = SdsNewBuf(buf)
	return nil
}

// TransClientContent convert the content from client into parameters
func (c *Client) TransClientContent() error {
	decoder := NewDecoder(bytes.NewReader(c.Query.SdsGetBuf()))
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

func addReplyError(c *Client, v string) {
	e := NewError([]byte(v))
	addReply(c, e)
}

func addReply(c *Client, e *EncodeData) {
	if s, err := EncodeMultiBulk(e); err == nil {
		c.Buf = SdsNewBuf(s)
	}
}
