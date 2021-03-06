package godis

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const (
	DefaultSCFFile = "./SCF/01.scf"
)

type cmdFunc func(c *Client, s *Server)

// Client stores client info
type Client struct {
	Query       *Sdshdr
	Command     *GodisCommand
	Argc        int
	Argv        []*Object
	Db          *GodisDB
	Buf         *Sdshdr
	VirtualFlag bool
}

// GodisDB ...
type GodisDB struct {
	Dt      *Dict // stores keys
	Expires *Dict // for timeout keys
	ID      int   // DB id
}

// Server stores server info
type Server struct {
	Db          []*GodisDB
	DbNum       int
	Port        int
	Clients     int32
	Pid         int
	Commands    *Dict
	Dirty       int64
	SCFFileName string
}

// GodisCommand ...
type GodisCommand struct {
	Name *Sdshdr
	Proc cmdFunc
}

// CreateClient create a client
func (s *Server) CreateClient() *Client {
	return &Client{
		Db: s.Db[0],
	}
}

// InitServer ...
func InitServer() *Server {
	ParseConf()
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
	s.SCFFileName = DefaultSCFFile
	addCmdFuncs(s)

	LoadData(s)

	return s
}

// LoadData ...
func LoadData(s *Server) {
	c := s.CreateClient()
	c.VirtualFlag = true
	cmds := ReadSCF(s.SCFFileName)
	for _, cmd := range cmds {
		c.Query = SdsNewString(cmd)
		err := c.TransClientContent()
		if err != nil {
			log.Errorf("process input buffer error:%+v", err)
		}
		s.ProcessCommand(c)
	}
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
	dirty := s.Dirty
	c.Command.Proc(c, s)
	if dirty < s.Dirty && !c.VirtualFlag {
		AppendToSCF(s.SCFFileName, *(c.Query.SdsGetString()))
	}
}

// LookUpCommand return the cmd if name is a cmd
func (s *Server) LookUpCommand(name string) *GodisCommand {
	if v := s.Commands.Get(NewObject(OBJSDS, SdsNewString(name))); v != nil {
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

func addReplyStatus(c *Client, v string) {
	e := NewStatus([]byte(v))
	addReply(c, e)
}

func addReplyInt(c *Client, v int64) {
	s := fmt.Sprintf("(integer) %d", v)
	e := NewInt([]byte(s))
	addReply(c, e)
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

func addCmdFuncs(s *Server) {
	cmds := []GodisCommand{
		GodisCommand{
			Name: SdsNewString("lpush"),
			Proc: LPushCommand,
		},
		GodisCommand{
			Name: SdsNewString("rpush"),
			Proc: RPushCommand,
		},
		GodisCommand{
			Name: SdsNewString("llen"),
			Proc: LLenCommand,
		},
		GodisCommand{
			Name: SdsNewString("lrange"),
			Proc: LRangeCommand,
		},
		GodisCommand{
			Name: SdsNewString("set"),
			Proc: SetCommand,
		},
		GodisCommand{
			Name: SdsNewString("get"),
			Proc: GetCommand,
		},
		GodisCommand{
			Name: SdsNewString("incr"),
			Proc: IncrCommand,
		},
		GodisCommand{
			Name: SdsNewString("zadd"),
			Proc: ZaddCommand,
		},
		GodisCommand{
			Name: SdsNewString("zscore"),
			Proc: ZscoreCommand,
		},
		GodisCommand{
			Name: SdsNewString("zrange"),
			Proc: ZrangeCommand,
		},
		GodisCommand{
			Name: SdsNewString("zrank"),
			Proc: ZrankCommand,
		},
		GodisCommand{
			Name: SdsNewString("zrem"),
			Proc: ZremCommand,
		},
	}
	for i := range cmds {
		s.Commands.Add(NewObject(OBJSDS, cmds[i].Name), NewObject(OBJCommand, &cmds[i]))
	}
}
