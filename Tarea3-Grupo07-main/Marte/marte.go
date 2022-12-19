// Tierra Marte y Titan son iguales! Con la diferencia de que Tierra será nuestro server dominante para el merge.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea3-Grupo07/Proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

var currentPlanet = "Marte"
var relojVectorial = []int32{0, 0, 0} //! PARA PRUEBAS SE DEJARA CON UN 1 EN LA POSICION DEL PLANETA (TIERRA, TITAN, MARTE)

func (s *server) AgregarBase(ctx context.Context, msg *pb.Crear) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base: %s. Valor: %s \n", msg.Sector, msg.Base, msg.Valor)
	escribirtxt("Marte/"+msg.Sector, msg.Sector+" "+msg.Base+" "+msg.Valor)
	escribirlog("Marte", "AgregarBase", msg.Sector, msg.Base, msg.Valor)

	relojVectorial[2] = relojVectorial[2] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) RenombrarBase(ctx context.Context, msg *pb.CambioNombre) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Nombre original: %s. Nombre nuevo: %s \n", msg.Sector, msg.NombreOriginal, msg.NombreNuevo)
	//if existe_base(msg.Sector,msg.Base){} Si existe, hacemos todo, sino existe, respondemos que no se pudo hacer nada xd
	update_base("RenombrarBase", msg.Sector, msg.NombreOriginal, msg.NombreNuevo)
	escribirlog("Marte", "RenombrarBase", msg.Sector, msg.NombreOriginal, msg.NombreNuevo)

	relojVectorial[2] = relojVectorial[2] + 1
	//defer serv.Stop()

	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) ActualizarValor(ctx context.Context, msg *pb.Actualizar) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base: %s. Nuevo valor: %s \n", msg.Sector, msg.Base, msg.NuevoValor)
	//if existe_base(msg.Sector,msg.Base){} Si existe, hacemos todo, sino existe, respondemos que no se pudo hacer nada xd
	update_base("ActualizarValor", msg.Sector, msg.Base, msg.NuevoValor)
	escribirlog("Marte", "ActualizarValor", msg.Sector, msg.Base, msg.NuevoValor)

	relojVectorial[2] = relojVectorial[2] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) BorrarBase(ctx context.Context, msg *pb.Borrar) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base a borrar: %s. \n", msg.Sector, msg.Base)
	//if existe_base(msg.Sector,msg.Base){} Si existe, hacemos todo, sino existe, respondemos que no se pudo hacer nada xd
	escribirlog("Marte", "RenombrarBase", msg.Sector, msg.Base, "")
	update_base("BorrarBase", msg.Sector, msg.Base, "") //! Lo dejo vacio ya que NO nos interesa cuantos soldados hay, solo nos interesa que mueran xd

	relojVectorial[2] = relojVectorial[2] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) SoldadosQuery(ctx context.Context, msg *pb.SoldadosRequest) (*pb.SoldadosAnswer, error) {
	fmt.Printf("Guardianes solicitan a %s los soldados del sector %s, base %s.\n", currentPlanet, msg.Sector, msg.Base)
	fmt.Println("Reloj:", relojVectorial) //! Solo linea de prubea, debe borrarse
	//defer serv.Stop()
	var cantidadSoldados int
	registroEncontrado := false
	baseEncontrada := false

	if _, err := os.Stat("Marte/" + msg.Sector + ".txt"); err == nil {
		//! Si existe entramos aqui

		archivo, err := os.Open("Marte/" + msg.Sector + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer archivo.Close()
		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() { // Recorreremos linea por linea del registro

			decomposed := strings.Split(scanner.Text(), " ") // Separaremos cada linea en sus 3 componentes
			for i := 0; i < len(decomposed); i++ {
				if decomposed[0] == msg.Sector && decomposed[1] == msg.Base {
					cantidadSoldados, err = strconv.Atoi(decomposed[2])
					if err != nil {
						log.Fatal(err)
					}
					baseEncontrada = true
					if cantidadSoldados >= 0 {
						registroEncontrado = true
					}
					break
				}
			}
			if registroEncontrado {
				break
			}
		}
		if !registroEncontrado {
			cantidadSoldados = -1 // Si es que dentro de los registros no se encontro el registro de ese sector y base, retorna -1
			if !baseEncontrada {
				cantidadSoldados = -3
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		//! Si no existe el archivo entramos aqui
		fmt.Printf("El archivo no existe - no se pueden obtener soldados.")

		//! Si no existe retornamos un -1 como cantidad de soldados (?) - Si es -1 tiene que hacer algo despues D:
		return &pb.SoldadosAnswer{CantidadSoldados: int32(-2), PlanetaEmisor: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna

	}

	return &pb.SoldadosAnswer{CantidadSoldados: int32(cantidadSoldados), PlanetaEmisor: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}

func (s *server) MergeRequest(ctx context.Context, msg *pb.Signal) (*pb.MergeAnswerString, error) {
	compendio := ""
	fmt.Printf(" ..:: Llego aviso de merge ::.. \n")
	archivo, err := os.Open("Marte/log_Marte.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {

		compendio = compendio + scanner.Text() + ";"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &pb.MergeAnswerString{Logs: compendio, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}

func (s *server) Consistencia(ctx context.Context, msg *pb.Consistency) (*pb.Ack, error) {

	actualizacion := strings.Split(msg.Actualizar, ";")
	for linea := 0; linea < len(actualizacion); linea++ {

		infoLinea := strings.Split(actualizacion[linea], " ")
		escribirtxt("Marte/"+infoLinea[0], actualizacion[linea])
	}
	relojVectorial = msg.RelojVectorial

	return &pb.Ack{Recibo: 1}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}

/*
Funcion que se encarga escribir en un archivo txt.
Receptor:
-	No hay

Parametros:
-	path: El path donde se va a crear el archivo.
-	la información a escribir en el archivo

Retorno:
-	No hay.
*/
func escribirtxt(path_archivo string, info_registro string) {
	file, err := os.OpenFile(path_archivo+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //!path sera donde se quiere escribir. base.txt o puede ser Tierra/base.txt por ej <- depende como lo armemos

	//Manejo de errores
	if err != nil {
		panic(err)
	}

	data := bufio.NewWriter(file)

	ftemp, _ := file.Stat() //Nos devuelve la estructura del archivo

	defer file.Close() //Defer para que se cierre solo

	//revisamos el archivo! Existe pero no tiene nada escrito:
	if ftemp.Size() == 0 {
		_, err2 := data.WriteString(info_registro + "\n") //Escribe en el archivo
		if err2 != nil {
			panic(err2)
		}

	} else { //Si ya existe y tiene escrito algo!
		_, err2 := data.WriteString(info_registro + "\n")
		if err2 != nil {
			panic(err2)
		}
	}
	data.Flush() //Libera data
	//No hay return porque solo escribe en el archivo!
}

// Este escribe el log (tambien es txt)
func escribirlog(nombre_server_planeta string, nombre_accion string, nombre_sector string, nombre_base string, data_extra string) {

	escribirtxt("Marte/log_"+nombre_server_planeta, nombre_accion+" "+nombre_sector+" "+nombre_base+" "+data_extra) //! data_extra puede ser el nuevo nombre de la base, un valor extra de soldados, o vacio sin nada (depende de la accion que se vaya a hacer)

}

//Revisar si existe la base!
/*
func existe_base(nombre_sector string, nombre_base string) bool {
	f, err := os.Open("Testing/" + nombre_sector + ".txt")
	if err != nil {
		return false
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), nombre_base) { //!Se revisa el texto hasta encontrar coincidencia.
			lineaSeparada := strings.Split(scanner.Text(), " ")
			fmt.Printf("[Buscando base %s] La base %s fue encontrada con el valor %s.\n", nombre_base, lineaSeparada[1], lineaSeparada[2])
			return true
		}
	}

	return false
}
*/

/*
Funcion que se encarga de Updatear La base dependiendo de la acción correspondiente
Receptor:
-	No hay

Parametros:
-	accion a realizar
-	nombre del sector
-	nombre de la base
-	nuevo valor a cambiar (Nuevo Nombre base, nuevo valor o borrar base)

Retorno:
-	No hay.
*/

func update_base(accion string, nombre_sector string, nombre_base string, nuevo_valor string) bool { //! Le agregue la accion para ver caso x caso c:
	var updateExitoso bool
	//* Abrimos el archivo del sector
	f, err := os.Open("Marte/" + nombre_sector + ".txt")
	if err != nil {
		return false
	}
	//* Terminamos de abrir el archivo del sector

	//* Comienza creación de archivo de actualizacion
	//! OJO que el nombre del archivo de actualizacion es el mismo que el original solo que tiene un 2, revisar en siguiente linea de codigo
	file, err := os.OpenFile("Marte/"+nombre_sector+"V2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //path sera donde se quiere escribir. base.txt o puede ser Tierra/base.txt por ej <- depende como lo armemos
	if err != nil {
		panic(err)
	}
	data := bufio.NewWriter(file)
	//* Terminamos de crear el archivo de actualizacion

	//*Recorreremos el archivo original linea por linea, buscando la linea a actualizar, las lineas que no se actualizan, se copian tal cual al archivo nuevo, la linea a actualizar, se rescata sector y base y se le concatena el valor nuevo para luego ser escrita en el archivo nuevo
	scanner := bufio.NewScanner(f)
	for scanner.Scan() { // Iremos linea por linea

		if accion == "ActualizarValor" {
			var info_registro string                           // variable string que contiene linea a escribir en archivo actualizado
			if strings.Contains(scanner.Text(), nombre_base) { // Se revisa si es que la linea contiene la base buscada.
				lineaSeparada := strings.Split(scanner.Text(), " ") // Separamos el string, sabemos que el tercer elemento es el que descartaremos [Sector, Base, Valor]
				info_registro = lineaSeparada[0] + " " + lineaSeparada[1] + " " + nuevo_valor
				updateExitoso = true //? Si es que todo hasta aca salio bien, marcamos como update exitoso
			} else {
				info_registro = scanner.Text() // Se copia la linea tal cual
			}

			//Si es que el archivo ya existe y tiene escrito algo!
			_, err2 := data.WriteString(info_registro + "\n")
			if err2 != nil {
				panic(err2)

			}
			data.Flush() //Libera data

		} else if accion == "RenombrarBase" {
			var info_registro string                           // variable string que contiene linea a escribir en archivo actualizado
			if strings.Contains(scanner.Text(), nombre_base) { // Se revisa si es que la linea contiene la base buscada.
				lineaSeparada := strings.Split(scanner.Text(), " ") // Separamos el string, sabemos que el segundo elemento es el que descartaremos [Sector, Base, Valor]
				info_registro = lineaSeparada[0] + " " + nuevo_valor + " " + lineaSeparada[2]
				updateExitoso = true //? Si es que todo hasta aca salio bien, marcamos como update exitoso
			} else {
				info_registro = scanner.Text() // Se copia la linea tal cual
			}

			//Si es que el archivo ya existe y tiene escrito algo!
			_, err2 := data.WriteString(info_registro + "\n")
			if err2 != nil {
				panic(err2)

			}
			data.Flush() //Libera data

		} else { //! BorrarBase
			var info_registro string                           // variable string que contiene linea a escribir en archivo actualizado
			if strings.Contains(scanner.Text(), nombre_base) { // Se revisa si es que la linea contiene la base buscada.
				updateExitoso = true //!Si entramos aqui, significa que vamos a ignorar la existencia de la base, es decir, no lo escribiremos en el archivo y cambiamos a true
			} else {
				info_registro = scanner.Text()                    // Se copia la linea tal cual
				_, err2 := data.WriteString(info_registro + "\n") //! Solo creo, que en este escenario, debiese escribirse dentro del else.
				//! Ya que si lo escribimos afuera, la base a eliminar, será escrita nuevamente.
				if err2 != nil {
					panic(err2)

				}
			}

			data.Flush() //Libera data
		}

	}
	//* En este punto el archivo nuevo tiene toda la info actualizada, falta eliminar el archivo original y renombrar el nuevo con el nombre del original (quitarle el 2)

	f.Close()    // Cerramos archivo para poder eliminarlo
	file.Close() // Cerramos archivo actualizado para poder cambiarle el nombre
	change_namearch(nombre_sector)
	fmt.Println("La base ha sido cambiada/eliminada de manera exitosa :D")
	return updateExitoso
}

// Esta funcion solamente elimina el archivo antiguo y renombra el archivo nuevo.
func change_namearch(nombre_sector string) {
	eliminarArchivoOriginal := os.Remove("Marte/" + nombre_sector + ".txt") //! Esto lo elimina
	if eliminarArchivoOriginal != nil {
		log.Fatal(eliminarArchivoOriginal)
	}

	cambioNombre_ArchivoActualizado := os.Rename("Marte/"+nombre_sector+"V2.txt", "pruebafunc/"+nombre_sector+".txt") //! Este le actualiza el nombre.
	if cambioNombre_ArchivoActualizado != nil {
		log.Fatal(cambioNombre_ArchivoActualizado)
	}
}

func main() {
	fmt.Println("Escuchando...")
	listener, err := net.Listen("tcp", ":50054") // Listener conexion sincrona
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
