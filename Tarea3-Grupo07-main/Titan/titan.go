// Tierra Marte y Titan son iguales! Con la diferencia de que Tierra será nuestro server dominante para el merge.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea3-Grupo07/Proto"
	"google.golang.org/grpc"
)

type register struct {
	nombreSector   string
	planeta        string
	base           string
	relojVector    []int32
	conteoSoldados int
}

var history []register //? Aca se guardaran las consultas anteriores
type server struct {
	pb.UnimplementedMessageServiceServer
}

var currentPlanet = "Titan"
var relojVectorial = []int32{0, 0, 0} //! PARA PRUEBAS SE DEJARA CON UN 1 EN LA POSICION DEL PLANETA (TIERRA, TITAN, MARTE)

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

func (s *server) AgregarBase(ctx context.Context, msg *pb.Crear) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base: %s. Valor: %s \n", msg.Sector, msg.Base, msg.Valor)
	escribirtxt("Titan/"+msg.Sector, msg.Sector+" "+msg.Base+" "+msg.Valor)
	escribirlog("Titan", "AgregarBase", msg.Sector, msg.Base, msg.Valor)
	relojVectorial[1] = relojVectorial[1] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) RenombrarBase(ctx context.Context, msg *pb.CambioNombre) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Nombre original: %s. Nombre nuevo: %s \n", msg.Sector, msg.NombreOriginal, msg.NombreNuevo)
	relojVectorial[1] = relojVectorial[1] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) ActualizarValor(ctx context.Context, msg *pb.Actualizar) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base: %s. Nuevo valor: %s \n", msg.Sector, msg.Base, msg.NuevoValor)
	relojVectorial[1] = relojVectorial[1] + 1
	//defer serv.Stop()
	return &pb.CodigoRespuesta{Codigo: 1, Planeta: currentPlanet, RelojVectorial: relojVectorial}, nil //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}
