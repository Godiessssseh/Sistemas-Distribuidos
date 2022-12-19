package main

import (
	"bufio"
	"context"
	"errors" ////////////////////////
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings" ///////////////////
	"time"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea2-Grupo07/Proto"

	"google.golang.org/grpc"
)

var Nombre string
var Variable string

type server struct {
	pb.UnimplementedMessageServiceServer
}

//Función que ya utilizamos en la Tarea 1 pero con funcionalidades extras!
// Si se obtiene una R -> Proviene de Rebeldes, es decir, que tenemos que buscar en los txt el ID -> Ese ID se busca en los datanodes -> Obtenemos la info!
//Si se obtiene una C o cualquier otro caso -> Proviene de Combine -> Se agrega al TXT

func (s *server) Intercambio(ctx context.Context, msg *pb.Message) (*pb.Answer, error) {
	//! Caso cuando se recibe un mensaje de rebeldes
	if msg.Data == "R" {
		fmt.Println("Se recibe la consulta de los Rebeldes - La categoria es:", msg.Body)
		//Revisar las consultas de los rebeldes!
		path := "DATA.txt"
		//Abrir archivos y buscar con el ID y a que nodo fue enviado, para obtener la informacion de los datanodes
		id_rebeldes := ""
		port := ""
		IP := ""

		leer, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Solo se puede leer!
		if err != nil {
			log.Fatal(err)
		}
		defer leer.Close()
		read := bufio.NewScanner(leer)
		read.Split(bufio.ScanLines) //Lo corta por linea tras linea!
		for read.Scan() {           //Mientras exista informacion para leer
			//!La informacion que recibe msg.body es si es Logistica, Militar o Financiera
			separarinfolinea := strings.Split(read.Text(), " ")

			if separarinfolinea[0] == msg.Body { //Revisamos si el tipo es el mismo!
				if separarinfolinea[2] == "DataNodeGrunt" {
					id_rebeldes = separarinfolinea[1] //Usamos el ID para obtener la info
					port = ":50052"
					IP = "dist026"

				} else if separarinfolinea[2] == "DataNodeSynth" {
					id_rebeldes = separarinfolinea[1]
					port = ":50053"
					IP = "dist027"
				} else {
					id_rebeldes = separarinfolinea[1]
					port = ":50054"
					IP = "dist028"
				} //Luego trabajamos en la conexion con los datanode para utilizar el ID y obtener la informacion!
				connS, err := grpc.Dial(IP+port, grpc.WithInsecure()) //crea la conexion sincrona con el datanode correspondiente

				//Manejo de errores
				if err != nil {
					panic("Error al conectarse con el Datanode D: " + err.Error())
				}
				serviceCliente := pb.NewMessageServiceClient(connS)
				res, _ := serviceCliente.Informacion(context.Background(), //Busca la informacion utilizando el ID
					&pb.Answer{
						Body: separarinfolinea[1],
					})
				id_rebeldes = id_rebeldes + res.Body + "\n" //Se une todo hasta que se acabe el for!
				connS.Close()
			}
		}
		info_rebeldes := id_rebeldes
		id_rebeldes = ""                            //Se resetea en caso de :D
		return &pb.Answer{Body: info_rebeldes}, nil //Retornar la data a rebeldes

		//! Caso cuando se recibe un mensaje del combine

	} else {
		fmt.Println("Se recibe el mensaje de Combine - El mensaje es:", msg.Body)

		msg := msg.Body //? Guardar el mensaje en una variable global!
		//este mensaje se debe guardar en un txt! y se debe dar una respuesta!
		//Escribire el txt aca!
		TIPO := ""
		ID := ""
		//Mi plan es hacer un split con " :", debido a que el formato es -> "LOGISTICA : ID : MSJE "

		split := strings.Split(msg, " : ") //! Ahora tengo Logistica, ID separados! msje no lo necesito!
		TIPO = split[0]
		ID = split[1]

		path := "DATA.txt"
		Add := TIPO + ":" + ID + ":" + "DataNode" + ":" + Nombre //El mensaje a guardar dentro del archivo!
		//Crear archivo y escribir en el!
		_, err := os.Stat(path)             //Para revisar si existe o no!
		if errors.Is(err, os.ErrNotExist) { //Si el archivo NO existe entramos al if y lo creamos
			f, err := os.Create(path) // Se crea el archivo respectivo!
			//Manejo de errores!
			if err != nil {
				log.Fatal("El archivo no pudo ser creado D: ", err)
			}
			defer f.Close()
		}
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Solo nos interesa el append, para agregar informacion al archivo!
		//Manejo de errores!
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(Add + "\n") // Se escribe y se le da un salto de línea!
		if err != nil {
			log.Fatal(err)
		}
		Variable = Add
		fmt.Println("Informacion guardada con exito en el archivo")
		return &pb.Answer{Body: "Perfecto"}, nil
	}
}

/*
Aca definimos la funcion que creara la conexion sincrona
Parametros:
- ipe -> Tiene el ip a donde se tiene que enviar el mensaje
- puerto -> tiene el puerto a donde se tiene que enviar el mensaje
- msg -> Mensaje
- nombre -> Nombre del datanode al que se esta enviando la info
Retorno:
- *pb.Message: Puntero a un tipo Message, que es el retorno que definimos en message.proto
- error: Estandar en caso de error al retornar
*/

func conexionnymsg(ipe string, puerto string, msg string, nombre string) {
	connS, err := grpc.Dial(ipe+puerto, grpc.WithInsecure()) //crea la conexion sincrona con datanode
	//Manejo de errores!

	if err != nil {
		panic("No se puso establecer la conexion con el datanode D:" + err.Error())
	}
	fmt.Println("Enviando la informacion al datanode: " + nombre)
	serviceCliente := pb.NewMessageServiceClient(connS)
	res, _ := serviceCliente.Nodo(context.Background(),
		&pb.Answer{
			Body: msg,
		})
	//check a esto? no se si debe ir o no, ayuda
	if res.Body == "Perfecto" {
		fmt.Println("Informacion enviada de manera correcta :D")
	}
	connS.Close()
}

func main() {
	//revisar que puerto
	port := ""
	IP := ""
	Nombre := ""
	for {
		//tenemos que crear/elegir un valor al azar para cada nodo!!
		//0 == grunt, 1 == Synth, 2 == Cremator

		rand.Seed(time.Now().UnixNano())
		numero_azar := rand.Intn(3) //!No se si esto toma del 0 al 2 o si toma otros valores

		//comparar el numero al azar!
		if numero_azar == 0 {
			port = ":50052" //Puerto del Grunt
			Nombre = "Grunt"
			IP = "dist026"

		} else if numero_azar == 1 {
			port = ":50053" //Puerto del Synth
			Nombre = "Synth"
			IP = "dist027"

		} else {
			port = ":50054" //Puerto del Cremator
			Nombre = "Cremator"
			IP = "dist028"
		}

		//conexion sincrona para el Combine y Rebelde
		listener, err := net.Listen("tcp", ":50051")
		fmt.Println("Antes de manejo de errores")
		//manejo de errores!
		if err != nil {
			panic("Error en la conexion, no se pudo crear :c" + err.Error())
		}
		serv := grpc.NewServer()
		go conexionnymsg(IP, port, Variable, Nombre)
		for {
			pb.RegisterMessageServiceServer(serv, &server{})
			if err = serv.Serve(listener); err != nil {
				break
			}
		}
	}
}
