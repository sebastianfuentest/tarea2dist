package chat

import (
	context "context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	grpc "google.golang.org/grpc"
)

//Server is
type Server struct {
	mut sync.Mutex
	id  int
}

//Chunk is struct de chunks nodeOrder es la lista de en que orden estan repartidas la distintas piezas en cada nodo
type Chunk struct {
	partes    uint64
	libro     string
	nodeOrder []int32
}

var path = "log.txt"
var listaChunks []Chunk
var logflag = false

func crearPropuesta(lista *ListaChunks) *Propuesta {
	rand.Seed(time.Now().Unix())

	var i int32 = 0
	var listaaux = lista.GetLista()
	var dn1, dn2, dn3, dnt = []int32{}, []int32{}, []int32{}, []int32{}
	for i = 0; i < int32(len(listaaux)); i++ {
		n := rand.Int() % 3
		fmt.Println(n)
		switch n {
		case 0:
			dn1 = append(dn1, i)
			dnt = append(dnt, int32(n))
		case 1:
			dn2 = append(dn2, i)
			dnt = append(dnt, int32(n))
		case 2:
			dn3 = append(dn3, i)
			dnt = append(dnt, int32(n))
		}
	}
	ret := Propuesta{
		Nombrelibro: lista.GetLibro(),
		Cnod1:       int32(len(dn1)),
		Cnod2:       int32(len(dn2)),
		Cnod3:       int32(len(dn3)),
		Lnod1:       dn1,
		Lnod2:       dn2,
		Lnod3:       dn3,
		Lnodt:       dnt,
	}

	return &ret
}

func createFile() {
	// check if file exists

	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	fmt.Println("File Created Successfully", path)
	logflag = true
}
func writeFile(prop *Propuesta) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	var i int32 = 0
	var ldn = prop.Lnodt
	if err != nil {
		return
	}
	defer file.Close()
	var lacosa string
	for i = 0; i < prop.Cnod1+prop.Cnod2+prop.Cnod3; i++ {
		lacosa = prop.Nombrelibro + " parte" + strconv.Itoa(int(i)) + " esta el nodo " + strconv.Itoa(int(ldn[i])) + "\n"
		_, err = file.WriteString(lacosa)
		if err != nil {
			return
		}
	}
	// Write some text line-by-line to file.
	// Save file changes.
	err = file.Sync()
	if err != nil {
		return
	}

	fmt.Println("File Updated Successfully.")
}

//SubirLibro is
func (s *Server) SubirLibro(ctx context.Context, message *ListaChunks) (*Message, error) {
	//llamar a namenode
	var connN *grpc.ClientConn
	aux2 := ListaChunks{
		Lista: message.GetLista(),
		Libro: message.GetLibro(),
	}
	connN, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer connN.Close()
	cn := NewChatServiceClient(connN)
	for {
		var prop = crearPropuesta(&aux2)
		aux1, _ := cn.Proponer(context.Background(), prop)
		if aux1.Confirmacion == 1 {
			fmt.Println("sali del loop de propuesta")
			break
		}
		fmt.Println("falle en crear una propuesta, reintentando ...")
	}
	ret := Message{
		Body: "guardado mi rey",
	}
	return &ret, nil
}

//Proponer is
func (s *Server) Proponer(ctx context.Context, message *Propuesta) (*Message, error) {
	//revvisar datanode1
	if message.Cnod1 > 0 {
		var conn1 *grpc.ClientConn
		conn1, err1 := grpc.Dial(":9001", grpc.WithInsecure())
		if err1 != nil {
			log.Fatalf("Could not connect: %s", err1)
			ret := Message{Body: "ta malo larva 1", Confirmacion: 0}
			return &ret, nil
		}
		defer conn1.Close()

	}
	//revisar datanode2
	if message.Cnod2 > 0 {
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())
		if err2 != nil {
			log.Fatalf("Could not connect: %s", err2)
			ret := Message{Body: "ta malo larva 2"}
			return &ret, nil
		}
		defer conn2.Close()

	}
	//revisar datanode3
	if message.Cnod3 > 0 {
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
		if err3 != nil {
			log.Fatalf("Could not connect: %s", err3)
			ret := Message{Body: "ta malo larva 3"}
			return &ret, nil
		}
		defer conn3.Close()

	}
	ret := Message{Body: "ta bien larva procede/do a guardarlo", Confirmacion: 1}
	if logflag == false {
		createFile()
	}
	writeFile(message)
	return &ret, nil

}

/*
//Repartir is
func (s *Server) Repartir(ctx context.Context, message *Propuesta) (*Message, error) {
	var dn1, dn2, dn3 = false, false, false

	//revisar datanode1
	if message.Cnod1 > 0 {
		var conn1 *grpc.ClientConn
		conn1, err1 := grpc.Dial(":9001", grpc.WithInsecure())
		if err1 != nil {
			log.Fatalf("Could not connect: %s", err1)
			ret := Message{Body: "ta malo larva 1", Confirmacion: 0}
			return &ret, nil
		}
		defer conn1.Close()
		dn1 = true
	}
	//revisar datanode2
	if message.Cnod2 > 0 {
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())
		if err2 != nil {
			log.Fatalf("Could not connect: %s", err2)
			ret := Message{Body: "ta malo larva 2"}
			return &ret, nil
		}
		defer conn2.Close()
		dn2 = true
	}
	//revisar datanode3
	if message.Cnod3 > 0 {
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
		if err3 != nil {
			log.Fatalf("Could not connect: %s", err3)
			ret := Message{Body: "ta malo larva 3"}
			return &ret, nil
		}
		defer conn3.Close()
		dn3 = true
	}

}*/
