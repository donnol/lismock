package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"reflect"
)

var (
	host, port string
	needVerify bool
)

func main() {
	log.Println("this is lis mock for test.")

	flag.StringVar(&host, "host", "127.0.0.1", "use --host to specify host")
	flag.StringVar(&port, "port", "8878", "use --port to specify port")
	flag.BoolVar(&needVerify, "verify", false, "must check the data received when enable")
	flag.Parse()

	// 监听tcp，接收hl7包并打印即可
	handleHL7ByTCP(host, port)
}

var (
	testHL7 = `MSH|^~\&|LCAMC|LABCORP|FAKE ACO|FAKE HOSPITAL|20170418093130||ORU^R01|M17108000000000001|P|2.3|||ER|ER` + "\r" +
		`PID|1|11111111111|11111111111|11111111111|Truong^Nicholas||19510610|M|||1511 MONTE VISTA ST^^PASADENA^CA^91106|||||||6337494512380` + "\r" +
		`PV1|1|O||||||1598879918^Khanh Hoang^Minh^^^^^^^^^^NPI|||||||||||||||||||||||||||||||SO` + "\r" +
		`ORC||633749453380|633749453380||||||20161202|||1598879918^Fake name^Fake given name^^^^^^^^^^NPI` + "\r" +
		`OBR|1||633749453380|005009^CBC WITH DIFFERENTIAL/PLATELET^L|||20161202|20161202||||||||1598879918^Fake name^Fake given name^^^^^^^^^^NPI||SO|||||||F` + "\r" +
		`OBX|1|ST|6690-2^LOINC^LN^005025^WBC^L||8.7|X10E3/UL|3.4-10.8||||F|||20161202|SO   ^^L` + "\r" +
		`OBX|1|ST|6690-2^LOINC^LN^005025^WBC^L||8.7|X10E3/UL|3.4-10.8||||F|||20161203|SO   ^^L` + "\r"

	testHL7DoubleQuote = "MSH|^~\\&|LCAMC|LABCORP|FAKE ACO|FAKE HOSPITAL|20170418093130||ORU^R01|M17108000000000001|P|2.3|||ER|ER\rPID|1|11111111111|11111111111|11111111111|Truong^Nicholas||19510610|M|||1511 MONTE VISTA ST^^PASADENA^CA^91106|||||||6337494512380\rPV1|1|O||||||1598879918^Khanh Hoang^Minh^^^^^^^^^^NPI|||||||||||||||||||||||||||||||SO\rORC||633749453380|633749453380||||||20161202|||1598879918^Fake name^Fake given name^^^^^^^^^^NPI\rOBR|1||633749453380|005009^CBC WITH DIFFERENTIAL/PLATELET^L|||20161202|20161202||||||||1598879918^Fake name^Fake given name^^^^^^^^^^NPI||SO|||||||F\rOBX|1|ST|6690-2^LOINC^LN^005025^WBC^L||8.7|X10E3/UL|3.4-10.8||||F|||20161202|SO   ^^L\rOBX|1|ST|6690-2^LOINC^LN^005025^WBC^L||8.7|X10E3/UL|3.4-10.8||||F|||20161203|SO   ^^L\r"
)

var (
	BeginByte = []byte{11}
	EndBytes  = []byte{13, 28, 13}
)

func WrapMessage(data []byte) []byte {
	res := make([]byte, 0, len(data)+4)
	res = append(res, BeginByte...)
	res = append(res, data...)
	res = append(res, EndBytes...)
	return res
}

func UnwrapMessage(data []byte) []byte {
	l := len(data)
	return data[len(BeginByte) : l-len(EndBytes)]
}

func handleHL7ByTCP(host, port string) {
	addr := host + ":" + port
	log.Printf("listen: %s\n", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept failed: %+v\n", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			var resp []byte
			for {
				buf := make([]byte, 2048)
				n, err := conn.Read(buf[:])
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				resp = append(resp, buf[:n]...)
				if bytes.Contains(buf[:n], EndBytes) {
					break
				}
			}

			rawresp := UnwrapMessage(resp)
			log.Printf("read len: %d, data: %q\n", len(resp), rawresp)

			if needVerify && !reflect.DeepEqual(rawresp, []byte(testHL7)) {
				log.Printf("bad result: %q \n!=\n %q\n", rawresp, testHL7)
			}
		}(conn)
	}
}
