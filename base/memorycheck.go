package base

import (
	_ "net/http/pprof"
	"log"
	"net/http"
)

type(
	MemoryCheck struct {

	}
)

//http://localhost:6060/debug/pprof/
//http://localhost:6060/debug/pprof/heap
//go tool pprof -inuse_space http://localhost:6060/debug/pprof/heap
//go tool pprof http://localhost:6060/debug/pprof/heap?debug=1
func (this *MemoryCheck) Init(){
	go func() {
	     log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}


// f, err := os.OpenFile("./tmp/cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
//    if err != nil {
//        log.Fatal(err)
//    }
//    defer f.Close()
//    pprof.StartCPUProfile(f)
//    defer pprof.StopCPUProfile()