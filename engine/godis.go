package godis

import (
	"bytes"
	"net"
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
