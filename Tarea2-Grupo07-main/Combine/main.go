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

	hostS := "localhost" //IP del NameNode -> Creo que es localhost

	for {
		fmt.Println("Bienvenido, estás listx?")
		fmt.Println("Ingresar el tipo de información: Logística/Financiera/Militar:")
		check := bufio.NewScanner(os.Stdin)
		check.Scan()
		TIPO := check.Text()
		fmt.Println("Ingresar el ID de su información (No debe repetirse con IDs anteriores)")
		check2 := bufio.NewScanner(os.Stdin)
		check2.Scan()
		ID := check2.Text()
		fmt.Println("Ingresar la DATA: ")
		check3 := bufio.NewScanner(os.Stdin)
		check3.Scan()
		DATA := check3.Text()
		formato := TIPO + " : " + ID + " : " + DATA

		//Aqui deberíamos mandar un mensaje grpc a NameNode!!

		connS, err := grpc.Dial(hostS+":50051", grpc.WithInsecure()) //!Ver que puerto elegimos para NameNode
		if err != nil {
			panic(" No se pudo conectar con el servidor" + err.Error())
		}

		serviceCliente := pb.NewMessageServiceClient(connS)
		res, _ := serviceCliente.Intercambio(context.Background(),
			&pb.Message{
				Body: formato,
				Data: "C",
			})
		if res.Body == "Perfecto" {
			fmt.Println("Mensaje enviado de manera correcta :D")
		}
		connS.Close()

	}

}