func (s *server) BorrarBase(ctx context.Context, msg *pb.Borrar) (*pb.CodigoRespuesta, error) {
	fmt.Printf("Sector: %s. Base a borrar: %s. \n", msg.Sector, msg.Base)
	relojVectorial[1] = relojVectorial[1] + 1
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

	if _, err := os.Stat("Titan/" + msg.Sector + ".txt"); err == nil {
		//! Si existe entramos aqui

		archivo, err := os.Open("Titan/" + msg.Sector + ".txt")
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

func MergeRequest(serviceCliente pb.MessageServiceClient, aviso int32) (logs string, reloj []int32) {
	fmt.Printf(" ..:: Se genero solicitud de merge, al planeta %d ::.. ", aviso)
	res, err := serviceCliente.MergeRequest(context.Background(),
		&pb.Signal{
			Aviso: aviso,
		})

	if err != nil {
		panic("No se puede crear el mensaje con el planeta " + strconv.Itoa(int(aviso)) + err.Error())
	}
	logs = res.Logs
	reloj = res.RelojVectorial
	return //? nil== valor no inicializado, representa "no hay error", en este caso el retorno error==nil, por lo que todo estara bien si es que es eso lo que se retorna
}

func Consistencia(serviceCliente pb.MessageServiceClient, consistente string, reloj []int32) (confirmacion int32) {
	res, err := serviceCliente.Consistencia(context.Background(),
		&pb.Consistency{
			Actualizar: consistente,
		})
	if err != nil {
		panic("No se puede crear el mensaje con el planeta " + err.Error())
	}
	return res.Recibo
}

/*

	REVISAR SI EXISTE EL SECTOR O SI EXISTE LA BASE

*/

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

	escribirtxt("Titan/log_"+nombre_server_planeta, nombre_accion+" "+nombre_sector+" "+nombre_base+" "+data_extra) //! data_extra puede ser el nuevo nombre de la base, un valor extra de soldados, o vacio sin nada (depende de la accion que se vaya a hacer)

}

func ImportarLogs(canalMensajesMarte pb.MessageServiceClient, canalMensajesTierra pb.MessageServiceClient) (relojTierra []int32, numRegistrosTierra int, relojMarte []int32, numRegistrosMarte int) {
	var stringTierra string
	var stringMarte string
	stringTierra, relojTierra = MergeRequest(canalMensajesTierra, 1)
	stringMarte, relojMarte = MergeRequest(canalMensajesMarte, 1)
	numRegistrosTierra = strings.Count(stringTierra, ";")
	numRegistrosMarte = strings.Count(stringMarte, ";")
	compiladoTierra := strings.Split(stringTierra, ";")                                                 // Dividimos el string, de tal manera que quede un arreglo donde cada elemento corresponde a una linea del log original
	compiladoMarte := strings.Split(stringMarte, ";")                                                   // Dividimos el string, de tal manera que quede un arreglo donde cada elemento corresponde a una linea del log original
	for i := 0; i+1 < int(math.Max(float64(len(compiladoTierra)), float64(len(compiladoMarte)))); i++ { // Escribiremos los logs, recorriendo los strings con limite en el tamaño del string mas grande
		if i+1 < int(math.Min(float64(len(compiladoTierra)), float64(len(compiladoMarte)))) { // Para evitar Segmentation Fault, recorreremos el string pequeño hasta su tamaño
			if len(compiladoTierra) >= len(compiladoMarte) {
				escribirtxt("Titan/log_Tierra", compiladoTierra[i])
				escribirtxt("Titan/log_Marte", compiladoMarte[i])
			} else {
				escribirtxt("Titan/log_Marte", compiladoMarte[i])
				escribirtxt("Titan/log_Tierra", compiladoTierra[i])
			}

		} else {
			if len(compiladoTierra) >= len(compiladoMarte) {
				escribirtxt("Titan/log_Tierra", compiladoTierra[i])
			} else {
				escribirtxt("Titan/log_Marte", compiladoMarte[i])
			}
		}
	}
	return
}

func deleteRegister(indice int) {
	// delete an element from the array
	newLength := 0
	for index := range history {
		if indice != index {
			history[newLength] = history[index]
			newLength++
		}
	}

	// reslice the array to remove extra index
	newArray := history[:newLength]
	history = newArray
}

func GenerarConsistencia(relojTierra []int32, numRegistrosTierra int, relojMarte []int32, numRegistrosMarte int) {

	fmt.Printf(" ..:: Comienza proceso de merge ::.. \n")
	archivoTierra, err := os.Open("Titan/log_Tierra.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer archivoTierra.Close()

	archivoTitan, err := os.Open("Titan/log_Titan.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer archivoTitan.Close()

	archivoMarte, err := os.Open("Titan/log_Marte.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer archivoMarte.Close()

	scannerTierra := bufio.NewScanner(archivoTierra)
	scannerTitan := bufio.NewScanner(archivoTitan)
	scannerMarte := bufio.NewScanner(archivoMarte)
	for scannerTierra.Scan() {
		var existeRegistro bool = false
		if len(history) != 0 { // Si es que el registro no esta vacio
			for indice, registro := range history { // iteramos sobre los registros
				lineaLeida := strings.Split((scannerTierra.Text()), " ")
				if registro.nombreSector == lineaLeida[1] && registro.base == lineaLeida[2] { // Si es que hay registro existente del planeta que se esta consultando
					existeRegistro = true
					deleteRegister(indice)

					//Y si revisamos los archivos, encontramos un problema -> Lo borramos yera? -> Iñaki dice que puede ser un workaround bueno uwu
					//Si no es eso no se que mas puede ser :c

				}
			}
		}
		if !existeRegistro { // Si es que no hay registro, se debe crear uno nuevo
			lineaLeida := strings.Split((scannerTierra.Text()), " ")
			if lineaLeida[0] == "AgregarBase" {
				newvalue, _ := strconv.Atoi(lineaLeida[3])
				var nuevoRegistro = register{lineaLeida[1], "Tierra", lineaLeida[2], relojTierra, newvalue}
				history = append(history, nuevoRegistro)
			} else if lineaLeida[0] == "RenombrarBase" {
				// Lo dejaremos pasar
			} else if lineaLeida[0] == "ActualizarValor" {
				// Lo dejaremos pasar
			} else {
				// De todos modos no existia
			}
		}
	}
	if err := scannerTierra.Err(); err != nil {
		log.Fatal(err)
	}

	for scannerTitan.Scan() {
		var existeRegistro bool = false
		if len(history) != 0 { // Si es que el registro no esta vacio
			for indice, registro := range history { // iteramos sobre los registros
				lineaLeida := strings.Split((scannerTitan.Text()), " ")
				if registro.nombreSector == lineaLeida[1] && registro.base == lineaLeida[2] { // Si es que hay registro existente del planeta que se esta consultando
					existeRegistro = true
					deleteRegister(indice)

					//Y si revisamos los archivos, encontramos un problema -> Lo borramos yera? -> Iñaki dice que puede ser un workaround bueno uwu
					//Si no es eso no se que mas puede ser :c

				}
			}
		}
		if !existeRegistro { // Si es que no hay registro, se debe crear uno nuevo
			lineaLeida := strings.Split((scannerTitan.Text()), " ")
			if lineaLeida[0] == "AgregarBase" {
				newvalue, _ := strconv.Atoi(lineaLeida[3])
				var nuevoRegistro = register{lineaLeida[1], "Tierra", lineaLeida[2], relojTierra, newvalue}
				history = append(history, nuevoRegistro)
			} else if lineaLeida[0] == "RenombrarBase" {
				// Lo dejaremos pasar
			} else if lineaLeida[0] == "ActualizarValor" {
				// Lo dejaremos pasar
			} else {
				// De todos modos no existia
			}
		}
	}
	if err := scannerTitan.Err(); err != nil {
		log.Fatal(err)
	}

	for scannerMarte.Scan() {
		var existeRegistro bool = false
		if len(history) != 0 { // Si es que el registro no esta vacio
			for indice, registro := range history { // iteramos sobre los registros
				lineaLeida := strings.Split((scannerMarte.Text()), " ")
				if registro.nombreSector == lineaLeida[1] && registro.base == lineaLeida[2] { // Si es que hay registro existente del planeta que se esta consultando
					existeRegistro = true
					deleteRegister(indice)

					//Y si revisamos los archivos, encontramos un problema -> Lo borramos yera? -> Iñaki dice que puede ser un workaround bueno uwu
					//Si no es eso no se que mas puede ser :c

				}
			}
		}
		if !existeRegistro { // Si es que no hay registro, se debe crear uno nuevo
			lineaLeida := strings.Split((scannerMarte.Text()), " ")
			if lineaLeida[0] == "AgregarBase" {
				newvalue, _ := strconv.Atoi(lineaLeida[3])
				var nuevoRegistro = register{lineaLeida[1], "Tierra", lineaLeida[2], relojTierra, newvalue}
				history = append(history, nuevoRegistro)
			} else if lineaLeida[0] == "RenombrarBase" {
				// Lo dejaremos pasar
			} else if lineaLeida[0] == "ActualizarValor" {
				// Lo dejaremos pasar
			} else {
				// De todos modos no existia
			}
		}
	}
	if err := scannerMarte.Err(); err != nil {
		log.Fatal(err)
	}
	e := os.Remove("Titan/log_Tierra.txt")
	if e != nil {
		log.Fatal(e)
	}
	file, err := os.OpenFile("Titan/log_Tierra.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//Manejo de errores
	if err != nil {
		panic(err)
	}
	_ = bufio.NewWriter(file)
	_, _ = file.Stat() //Nos devuelve la estructura del archivo
	defer file.Close() //Defer para que se cierre solo

	e = os.Remove("Titan/log_Titan.txt")
	if e != nil {
		log.Fatal(e)
	}
	file, err = os.OpenFile("Titan/log_Titan.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//Manejo de errores
	if err != nil {
		panic(err)
	}
	_ = bufio.NewWriter(file)
	_, _ = file.Stat() //Nos devuelve la estructura del archivo
	defer file.Close() //Defer para que se cierre solo

	e = os.Remove("Titan/log_Marte.txt")
	if e != nil {
		log.Fatal(e)
	}
	file, err = os.OpenFile("Titan/log_Marte.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//Manejo de errores
	if err != nil {
		panic(err)
	}
	_ = bufio.NewWriter(file)
	_, _ = file.Stat() //Nos devuelve la estructura del archivo
	defer file.Close() //Defer para que se cierre solo
	compendio := ""
	for registroActualizado := 0; registroActualizado < len(history); registroActualizado++ {
		escribirtxt("Titan/"+history[registroActualizado].nombreSector, history[registroActualizado].nombreSector+" "+history[registroActualizado].base+" "+strconv.Itoa(history[registroActualizado].conteoSoldados))
		compendio = compendio + history[registroActualizado].nombreSector + " " + history[registroActualizado].base + " " + strconv.Itoa(history[registroActualizado].conteoSoldados) + ";"
	}

	canalMensajesTierra := CrearConexion("localhost", ":50052") //? Servidor Tierra
	canalMensajesMarte := CrearConexion("localhost", ":50054")  //? Servidor Marte
	for posicion := 0; posicion < 3; posicion++ {
		if relojTierra[posicion] > relojVectorial[posicion] && relojTierra[posicion] > relojMarte[posicion] {
			relojVectorial[posicion] = relojTierra[posicion]
		} else if relojVectorial[posicion] > relojTierra[posicion] && relojVectorial[posicion] > relojMarte[posicion] {
			// se queda como esta
		} else {

			relojVectorial[posicion] = relojMarte[posicion]
		}
	}
	Consistencia(canalMensajesTierra, compendio, relojVectorial)
	Consistencia(canalMensajesMarte, compendio, relojVectorial)
}

func Merge() {
	for {
		time.Sleep(60 * time.Second)
		canalMensajesTierra := CrearConexion("dist026", ":50052") //? Servidor Tierra
		canalMensajesMarte := CrearConexion("dist028", ":50054")  //? Servidor Marte
		relojTierra, numRegistrosTierra, relojMarte, numRegistrosMarte := ImportarLogs(canalMensajesMarte, canalMensajesTierra)
		GenerarConsistencia(relojTierra, numRegistrosTierra, relojMarte, numRegistrosMarte)
	}

}

func main() {
	fmt.Println("Escuchando...")
	go Merge()
	listener, err := net.Listen("tcp", ":50053") // Listener conexion sincrona
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
