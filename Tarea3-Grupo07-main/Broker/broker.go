package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea3-Grupo07/Proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

func CrearConexion(direccion string, puerto string) pb.MessageServiceClient {
	hostS := direccion                                         //Host de un Laboratorio
	connS, err := grpc.Dial(hostS+puerto, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
	if err != nil {
		panic("No se pudo conectar con el servidor del planeta" + err.Error())
	}
	// defer connS.Close() // Se "aplaza" el cierre de la conexion
	serviceCliente := pb.NewMessageServiceClient(connS)
	return serviceCliente

}

/*
 */
func SoldadosQuery(serviceCliente pb.MessageServiceClient, sector string, base string) (int32, string, []int32) {
	res, err := serviceCliente.SoldadosQuery(context.Background(),
		&pb.SoldadosRequest{
			Sector: sector,
			Base:   base,
		})

	if err != nil {
		panic("No se puede crear el mensaje " + err.Error())
	}
	//println("Respuesta: %d \n", res.Codigo)//respuesta del
	if res.CantidadSoldados == -1 {
		fmt.Printf("Existe el archivo de registros del sector %s, en la base %s, pero no hay informacion sobre sus soldados.\n", sector, base)
	} else if res.CantidadSoldados == -2 {
		fmt.Printf("No existe un archivo de registros del sector %s\n", sector)
	} else if res.CantidadSoldados == -3 {
		fmt.Printf("Existe el archivo de registros del sector %s, pero no especificamente de la base %s\n", sector, base)
	} else {
		fmt.Println("Soldados recibidos desde : ", res.PlanetaEmisor, "El reloj recibido es: ", res.RelojVectorial)
	}
	return res.CantidadSoldados, res.PlanetaEmisor, res.RelojVectorial
}
func (s *server) GetSoldados(ctx context.Context, msg *pb.SoldadosRequest) (*pb.SoldadosAnswer, error) {
	fmt.Printf("Se solicitan la cantidad de soldados en el sector %s, base: %s. \n", msg.Sector, msg.Base)
	direccionPlaneta, puertoPlaneta := elegirazar()
	canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
	soldadosRespuesta, planetaRespuesta, relojPlaneta := SoldadosQuery(canalConPlaneta, msg.Sector, msg.Base)
	//defer serv.Stop()
	return &pb.SoldadosAnswer{CantidadSoldados: soldadosRespuesta, PlanetaEmisor: planetaRespuesta, RelojVectorial: relojPlaneta}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
	// Aca se retorna, entre llaves se deben entregar los elementos que componen al struct usado
}

func (s *server) GetPlanet(ctx context.Context, msg *pb.SolicitudPlaneta) (*pb.PlanetaDesignado, error) {
	fmt.Printf("Vanguardia solicita planeta para operacion de %s. \n", msg.Solicitud)
	DireccionPrueba, PuertoPrueba := elegirazar()

	//? Valores anteriores usados
	//DireccionPrueba := "localhost"
	//PuertoPrueba := ":50051"
	//defer serv.Stop()

	//! Dejo este print?
	fmt.Printf("Servidor enviado con direccion %s y puerto %s", DireccionPrueba, PuertoPrueba)
	return &pb.PlanetaDesignado{DireccionAsignada: DireccionPrueba, PuertoAsignado: PuertoPrueba}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
	// Aca se retorna, entre llaves se deben entregar los elementos que componen al struct usado
}

/*
Funcion que se encarga de elegir un server planeta al azar.
Receptor:
-	No hay

Parametros:
-	No hay

Retorno Doble:
- Retorna el host elegido con el puerto elegido.
*/
func elegirazar() (host string, puerto string) {
	l := rand.Intn(3) //Elige valor al azar entre 0, 1 y 2
	hostrand := ""
	puertorand := ""
	if l == 0 { //Servidor Dominante Tierra
		hostrand = "dist026" //!dist026 Cambiar despues!!
		puertorand = ":50052"
	} else if l == 1 { //Servidor Titan
		hostrand = "dist027" //!dist027
		puertorand = ":50053"
	} else { //Servidor Marte
		hostrand = "dist028" //!dist028
		puertorand = ":50054"
	}
	return hostrand, puertorand
}

func main() {
	fmt.Println("Escuchando...")
	rand.Seed(time.Now().UnixNano())             //No queremos la misma seed siempre
	listener, err := net.Listen("tcp", ":50051") // Listener conexion sincrona
	if err != nil {
		panic("La conexion no se pudo crear" + err.Error())
	}
	serv := grpc.NewServer()
	for {
		pb.RegisterMessageServiceServer(serv, &server{})
		if err = serv.Serve(listener); err != nil {
			panic("El server no se pudo iniciar" + err.Error())
		}
	}

}
