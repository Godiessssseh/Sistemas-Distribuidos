/*
No debería quedar tan largo.
Guardianes
-----------------------------
->	Monotonic Reads
->	Obtener información de los soldados en algún planeta.
->	GetSoldados (name_sector, name_base)
*/

package main

import (
	"context"
	"fmt"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea3-Grupo07/Proto"
	"google.golang.org/grpc"
)

type register struct {
	nombreSector   string
	planeta        string
	relojVector    []int32
	conteoSoldados int
}

var history []register //? Aca se guardaran las consultas anteriores

func IsNewerThan(relojRegistro []int32, relojRespuesta []int32, posicionPlaneta int) bool {
	primerPlaneta := relojRegistro[0] < relojRespuesta[0]  // Si es que el vector recibido indica ser mas nuevo respecto a la Tierra
	segundoPlaneta := relojRegistro[1] < relojRespuesta[1] // Si es que el vector recibido indica ser mas nuevo respecto a Titan
	tercerPlaneta := relojRegistro[2] < relojRespuesta[2]  // Si es que el vector recibido indica ser mas nuevo respecto a Marte
	if primerPlaneta && segundoPlaneta && tercerPlaneta {  // Si es que el vector es mas reciente en los tres planetas
		return true // Dada la condicion cierta, lo aceptamos
	} else {
		return false
	}
}
func GetSoldados(serviceCliente pb.MessageServiceClient, sector string, base string) int32 {
	res, err := serviceCliente.GetSoldados(context.Background(),
		&pb.SoldadosRequest{
			Sector: sector,
			Base:   base,
		})

	if err != nil {
		panic("No se puede crear el mensaje " + err.Error())
	}
	fmt.Printf("Recibiendo soldados desde el planeta %s....\n", res.PlanetaEmisor)

	if res.CantidadSoldados == -1 {
		// el archivo existe pero no hay informacion de los soldados
		fmt.Printf("Existe el archivo de registros del sector %s, en la base %s, pero no hay informacion sobre sus soldados. Codigo de error: %d.\n", sector, base, res.CantidadSoldados)
	} else if res.CantidadSoldados == -2 {
		fmt.Printf("No existe un archivo de registros del sector %s. Codigo de error: %d.\n", sector, res.CantidadSoldados)
	} else if res.CantidadSoldados == -3 {
		fmt.Printf("Existe un archivo de registros del sector %s, pero no de la base %s. Codigo de error: %d.\n", sector, base, res.CantidadSoldados)
	} else {
		fmt.Println("Soldados recibidos desde : ", res.PlanetaEmisor, "El reloj recibido es: ", res.RelojVectorial)

		var existeRegistro bool = false
		if len(history) != 0 { // Si es que el registro no esta vacio
			for indice, registro := range history { // iteramos sobre los registros
				if registro.nombreSector == sector { // Si es que hay registro existente del planeta que se esta consultando
					var posicionPlaneta int
					switch res.PlanetaEmisor {
					case "Tierra":
						posicionPlaneta = 0
					case "Titan":
						posicionPlaneta = 1
					case "Marte":
						posicionPlaneta = 2
					}
					fmt.Println("Reloj que ya teniamos: ", registro.relojVector)
					fmt.Printf("Hay registro del planeta %s, sobre el sector %s.\n", registro.planeta, registro.nombreSector)
					fmt.Println("Reloj que nos llego:", res.RelojVectorial)
					if IsNewerThan(registro.relojVector, res.RelojVectorial, posicionPlaneta) { // Revisamos si es que el registro que estamos recibiendo es mas nuevo respecto al que ya teniamos
						history[indice].planeta = res.PlanetaEmisor
						history[indice].conteoSoldados = int(res.CantidadSoldados)
						history[indice].relojVector = res.RelojVectorial
						existeRegistro = true
						fmt.Println("Actualizacion exitosa, ahora el reloj es:", history[indice].relojVector, ".")
						break
					} else {
						existeRegistro = true
					}
				}
			}
		}
		if !existeRegistro { // Si es que no hay registro, se debe crear uno nuevo
			var nuevoRegistro = register{sector, res.PlanetaEmisor, res.RelojVectorial, int(res.CantidadSoldados)}
			history = append(history, nuevoRegistro)
			fmt.Printf("Registro exitoso.\n")
		}

	}
	return res.CantidadSoldados
}

func CrearConexion(direccion string, puerto string) pb.MessageServiceClient {
	hostS := direccion                                         //Host de un Laboratorio
	connS, err := grpc.Dial(hostS+puerto, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
	//! connS2, err2 := grpc.Dial(hostS+puerto, grpc.WithTransportCredentials(insecure.NewCredentials())) creo que este es el workaround para quitar el warning.
	if err != nil {
		panic("No se pudo conectar con el servidor" + err.Error())
	}
	// defer connS.Close() // Se "aplaza" el cierre de la conexion
	serviceCliente := pb.NewMessageServiceClient(connS)
	return serviceCliente

}

func main() {

	fmt.Println("--- Bienvenido ---")
	canalMensajes := CrearConexion("dist025", ":50051")
	var sector, base string
	flag := true

	for flag {

		fmt.Println("Ingrese - Nombre Sector -")
		fmt.Scanln(&sector)

		fmt.Println("Ingrese - Nombre Base -")
		fmt.Scanln(&base)
		//?sector := "Chile"
		//?base := "Santiago"
		cantidadSoldados := GetSoldados(canalMensajes, sector, base)
		fmt.Printf("En el sector de %s, en la base %s, hay %d soldados. \n", sector, base, cantidadSoldados)
	}

}
