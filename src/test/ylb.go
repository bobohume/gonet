package main

import (
	"actor"
	"base"
	"db"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"message"
	"net"
	"net/http"
	"network"
	"runtime"
	"server/common"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type datetime int64

func sender(conn net.Conn) {
	worlds := "helloworld!"
	conn.Write([]byte(worlds))
	fmt.Println("send over")
}

type Node struct{
	i int
	pNext *Node
}

type Sqltest1 struct{
	MM int8
	MM1 uint8
}

type Sqltest struct{
	In *uint8 `sql:"primary;name:_i" redis:"name:_i"`
	J int8 `sql:"primary;name:_j" redis:"name:_i"`
	K string
	I2 []uint
	J2 []int
	Sqltest11 [1]Sqltest1
	K2 []string
	Bbb2 [2]bool
	Bbb21 [2]float64
	T1 int64 `sql:"datetime"`
	T int64 `sql:"datetime"`
}
type Server struct {
	ServerName string
	ServerIP   string
	ServerPort []int
	TestMM 		[2]int8
	TestMM1 	[2]string
	TestMM3		[2]int
	TestMM4		[]int
	Testtss     []string
	//sqltest []sqltest
}

func (this *Server) ReadData(b *base.BitStream){
	base.ReadData(this, b)
}

func (this *Server) WriteData(b *base.BitStream){
	base.WriteData(this, b)
}

type Serverslice struct {
	Servers []Server
}

type (
	CmdProcess struct {
		actor.Actor
	}

	ICmdProcess interface {
		actor.IActor
	}
)



func (this *CmdProcess) Init(num int) {
	for{
		i := 0
		i += 1
	}
	this.Actor.Init(num)
	this.RegisterCall("cpus", func() {
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("routines", func() {
		fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
	})

	this.RegisterCall("setcpus", func(args string) {
		n, _ := strconv.Atoi(args)
		runtime.GOMAXPROCS(n)
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("startgc", func() {
		runtime.GC()
		fmt.Println("gc finished")
	})

	this.Actor.Start()
}

type A struct {
	a int
	b int
}

type TOPRANKMAP [] A

func (this TOPRANKMAP) Len() int{
	return len(this)
}

func (this TOPRANKMAP) Less(i, j int) bool{
	return this[i].a > this[j].b
}

func (this TOPRANKMAP)Swap(i, j int){
	this[i], this[j] = this[j], this[i]
}

type TOPRANKSET map[int] int//排行榜队列

var lang = flag.String("lang", "golang", "the lang of the program")
func main() {
	ttt := make(map[int64] int64)
	base.UUID.Init(0)
	time1 := time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < 10000000; i++{
		n := base.UUID.UUID()
		_, bEx := ttt[n]
		if bEx{
			fmt.Println("重复uid")
			continue
		}
		ttt[n] = n
	}
	fmt.Println("end", time.Now().UnixNano() / int64(time.Millisecond) - time1)

	for {
		time.Sleep(1)
	}

	str := []byte("我是大另议11111222")
	fmt.Println(str)
	J := uint8(1)
	var1 :=Sqltest{&J, 2, "test", []uint{1, 2}, []int{3,4}, [1]Sqltest1{Sqltest1{1, 1}},  []string{"tes21", "tes31"}, [2]bool{false,true},[2]float64{1, 2.2}, time.Now().Unix(), time.Now().Unix(), }
	fmt.Println(db.UpdateSql(var1, "tb_test"))
	fmt.Println(db.UpdateSqlEx(var1, "tb_test", "_i", "J2"))
	fmt.Println(db.LoadSql(var1, "tb_test","playerid = 111"))
	fmt.Println(db.LoadSql(var1, "tb_test",""))
	fmt.Println(db.LoadSqlEx(var1,  "tb_test","playerid = 111", "_i", "J2",))
	fmt.Println(db.DeleteSql(var1, "tb_test"))
	fmt.Println(db.DeleteSqlEx(var1,  "tb_test", "_i", "J2",))
	fmt.Println(db.InsertSql(var1, "tb_test"))
	fmt.Println(db.InsertSqlEx(var1,  "tb_test", "_i", "J2",))
	//row := &db.Row{}
	//row.Init()
	//row.Set("_i", "100")
	//row.Set("_j", "101")
	//row.Set("mm", "102")

	//row.Obj(&var1)
	fmt.Println(var1)




	//这里实现了远程获取pprof数据的接口
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()


	map1 := make(map[int] int)
	mutex := sync.Mutex{}
	fmt.Println("111", time.Now().UnixNano())
	go func() {
		for i := 0; i < 100000; i++{
			mutex.Lock()
			map1[1] = 1000
			mutex.Unlock()
		}
		fmt.Println("222",time.Now().UnixNano())
	}()

	go func() {
		for i := 0; i < 100000; i++{
			mutex.Lock()
			map1[1] = 1000
			mutex.Unlock()
		}
		fmt.Println("333",time.Now().UnixNano())
	}()

	go func() {
		for i := 0; i < 100000; i++{
			mutex.Lock()
			map1[1]++
			mutex.Unlock()
		}
		fmt.Println("4444",time.Now().UnixNano())
	}()

	go func() {
		for i := 0; i < 100000; i++{
			mutex.Lock()
			map1[1]++
			mutex.Unlock()
		}
		fmt.Println("5555",time.Now().UnixNano())
	}()

	var a [3]int
	//fmt.Sscanf("1,2,3", "%d,%d,%d", &a[0], &a[1], &a[2])
	fmt.Sscanf("1,2,3", "%d,%d,%d", &a)
	fmt.Println(a)
	fmt.Println(base.GetNextTime(0))
	fmt.Println(base.GetNextTime(1))
	fmt.Println(base.GetNextTime(2))
	fmt.Println(base.GetNextTime(3))
	var map111 [2]TOPRANKSET
	map111[0] = make(TOPRANKSET)
	fmt.Println(map111[0])

	fmt.Printf("testte%d \n", int64(23372036854775808))
	fmt.Println("111")
	pCmd := &CmdProcess{}
	pCmd.Init(1)
	funcName := common.StartConsole
	var funcName1  *func(actor.IActor)
	fmt.Println(funcName, funcName1)
	ponit := unsafe.Pointer(&funcName)
	nPoint := (*int)(unsafe.Pointer(ponit))
	fmt.Println(*nPoint)
	func1 := (*func (actor.IActor))(unsafe.Pointer(nPoint))
	fmt.Println(func1, funcName)
	(*func1)(pCmd)
	//func1 := (funcName1)(ponit)
	//common.StartConsole(pCmd)
	flag.Parse()
	//var sss base.BitStream

	test := &message.C_A_LoginRequest1{
		Login:&message.C_A_LoginRequest{		PacketHead:message.BuildPacketHead(0,0 ),
			AccountName:proto.String("testt"), BuildNo:proto.String("test112"),},
	}
	//var sss1 base.Message
	sss1 := &Server{"test", "127.0.0.1", []int{1000, 2000}, [2]int8{10,20}, [2]string{"test11", "test22"}, [2]int{1,2}, []int{1022,123333}, []string{"111", "222"}}
	bs  := base.NewBitStream(make([]byte, 1024), 1024)
	/*base.RegisterMessage(&Server{}, func() base.Message{
		return &Server{}
	})*/
	base.RegisterMessage(&Server{})
	sss1.WriteData(bs)
	bs1  := base.NewBitStream(bs.GetBuffer(), 1024)
	packet := base.GetMessage("server")
	packet.ReadData(bs1)
	fmt.Println(packet, sss1)
	fmt.Println(message.GetPakcetHead(test))


	var server network.IServerSocket
	server = new(network.ServerSocket)
	server.Init("127.0.0.1", 11028)
	server.Start()
	//server.Stop()

	buf := make([]byte, 256)
	bitstream := base.NewBitStream(buf, 256)
	bitstream.WriteInt(1000, 16)
	bitstream.WriteInt(200000, 32)
	bitstream.WriteFlag(true)
	bitstream.WriteFlag(true)
	bitstream.WriteFlag(true)
	bitstream.WriteFlag(false)
	bitstream.WriteString("123456我的")
	bitstream.WriteInt64(1222222, 32)
	bitstream.WriteFloat(2.1)

	bitstream1 := base.NewBitStream(buf, 256)
	fmt.Println(bitstream1.ReadInt(16))
	fmt.Println(bitstream1.ReadInt(32))
	fmt.Println(bitstream1.ReadFlag())
	fmt.Println(bitstream1.ReadFlag())
	fmt.Println(bitstream1.ReadFlag())
	fmt.Println(bitstream1.ReadFlag())
	fmt.Println(bitstream1.ReadString(), "www")
	fmt.Println(bitstream1.ReadInt64(32))
	fmt.Println(bitstream1.ReadFloat())

	for {
		time.Sleep(1000)
	}
}
