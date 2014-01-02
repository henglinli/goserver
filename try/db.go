package db
//
import (
	//	"fmt"
	"code.google.com/p/leveldb-go/leveldb"
	//"code.google.com/p/leveldb-go/leveldb/db"
	//"code.google.com/p/leveldb-go/leveldb/memfs"
)
//
const (
	kData = "data"
)
/*
func main() {
	//fs := memfs.New()
	db, error := leveldb.Open("testdb", nil)
	if nil != error {
		fmt.Println("open testdb done")
		db.Close()
	}
}
*/
func Open(
