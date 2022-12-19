package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea1-Grupo07/Proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

// Variables globales a usar
var resuelto string // Define si es que una vez se declaro el estallido, se resolvio o no.
var serv *grpc.Server

/*
En este caso estamos trabajando sobre el laboratorio, que sera el servidor principal, por lo que tenemos que
definir el nuevo tipo de dato con el cual estaremos trabajando, en este caso es un struct de respuesta, el
cual no esta implementado.
*/

type server struct {
	pb.UnimplementedMessageServiceServer
}

/*
Aca definimos la funcion receptora del servicio que usaremos

Receptor:
- s: Puntero tipo server

Parametros:
- ctx: Contexto
- msg: Protobuff mensaje -> Tiene que ver con lo que se definio como parametro en el message.proto

Retorno:
- *pb.Message: Puntero a un tipo Message, que es el retorno que definimos en message.proto
- error: Estandar en caso de error al retornar
*/
func (s *server) Intercambio(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	fmt.Printf("Llega Escuadrón %d, conteniendo estallido...\n", msg.Equipo)
	if resolucion() == 1 {
		fmt.Println("Revisando estado Escuadrón: [LISTO]")
		resuelto = "SI"
		fmt.Printf("Estallido contenido, Escuadrón %d Retornando...\n", msg.Equipo)
		//defer serv.Stop()
	} else {
		fmt.Println("Revisando estado Escuadrón: [NO LISTO]")
		resuelto = "NO"
	}
	return &pb.Message{Body: resuelto, Equipo: msg.Equipo}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
	// Aca se retorna, entre llaves se deben entregar los elementos que componen al struct usado
}

/*
Aca definimos la funcion que revisará si el equipo está listo o no

Receptor:
- s: Puntero tipo server

Parametros:
- ctx: Contexto
- msg: Protobuff mensaje -> Tiene que ver con lo que se definio como parametro en el message.proto

Retorno:
- *pb.Message: Puntero a un tipo Message, que es el retorno que definimos en message.proto
- error: Estandar en caso de error al retornar
*/
func (s *server) Revision(ctx context.Context, msg *pb.Revisando) (*pb.Revisando, error) {
	fmt.Printf("Consulta para el equipo %d: %s \n", msg.Equipo, msg.Mensaje)
	if resolucion() == 1 {
		fmt.Println("Revisando estado Escuadrón: [LISTO]")
		resuelto = "SI"
		fmt.Printf("Estallido contenido, Escuadrón %d Retornando...\n", msg.Equipo)
		//defer serv.Stop()
	} else {
		fmt.Println("Revisando estado Escuadrón: [NO LISTO]")
		resuelto = "NO"
	}
	return &pb.Revisando{Mensaje: resuelto, Equipo: msg.Equipo}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
	// Aca se retorna, entre llaves se deben entregar los elementos que componen al struct usado
}

/*
func estallido()
Se define la funcion estallido que nos dice si un estallido va a ocurrir o no.

Receptor:
- No hay

Parametros:
- No hay

Retorno:
- Retorna un valor entero. 1 para estallidos, 0 por si no hubo alguno.
*/
func estallido() int {
	rand.Seed(time.Now().UTC().UnixNano())
	chance := rand.Intn(10) //comienza del 0 hasta el 9
	if chance <= 7 {
		return 1 // Return 1 para estallidos
	}
	return 0 // 0 si no hay Estallidos.
}

/*
func resolucion()
Se define la funcion resolucion, que nos dice si un estallido es contenido o no
Receptor:
- No hay

Parametros:
- No hay

Retorno:
- Retorna un valor entero. 1 para equipo listo para contener, 0 por si no esta listo el equipo.
*/
func resolucion() int {
	rand.Seed(time.Now().UTC().UnixNano())
	chance := rand.Intn(10) //comienza del 0 hasta el 9
	if chance <= 5 {
		return 1 //? return 1 para LISTO
	}
	return 0 //? return 0 para NO LISTO

}

/*
Función que manda un mensaje a través de rabbitmq para que siga funcionando

Receptor:
- No hay

Parametros:
- ch *amqp.Channel: Se utiliza para mandar el mensaje de rabbitmq
- qName string: El nombre de la cola
- LabName string: el nombre del laboratorio

Retorno:
- *pb.Message: Puntero a un tipo Message, que es el retorno que definimos en message.proto
- error: Estandar en caso de error al retornar
*/
func msjerabbit(ch *amqp.Channel, qName string, LabName string) {

	ch.Publish("", qName, false, false,
		amqp.Publishing{
			Headers:     nil,
			ContentType: "text/plain",
			Body:        []byte(LabName), //Contenido del mensaje convertido a bytes
		})
}

func main() {
	LabName := "Laboratorio Kampala - Uganda"                      //nombre del laboratorio
	hostQ := "dist025"                                             //ip del servidor de RabbitMQ 172.17.0.1
	qName := "Emergencias"                                         //nombre de la cola
	connQ, err := amqp.Dial("amqp://test:test@" + hostQ + ":5672") //conexion con RabbitMQ

	if err != nil {
		log.Fatal(err)
	} //? Manejo de errores de conexion con RabbitMQ
	defer connQ.Close()

	ch, err := connQ.Channel() //? Abre un canal de servidor único y concurrente para procesar la mayor parte de los mensajes.
	if err != nil {
		log.Fatal(err)
	} //? Manejo de errores de apertura de canal para procesamiento de msjes.
	defer ch.Close() // Se "aplaza" el cierre de la conexion

	for {
		if estallido() == 1 {
			resuelto = "NO"
			fmt.Println("Analizando estado Laboratorio: [ESTALLIDO]")
			fmt.Println("SOS Enviado a Central. Esperando respuesta...")
			//Mensaje enviado a la cola de RabbitMQ (Llamado de emergencia)
			err = ch.Publish("", qName, false, false,
				amqp.Publishing{
					Headers:     nil,
					ContentType: "text/plain",
					Body:        []byte(LabName), //Contenido del mensaje convertido a bytes
				})
			if err != nil {
				log.Fatal(err)
			}
			// Fin de mensaje, RabbitMQ-mensaje lab->central ~ asincrono

			// Inicio respuesta, gRPC-mensaje central->lab ~ sincrono
			listener, err := net.Listen("tcp", ":50051") // Listener conexion sincrona
			if err != nil {
				panic("La conexion no se pudo crear" + err.Error())
			}

			for {
				serv = grpc.NewServer()
				pb.RegisterMessageServiceServer(serv, &server{}) // Aca se une el servidor levantado recien, con la funcion que recibe las solicitudes
				go serv.Serve(listener)
				time.Sleep(5 * time.Second) //espera de 5 segundos
				defer serv.Stop()           //Cierra la conexión
				msjerabbit(ch, qName, LabName)
			}
		} else {
			fmt.Println("Analizando estado Laboratorio: [OK]")
		}
	}
}
