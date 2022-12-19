/*No debería quedar tan largo.
Vanguardia
-----------------------------
->	Read your write
->	AgregarBase (name_sector, name_base)[valor] -> Utiliza nombre base para crearlo
->	RenombrarBase (name_sector, name_base, nuevo_nombre_base) -> cambiar el nombre de la base por uno nuevo
->	ActualizarValor (name_sector, name_base, valor) -> actualizar soldados en la base correspondiente.
->	BorrarBase (name_sector, name_base) -> Borra el txt de la base + borrarse de los logs (ta dura esta)
*/

package main

import (
	"context"
	"fmt"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea3-Grupo07/Proto"
	"google.golang.org/grpc"
)

//Esta comentada porque no estoy seguro si funciona, estoy analizando aún!

// Usar para read your writes! (o eso espero)
type Read struct {
	ultimo_sector     string
	ultima_base       string
	ultimo_movimiento string
	relojVector       []int32
	ultimo_servidor   string
}

var contador int = 1 //Se utiliza para crear el archivo!
var leer []Read      //? Guardar aquí las consultas a read your writes!

//? Como hicimos funcs separadas para cada opcion (agregar, eliminar, actualizarvalor, cambiarnombre)
//? Mejor una funcion global que sea invocada por c/u para que comparemos el último agregado!

func Add(sector string, base string, movimiento string, reloj []int32, servidor string) {
	var existeLeer bool = false
	if len(leer) != 0 { // Si es que el registro no esta vacio
		if leer[0].ultimo_servidor != servidor { //! Revisamos si el server que usaremos (?) es el mismo que ya tenemos guardado para hacer todo más rápido?
			//Como solo guarda el ultimo movimiento, hay que cambiar valores solamente
			fmt.Println("Último sector registrado", leer[0].ultimo_sector)
			fmt.Printf("Última base registrada %s\n", leer[0].ultima_base)
			fmt.Println("Último movimiento registrado", leer[0].ultimo_movimiento)
			fmt.Println("Último reloj registrado", leer[0].relojVector)
			fmt.Printf("Último registro del servidor planeta %s\n", leer[0].ultimo_servidor)

			leer[0].ultimo_sector = sector
			leer[0].ultima_base = base
			leer[0].ultimo_movimiento = movimiento
			leer[0].relojVector = reloj
			leer[0].ultimo_servidor = servidor
			fmt.Println("Actualizacion al registro realizada.")
			fmt.Println("")
			existeLeer = true
		} else {
			//! Si los servidores son iguales, esto implica que solo actualizamos el reloj y el movimiento (utilizamos actualizar base para fixear un error xd )
			leer[0].ultima_base = base
			leer[0].ultimo_movimiento = movimiento
			leer[0].relojVector = reloj
			fmt.Println("Movimiento actualizado.", leer[0].ultimo_movimiento)
			fmt.Println("Reloj actualizado. \n", leer[0].relojVector)
			existeLeer = true
		}
	}
	if !existeLeer {
		var nuevoLeer = Read{sector, base, movimiento, reloj, servidor} //!Check a esto!
		leer = append(leer, nuevoLeer)
		fmt.Printf("Registro creado por primera vez de manera exitosa.\n")
	}
}

// Para read your writes!
func gethost(planeta string) (host string, puerto string) {
	hostrand := ""
	puertorand := ""
	if planeta == "Tierra" {
		hostrand = "localhost" //!dist026 Cambiar despues!!
		puertorand = ":50052"
	} else if planeta == "Marte" {
		hostrand = "localhost" //!dist028
		puertorand = ":50054"
	} else {
		hostrand = "localhost" //!dist027
		puertorand = ":50053"
	}
	return hostrand, puertorand
}

func GetPlanet(serviceCliente pb.MessageServiceClient, tipoSolicitud string) (direccion string, puerto string) {
	res, err := serviceCliente.GetPlanet(context.Background(),
		&pb.SolicitudPlaneta{
			Solicitud: tipoSolicitud,
		})

	if err != nil {
		panic("No se puede crear el mensaje " + err.Error())
	}
	direccion = res.DireccionAsignada
	puerto = res.PuertoAsignado
	return
}

/*	AgregarBase recibe la conexion, Chile, Santiago, 0 (formato de obtener valores)
*	invoca a get planeta -> broker usara getplaneta y retornara una direccion y un puerto al azar (planeta a conectar)
*	con el host y puerto obtenidos, creamos la conexion
*	despues de la conexion, se trabaja con AgregarBase!
 */
