package main
//
import (
	"log"
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
func main() {
	// open
	db_folder := "tmp.db"
	db, err := leveldb.Open(db_folder, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// set
	err = db.Set([]byte("lee"), []byte("lihenglin"), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// get
	var result []byte
	result, err = db.Get([]byte("lee"), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(result))
	db.Close()
}
