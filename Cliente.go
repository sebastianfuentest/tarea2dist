package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"papa.com/chat"
)

//Chunkear is funcion que chunkea
func Chunkear(Libro string) (ret [][]byte) {
	aux := [][]byte{}
	fileToBeChunked := "./" + Libro

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 256000

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := Libro + "_" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		aux = append(aux, partBuffer)
		fmt.Println("Split to : ", fileName)
	}
	return aux
}

//DeChunkear is funcion que reconstruye los chunks
func DeChunkear(Libro string, ctot uint64) {

	newFileName := "reconstituido" + Libro
	file, err := os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	const fileChunk = 256000
	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int64 = 0

	for j := uint64(0); j < ctot; j++ {

		//read a chunk
		currentChunkFileName := Libro + "_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

		n, err := file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		fmt.Println("Written ", n, " bytes")

		fmt.Println("Recombining part [", j, "] into : ", newFileName)
	}

	// now, we close the newFileName
	file.Close()

}
func main() {
	rand.Seed(time.Now().Unix())
	var conn1 *grpc.ClientConn
	conn1, err1 := grpc.Dial("dist50:9001", grpc.WithInsecure())

	if err1 != nil {
		log.Fatalf("Could not connect: %s", err1)
	}
	defer conn1.Close()
	c1 := chat.NewChatServiceClient(conn1)

	var connN *grpc.ClientConn
	connN, errn := grpc.Dial("dist49:9000", grpc.WithInsecure())

	if errn != nil {
		log.Fatalf("Could not connect: %s", errn)
	}
	defer connN.Close()
	cn := chat.NewChatServiceClient(connN)
	/*
		/*
			var conn2 *grpc.ClientConn
			conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())

			if err2 != nil {
				log.Fatalf("Could not connect: %s", err2)
			}
			defer conn2.Close()

			var conn3 *grpc.ClientConn
			conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())

			if err3 != nil {
				log.Fatalf("Could not connect: %s", err3)
			}
			defer conn3.Close()
	*/
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese una opción \n1.Subir archivo\n2.Descargar archivo\n3.Pedir lista de libros\n4.Terminar el programa")
	log.Printf("---------------------")
	var libch = [][]byte{}
	var i = 0
	powar := chat.ListaChunks{}
	for {
		//powar := chat.ListaChunks{}
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.Trim(text, " \r\n"))

		if strings.Compare(text, "1") == 0 {
			fmt.Println("Ingrese nombre del libro a subir\n")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, " \r\n")
			libch = Chunkear(text)

			powar.Libro = text

			for i = 0; i < len(libch); i++ {
				amandar := chat.BookChunk{
					Data:  libch[i],
					Pieza: int32(i),
				}

				powar.Lista = append(powar.Lista, &(amandar))

				fmt.Println(len(powar.Lista))
			}

			c1.SubirLibro(context.Background(), &powar)

			/*n := rand.Int() % 3
			switch n {
			case 1:
				c1.SubirLibro(context.Background(), powar)
			}*/

		}
		if strings.Compare(text, "2") == 0 {
			fmt.Println("Ingrese nombre del libro a descargar\n")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\r\n")
			//fmt.Println(powar.Lista[1])
			DeChunkear(text, uint64(len(powar.Lista)))
		}
		if strings.Compare(text, "4") == 0 {
			fmt.Println("bye bye\n")
			break
		}
		if strings.Compare(text, "3") == 0 {
			uwu := chat.Basura{
				Body: "UnU",
			}
			var LALISTA, _ = cn.ListarLibros(context.Background(), &uwu)
			for _, eachln := range LALISTA.Libros {
				fmt.Println(eachln)
			}
		}

		fmt.Println("Ingrese una opción \n1.Subir archivo\n2.Descargar archivo\n3.Pedir lista de libros\n4.Terminar el programa")
		log.Printf("---------------------")
	}

}
