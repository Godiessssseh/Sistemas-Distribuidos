package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea1-Grupo07/Proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

// Definición de una Cola
type Cola struct {
	Elements []int
	Size     int
}

/*
Funcion que encola el elemento indicado en la cola
Receptor:
- q: Puntero tipo Cola

Parametros:
- elem int: valor a posicional en la cola

Retorno:
- No hay
*/
func (q *Cola) Enqueue(elem int) {
	q.Elements = append(q.Elements, elem)
}

/*
Funcion que desencola y aplica pop al elemento indicado en la cola
Receptor:
- q: Puntero tipo Cola

Parametros:
- No hay
Retorno:
- El elemento que se saco de la cola
*/
func (q *Cola) Dequeue() int {
	if q.IsEmpty() {
		fmt.Println("UnderFlow")
		return 0
	}
	element := q.Elements[0]
	if q.GetLength() == 1 {
		q.Elements = nil
		return element
	}
	q.Elements = q.Elements[1:]
	return element
}

/*
Funcion que obtiene el primer valor de la cola
Receptor:
- q: Puntero tipo Cola

Parametros:
- No hay
Retorno:
- El primer elemento de la cola
*/
func (q *Cola) GetTop() int {
	return q.Elements[0]
}

/*
Funcion que obtiene el largo de la cola
Receptor:
- q: Puntero tipo Cola

Parametros:
- No hay

Retorno:
- El largo de la cola
*/
func (q *Cola) GetLength() int {
	return len(q.Elements)
}

/*
Funcion que revisa si la cola está vacia o no
Receptor:
- q: Puntero tipo Cola

Parametros:
- No hay

Retorno:
- True si esta vacia, false si no lo esta.
*/
func (q *Cola) IsEmpty() bool {
	return len(q.Elements) == 0
}

// Definir variables globales
var candado sync.Mutex //Para sincronizar y bloquear variables
var Equipos Cola       //Inicializamos la cola que se usara

/*
Funcion que resuelve las emergencias, recibiendo el estado de la escuadra - el status y el retorno de la misma.
Receptor:
- No hay.

Parametros:
- f *os.File: Para escribir en el archivo
- emergencia ampq.Delivery: Se utiliza para obtener el nombre del laboratorio
- serviceCLiente pb.MessageServiceCLient: Mensajería usando el protoc
- equipoDesignado int32: Numero del equipo designado
- wg2 *sync.Waitgroup: Para cerrar la espera del wait group 2 (wg2.Done())
- LabName string: Nombre del Laboratorio
- wg *sync.Waitgroup: Para cerrar la espera del wait group 1 (wg.Done())

Retorno:
- No hay.
*/
func ResolverEmergencia(f *os.File, emergencia amqp.Delivery, serviceCliente pb.MessageServiceClient, equipoDesignado int32, wg2 *sync.WaitGroup, LabName string, wg *sync.WaitGroup) {
	candado.Lock()
	soli := 0
	fmt.Printf("Se envía Escuadra %d a "+string(emergencia.Body)+"\n", equipoDesignado)
	res, err := serviceCliente.Intercambio(context.Background(),
		&pb.Message{
			Body:   "Equipo Designado",
			Equipo: equipoDesignado,
		})

	if err != nil {
		panic("No se puede crear el mensaje " + err.Error())
	}
	fmt.Printf("Status Escuadra %d: "+res.Body+"\n", res.Equipo)
	soli += 1
	if res.Body == "SI" {
		fmt.Printf("Retorno a Central Escuadra %d, Conexion "+string(emergencia.Body)+" Cerrada"+"\n\n", res.Equipo)
		Equipos.Enqueue(int(res.Equipo))
	} else {
		for {
			res, err := serviceCliente.Revision(context.Background(),
				&pb.Revisando{
					Mensaje: "Equipo listo?",
					Equipo:  res.Equipo,
				})

			if err != nil {
				panic("No se puede crear el mensaje " + err.Error())
			}
			soli += 1
			fmt.Printf("Status Escuadra %d: %s \n", res.Equipo, res.Mensaje)
			if res.Mensaje == "SI" {
				fmt.Printf("Retorno a Central Escuadra %d, Conexion "+string(emergencia.Body)+" Cerrada"+"\n\n", res.Equipo)
				Equipos.Enqueue(int(res.Equipo))
				break
			}
			//respuesta del laboratorio
			time.Sleep(5 * time.Second) //espera de 5 segundos
		}
	}
	fmt.Fprintln(f, LabName+";"+strconv.Itoa(soli))
	candado.Unlock()
	wg.Done()
	wg2.Done()
}

