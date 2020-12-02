package chat

import (
	context "context"
	"fmt"
	"io/ioutil"
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

func guardarChunk(chu *BookChunk, libro string) {
	newFileName := libro + "_" + strconv.Itoa(int(chu.GetPieza()))
	//var fileSize int64
	const fileChunk = 256000
	//fileSize = chu.GetTam()
	//partSize := int(math.Min(fileChunk, float64(fileSize-int64(fileChunk))))
	_, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("no se pudo hacer el archivo " + newFileName)
	}
	ioutil.WriteFile(newFileName, chu.Data, os.ModeAppend)
}

//MandarAGuardar is funcion que guarda los chunks en memoria
func (s *Server) MandarAGuardar(ctx context.Context, chu *BookChunk) (*Message, error) {
	newFileName := chu.GetLibro() + "_" + strconv.Itoa(int(chu.GetPieza()))
	//var fileSize int64
	const fileChunk = 256000
	//fileSize = chu.GetTam()
	//partSize := int(math.Min(fileChunk, float64(fileSize-int64(fileChunk))))
	_, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("no se pudo hacer el archivo " + newFileName)
	}
	ioutil.WriteFile(newFileName, chu.Data, os.ModeAppend)

	ret := Message{
		Body:         "lo guardamos rey",
		Confirmacion: int32(1),
	}
	return &ret, nil
}

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

	var prop = crearPropuesta(&aux2)
	cn.Proponer(context.Background(), prop)

	ret := Message{
		Body: "guardado mi rey",
	}
	return &ret, nil
}

//Proponer is
func (s *Server) Proponer(ctx context.Context, message *Propuesta) (*Propuesta, error) {
	//revvisar datanode1
	var on1, on2, on3 = true, true, true
	if message.Cnod1 > 0 {
		var conn1 *grpc.ClientConn
		conn1, err1 := grpc.Dial(":9001", grpc.WithInsecure())
		if err1 != nil {
			on1 = false
		} else {
			defer conn1.Close()
		}

	}
	//revisar datanode2
	if message.Cnod2 > 0 {
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())
		if err2 != nil {
			on2 = false
		} else {
			defer conn2.Close()
		}
	}
	//revisar datanode3
	if message.Cnod3 > 0 {
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
		if err3 != nil {
			on3 = false
		} else {
			defer conn3.Close()
		}
	}
	rand.Seed(time.Now().Unix())

	var dn1, dn2, dn3, dnt = []int32{}, []int32{}, []int32{}, []int32{}

	if (message.Cnod1 > 0 && !on1) || (message.Cnod2 > 0 && !on2) || (message.Cnod3 > 0 && !on3) {
		//for que hace la prop con los nodos que tienen on = true
		for i := 0; i < len(message.Lnodt); i++ {
			n := rand.Int() % 3
			fmt.Println(n)
			// este switch igual debiera revisar si el nodo que se elige esta bueno o no pero no lo hecho todavia
			switch n {
			case 0:
				dn1 = append(dn1, int32(i))
				dnt = append(dnt, int32(n))
			case 1:
				dn2 = append(dn2, int32(i))
				dnt = append(dnt, int32(n))
			case 2:
				dn3 = append(dn3, int32(i))
				dnt = append(dnt, int32(n))
			}
		}
	}
	ret := Propuesta{
		Nombrelibro: message.GetNombrelibro(),
		Cnod1:       int32(len(dn1)),
		Cnod2:       int32(len(dn2)),
		Cnod3:       int32(len(dn3)),
		Lnod1:       dn1,
		Lnod2:       dn2,
		Lnod3:       dn3,
		Lnodt:       dnt,
	}
	if logflag == false {
		createFile()
	}
	writeFile(message)
	return &ret, nil

}

//Repartir is funcion que siempre es llamada por el datanode0
func (s *Server) Repartir(ctx context.Context, message *ListaPropuesta) (*Message, error) {
	var dn2, dn3 ChatServiceClient

	//revisar datanode1
	if message.Prop.Cnod2 > 0 {
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())
		if err2 != nil {
			log.Fatalf("Could not connect: %s", err2)
			ret := Message{Body: "ta malo larva 2"}
			return &ret, nil
		}
		defer conn2.Close()
		dn2 = NewChatServiceClient(conn2)
	}
	//revisar datanode2
	if message.Prop.Cnod3 > 0 {
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
		if err3 != nil {
			log.Fatalf("Could not connect: %s", err3)
			ret := Message{Body: "ta malo larva 3"}
			return &ret, nil
		}
		defer conn3.Close()
		dn3 = NewChatServiceClient(conn3)
	}

	for i := 0; i < len(message.Prop.GetLnodt()); i++ {
		switch message.Prop.GetLnodt()[i] {
		case 0:
			guardarChunk(message.Lista.Lista[i], message.Prop.GetNombrelibro())
		case 1:
			dn2.MandarAGuardar(context.Background(), message.Lista.Lista[i])
		case 2:
			dn3.MandarAGuardar(context.Background(), message.Lista.Lista[i])
		}
	}
	ret := Message{
		Body:         "lo logramos ivan",
		Confirmacion: 1,
	}
	return &ret, nil

}
