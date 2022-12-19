package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea2-Grupo07/Proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

//Esta funcion solo guarda la informacion del NameNode

func (s *server) Nodo(ctx context.Context, msg *pb.Answer) (*pb.Answer, error) {
	fmt.Println("Se ha recibido la informacion desde el DataNode")
	fmt.Println("Preparandose para guardar la informacion")
	guardarinfo(msg.Body)
	return &pb.Answer{Body: "Perfecto"}, nil
}

//Esta funcion obtiene lo que haya escrito en el txt

func (s *server) Informacion(ctx context.Context, msg *pb.Answer) (*pb.Answer, error) {
	fmt.Println("Es la hora *chan chan* de buscar en los archivos ")
	path := "DATA.txt"
	info := ""
	leer, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Solo se puede leer!
	//Manejo de errores
	if err != nil {
		log.Fatal(err)
	}
	defer leer.Close()
	//Leer el archivo
	check := bufio.NewScanner(leer)
	check.Split(bufio.ScanLines)
	for check.Scan() {
		split := strings.Split(check.Text(), ":") //TIPO:ID:NOMBRENODO
		id := strings.TrimSpace(split[1])
		if id == msg.Body {
			info = split[2]
		}
	}
	return &pb.Answer{Body: info}, nil
}

func guardarinfo(msg string) {
	path := "./Grunt/DATA.txt"
	_, err := os.Stat(path) // -> Busca si existe en esa direccion el archivo! (Lo guardamos en err)

	if errors.Is(err, os.ErrNotExist) { //? Si el archivo NO existe, entramos a este if
		f, err := os.Create(path) //Lo creamos :D

		//? Manejo de errores!
		if err != nil {
			log.Fatal("ERROR en el archivo D: ", err)
		}
		defer f.Close()
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Abrimos el archivo en modo append, para agregarle más informacion!

	//? Manejo de errores nuevamente!
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()                    //Defer ya que dejamos de usar el archivo y se cierra solo :D
	_, err = f.WriteString(msg + "\n") //Aqui se escribe el mensaje con ul salto de linea

	//Nuevamente hay que hacer manejo de errores
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("El mensaje fue escrito de manera correcta y satisfactoria :'D")
	return

}

func main() {

	for {

		//Conexion sincrona entre Grunt y NameNode
		listener, err := net.Listen("tcp", ":50052") //conexion sincrona

		if err != nil {
			panic("Error en la conexion, no se pudo crear :c" + err.Error())
		}
		serv := grpc.NewServer()
		fmt.Println("Escuchando...")
		for {
			pb.RegisterMessageServiceServer(serv, &server{})
			if err = serv.Serve(listener); err != nil {
				break
			}
		}
	}
}