func AgregarBase(serviceCliente pb.MessageServiceClient, sector string, base string, valor string) int32 {
	if leer[0].ultimo_sector == sector && leer[0].ultima_base == base {
		//? Usamos read your writes, si el sector y la base es la misma, nos conectamos directamente a ellos!
		direccionPlaneta, puertoPlaneta := gethost(leer[0].ultimo_servidor) //! No sé si se puede invocar así :c
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)   //Direccion planeta + puerto planeta!
		res, err := canalConPlaneta.AgregarBase(context.Background(),
			&pb.Crear{
				Sector: sector,
				Base:   base,
				Valor:  valor,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		Add(sector, base, "AgregarBase", res.RelojVectorial, leer[0].ultimo_servidor) //
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		fmt.Println("Base agregada en: ", res.Planeta)
		return res.Codigo

	} else {
		direccionPlaneta, puertoPlaneta := GetPlanet(serviceCliente, "AgregarBase")
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.AgregarBase(context.Background(),
			&pb.Crear{
				Sector: sector,
				Base:   base,
				Valor:  valor,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, base, "AgregarBase", res.RelojVectorial, res.Planeta) //En que base fue agregada la información!
		fmt.Println("Base agregada en: ", res.Planeta)

		return res.Codigo
	}
}

/*	RenombrarBase recibe la conexion, Chile, Santiago, 0 (formato de obtener valores)
*	invoca a get planeta -> broker usara getplaneta y retornara una direccion y un puerto al azar (planeta a conectar)
*	con el host y puerto obtenidos, creamos la conexion
*	despues de la conexion, se trabaja con RenombrarBase!
 */

func RenombrarBase(serviceCliente pb.MessageServiceClient, sector string, nombreOriginal string, nombreNuevo string) int32 {
	if leer[0].ultimo_sector == sector && leer[0].ultima_base == nombreOriginal {
		//? Usamos read your writes, si el sector y la base es la misma, nos conectamos directamente a ellos!
		direccionPlaneta, puertoPlaneta := gethost(leer[0].ultimo_servidor) //! No sé si se puede invocar así :c
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.RenombrarBase(context.Background(),
			&pb.CambioNombre{
				Sector:         sector,
				NombreOriginal: nombreOriginal,
				NombreNuevo:    nombreNuevo,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, nombreNuevo, "RenombrarBase", res.RelojVectorial, leer[0].ultimo_servidor) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base renombrada en: ", res.Planeta)
		return res.Codigo

	} else {
		direccionPlaneta, puertoPlaneta := GetPlanet(serviceCliente, "RenombrarBase")
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.RenombrarBase(context.Background(),
			&pb.CambioNombre{
				Sector:         sector,
				NombreOriginal: nombreOriginal,
				NombreNuevo:    nombreNuevo,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, nombreNuevo, "RenombrarBase", res.RelojVectorial, res.Planeta) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base renombrada en: ", res.Planeta)
		return res.Codigo
	}
}

/*	ActualizarValor recibe la conexion, Chile, Santiago, 0 (formato de obtener valores)
*	invoca a get planeta -> broker usara getplaneta y retornara una direccion y un puerto al azar (planeta a conectar)
*	con el host y puerto obtenidos, creamos la conexion
*	despues de la conexion, se trabaja con ActualizarValor!
 */

func ActualizarValor(serviceCliente pb.MessageServiceClient, sector string, base string, nuevoValor string) int32 {
	if leer[0].ultimo_sector == sector && leer[0].ultima_base == base {
		direccionPlaneta, puertoPlaneta := gethost(leer[0].ultimo_servidor) //! No sé si se puede invocar así :c
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.ActualizarValor(context.Background(),
			&pb.Actualizar{
				Sector:     sector,
				Base:       base,
				NuevoValor: nuevoValor,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, base, "ActualizarValor", res.RelojVectorial, leer[0].ultimo_servidor) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base actualizada en: ", res.Planeta)
		return res.Codigo

	} else {
		direccionPlaneta, puertoPlaneta := GetPlanet(serviceCliente, "ActualizarValor")
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.ActualizarValor(context.Background(),
			&pb.Actualizar{
				Sector:     sector,
				Base:       base,
				NuevoValor: nuevoValor,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, base, "ActualizarValor", res.RelojVectorial, res.Planeta) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base actualizada en: ", res.Planeta)
		return res.Codigo
	}

}

/*	BorrarBase recibe la conexion, Chile, Santiago, 0 (formato de obtener valores)
*	invoca a get planeta -> broker usara getplaneta y retornara una direccion y un puerto al azar (planeta a conectar)
*	con el host y puerto obtenidos, creamos la conexion
*	despues de la conexion, se trabaja con BorrarBase!
 */

func BorrarBase(serviceCliente pb.MessageServiceClient, sector string, base string) int32 {
	if leer[0].ultimo_sector == sector && leer[0].ultima_base == base {
		direccionPlaneta, puertoPlaneta := gethost(leer[0].ultimo_servidor) //! No sé si se puede invocar así :c
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.BorrarBase(context.Background(),
			&pb.Borrar{
				Sector: sector,
				Base:   base,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, base, "BorrarBase", res.RelojVectorial, leer[0].ultimo_servidor) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base borrada en: ", res.Planeta)
		return res.Codigo

	} else {
		direccionPlaneta, puertoPlaneta := GetPlanet(serviceCliente, "BorrarBase")
		canalConPlaneta := CrearConexion(direccionPlaneta, puertoPlaneta)
		res, err := canalConPlaneta.BorrarBase(context.Background(),
			&pb.Borrar{
				Sector: sector,
				Base:   base,
			})

		if err != nil {
			panic("No se puede crear el mensaje " + err.Error())
		}
		//println("Respuesta: %d \n", res.Codigo)//respuesta del laboratorio
		Add(sector, base, "BorrarBase", res.RelojVectorial, res.Planeta) //Se guarda nombre nuevo debido a que es la NUEVA base a guardar!
		fmt.Println("Base borrada en: ", res.Planeta)
		return res.Codigo
	}
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

// Esta func hace el trabajo de obtener los inputs por consola
// Accion -> Agregar, Renombrar, Actualizar, Borrar las bases.
// nombre del sector y base, con la cantidad de soldados ingresados x pantalla.s
func interfaz() []string {
	var eleccion, accion, sector, base, soldiers, newname string
	var prueba []string //Se guarda la accion + el sector + la base + cantidad
	//? Lista prueba -> [accion, sector, base, valor]
	flag := true

	for flag {
		fmt.Println("Ingrese el número del comando a utilizar")
		fmt.Println("-1- AgregarBase\n-2- RenombrarBase\n-3- ActualizarValor\n-4- BorrarBase")
		fmt.Scanln(&eleccion)

		if eleccion == "1" { //? AgregarBase -> (sector, base, soldados)
			fmt.Println("Ingrese - Nombre Sector -")
			fmt.Scanln(&sector)

			fmt.Println("Ingrese - Nombre Base -")
			fmt.Scanln(&base)

			fmt.Println("Ingrese - Cantidad Soldados -")
			fmt.Scanln(&soldiers)

			accion = "AgregarBase"
			if soldiers == "" { //? Si no ingresan cantidad, automatico es 0
				soldiers = "0"
			}

			prueba = append(prueba, accion, sector, base, soldiers)
			flag = false

		} else if eleccion == "2" { //? RenombrarBase -> (sector, base, soldados)
			fmt.Println("Ingrese - Nombre Sector -")
			fmt.Scanln(&sector)

			fmt.Println("Ingrese - Nombre Base -")
			fmt.Scanln(&base)

			fmt.Println("Ingrese - Nuevo nombre de Base -")
			fmt.Scanln(&newname)

			accion = "RenombrarBase"
			prueba = append(prueba, accion, sector, base, newname) //?AgregarBase
			flag = false

		} else if eleccion == "3" { //? ActualizarValor -> (sector, base, newvalor)
			fmt.Println("Ingrese - Nombre Sector -")
			fmt.Scanln(&sector)

			fmt.Println("Ingrese - Nombre Base -")
			fmt.Scanln(&base)

			fmt.Println("Ingrese - Nueva cantidad de soldados -")
			fmt.Scanln(&soldiers)

			accion = "ActualizarValor"
			if soldiers == "" { //? Si no ingresan cantidad, automatico es 0
				soldiers = "0"
			}

			prueba = append(prueba, accion, sector, base, soldiers)
			flag = false

		} else if eleccion == "4" {
			fmt.Println("Ingrese -Nombre Sector-")
			fmt.Scanln(&sector)

			fmt.Println("Ingrese -Nombre Base -")
			fmt.Scanln(&base)

			accion = "BorrarBase"
			prueba = append(prueba, accion, sector, base)
			flag = false

		} else {
			fmt.Println("No elegiste de manera correcta :c sadge")
		}
	}
	return prueba //return a la lista!
}

func main() {
	fmt.Println("-------- Bienvenidx a la Vanguardia --------")
	canalMensajes := CrearConexion("dist025", ":50051")
	flag := true
	for flag {
		ingreso := interfaz() //<- Esta func obtiene una lista: [0] -> accion, [1] -> sector, [2] -> base , [3] -> cantidad_soldados
		if contador == 1 {    //Pensando donde me meto esto xd
			Add(ingreso[0], ingreso[1], "prueba", []int32{0, 0, 0}, "prueba") //!Check a esto!
			contador = 0
		}
		if ingreso[0] == "AgregarBase" {
			if AgregarBase(canalMensajes, ingreso[1], ingreso[2], ingreso[3]) == 1 { //! Please check this !! Quizas cambiarlo por string?
				println("Base agregada con exito. \n")
			} else {
				println("La base no se agrego correctamente. D: \n")
			}
		} else if ingreso[0] == "RenombrarBase" {
			if RenombrarBase(canalMensajes, ingreso[1], ingreso[2], ingreso[3]) == 1 {
				println("La base ha sido renombrada con exito. \n")
			} else {
				println("La base no pudo ser renombrada correctamente. D: \n")
			}
		} else if ingreso[0] == "ActualizarValor" {
			if ActualizarValor(canalMensajes, ingreso[1], ingreso[2], ingreso[3]) == 1 { //! Same problem que arriba ojito
				println("Valor actualizado con exito. \n")
			} else {
				println("El valor no pudo ser agregado correctamente. T_T \n")
			}
		} else { //!BorrarBase
			if BorrarBase(canalMensajes, ingreso[1], ingreso[2]) == 1 {
				println("Base eliminada con exito. \n")
			} else {
				println("La base no se elimino correctamente. T_T \n")
			}
		}
	}

}
