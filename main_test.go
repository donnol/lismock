package main

import (
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8878")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testHL7, testHL7DoubleQuote) {
		t.Fatalf("%s != %s", testHL7, testHL7DoubleQuote)
	}

	t.Logf("%s", strconv.Quote(testHL7))
	t.Logf("%q", testHL7)
	t.Logf("%q", []byte(testHL7))
	t.Logf("%#q", testHL7)

	// 下面两种打印方式，在终端看时会少了东西
	t.Logf("%s", testHL7)
	t.Logf("%s", testHL7DoubleQuote)

	data := WrapMessage([]byte(testHL7))
	n, err := conn.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("write data len: %d, data: %s", n, data)

	time.Sleep(2 * time.Second)
}
