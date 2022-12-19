package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	pb "github.com/Sistemas-Distribuidos-2022-2/Tarea2-Grupo07/Proto"
	"google.golang.org/grpc"
)

func main() {

	hostS := "dist026" //!IP del NameNode -> Hay que elegir!

	for {
		fmt.Println("¿Qué categoría te gustaria consultar?")
		fmt.Println("Logistica/Financiera/Militar/Cerrar: ") //? Si elige cerrar, tenemos que matar todos los nodos uwu
		check := bufio.NewScanner(os.Stdin)
		check.Scan()
		categoria := check.Text()

		//Aqui deberíamos mandar un mensaje grpc a NameNode!!

		connS, err := grpc.Dial(hostS+":50051", grpc.WithInsecure())
		if err != nil {
			panic("Error en la conexion con el servidor :c" + err.Error())
		}

		serviceCliente := pb.NewMessageServiceClient(connS)
		res, _ := serviceCliente.Intercambio(context.Background(),
			&pb.Message{
				Body: categoria,
				Data: "R", //? Enviamos una R de rebeldes!
			})
		fmt.Println("Mensaje enviado de manera correcta :D")
		fmt.Println(res.Body)
		connS.Close()

	}

}
