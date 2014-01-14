package message

import (
	"code.google.com/p/goprotobuf/proto"
	"log"
	"reflect"
	"testing"
)

func TestExtension(t *testing.T) {
	log.Println("test Extension...")

	request := &Request{}
	proto.SetDefaults(request)
	log.Println("request type:", request.GetType())
	log.Println("request command:", request.GetCommand())
	// has
	log.Println("has extension:",
		proto.HasExtension(request, E_Login_Account))
	// set
	account := &Account{}
	proto.SetDefaults(account)
	account.Name = proto.String("lee")
	account.Token = proto.String("lee")
	log.Println(account.Type, account.GetName(), account.GetToken())

	err := proto.SetExtension(request, E_Login_Account, account)
	if err != nil {
		t.Fatal(err.Error())
	}
	// has
	log.Println("has extension:",
		proto.HasExtension(request, E_Login_Account))
	// get
	var new_account interface{}
	new_account, err = proto.GetExtension(request, E_Login_Account)
	if err != nil {
		t.Fatal(err.Error())
	}
	{
		account := new_account.(*Account)
		log.Println("getname", account.GetName())
	}
	data_type := reflect.TypeOf(new_account)

	log.Println("type:", data_type)
	switch new_account.(type) {
	case *Account:
	default:
		t.Fatal("Bad Type")
	}
	value := reflect.ValueOf(new_account)

	log.Println("value:", value)
	GetName := value.MethodByName("GetName")
	log.Println("is reflect.Func", reflect.Func == GetName.Kind())
	name := GetName.Call(nil)

	name_str := name[0].String()
	log.Println("GetName:", reflect.TypeOf(name_str).Kind())
}