/*
Funcion que manda a un equipo a resolver una emergencia.
Receptor:
- No hay.

Parametros:
- f *os.File: Para escribir en el archivo
- hostLab string: Tiene el host del laboratorio correspondiente
- emergencia ampq.Delivery: Se utiliza para obtener el nombre del laboratorio
- puerto string: numero del puerto entregado como string
- wg *sync.Waitgroup: Para cerrar la espera del wait group 1 (wg.Done())

Retorno:
- No hay.
*/
func AtenderEmergencia(f *os.File, hostLab string, emergencia amqp.Delivery, puerto string, wg *sync.WaitGroup) {

	connS, err := grpc.Dial(hostLab+puerto, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
	if err != nil {
		panic("No se pudo conectar con el servidor" + err.Error())
	}
	// defer connS.Close() // Se "aplaza" el cierre de la conexion
	serviceCliente := pb.NewMessageServiceClient(connS)

	//equipos disponibles, es posible acceder a ellos desde abajo? o insertarlos
	for {
		if !(Equipos.IsEmpty()) { // Si es que no hay equipos disponibles
			break
		}
	}
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	switch Equipos.Dequeue() {
	case 1:
		go ResolverEmergencia(f, emergencia, serviceCliente, 1, &wg2, string(emergencia.Body), wg)
		wg2.Wait()
	case 2:
		go ResolverEmergencia(f, emergencia, serviceCliente, 2, &wg2, string(emergencia.Body), wg)
		wg2.Wait()
	default:
		fmt.Println("Error al asignar equipo a la emergencia.")
	}
}

func main() {
	qName := "Emergencias" //Nombre de la cola
	hostQ := "localhost"   //Host de RabbitMQ 172.17.0.1
	dist025 := "localhost"
	port1 := ":50051"
	port2 := ":50052"
	port3 := ":50053"
	port4 := ":50054"
	dist026 := "dist026"
	dist027 := "dist027"
	dist028 := "dist028"
	//Host de un Laboratorio
	connQ, err := amqp.Dial("amqp://test:test@" + hostQ + ":5672") //Conexion con RabbitMQ
	Equipos.Enqueue(1)
	Equipos.Enqueue(2)
	if err != nil {
		log.Fatal(err)
	}
	defer connQ.Close()

	//Esto abre un canal de servidor único y concurrente para procesar la mayor parte de los mensajes.
	ch, err := connQ.Channel() //? ch=amqp channel
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(qName, false, false, false, false, nil) //Se crea la cola en RabbitMQ
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(q)

	fmt.Println("Esperando Emergencias")

	f, err := os.Create("Solicitudes.txt")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	chDelivery, err := ch.Consume(qName, "", true, false, false, false, nil) //? obtiene la cola de RabbitMQ
	if err != nil {
		log.Fatal(err)
	}

	// Ahora se revisara item por item de la cola
	for delivery := range chDelivery {
		fmt.Println("Mensaje asíncrono de " + string(delivery.Body) + " leído ") // Se avisa que se esta atendiendo un SOS
		wg := sync.WaitGroup{}
		wg.Add(1)
		if string(delivery.Body) == "Laboratorio Kampala - Uganda" {
			go AtenderEmergencia(f, dist025, delivery, port1, &wg)
			wg.Wait()
		} else if string(delivery.Body) == "Laboratorio Pohang - Korea" {
			go AtenderEmergencia(f, dist026, delivery, port2, &wg)
			wg.Wait()
		} else if string(delivery.Body) == "Laboratorio Renca - Chile" {
			go AtenderEmergencia(f, dist027, delivery, port3, &wg)
			wg.Wait()
		} else {
			go AtenderEmergencia(f, dist028, delivery, port4, &wg)
			wg.Wait()
		}
	}
}
