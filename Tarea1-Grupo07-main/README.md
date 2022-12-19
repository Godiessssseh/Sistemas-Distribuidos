☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. * ･

# ( ͡° ͜ʖ ͡°)ﾉ README Laboratorio 1 Sistemas Distribuidos (◕ᴥ◕ʋ)

>Fecha entrega: Jueves 15 de septiembre  del 2022  
>**Integrante 1**: Diego Rosales Leon, *201810531-7*, paralelo 200  
>**Integrante 2**: Alan Zapata Silva, *201956567-2*, paralelo 200  

### Rabbitmq

- Rabbitmq se debe ejecutar a través del código:
    `sudo docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management`

- Si ya se está ejecutando, no hay que hacer nada.

### Compilar proto

- De llegar a ser necesario, correr el proto dentro de cada máquina virtual para que los servidores y toda la tarea funcionen de manera correcta
    `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Proto/message.proto`

##### Instrucciones ejecución

- Existe un makefile para ejecutar cada parte de la tarea:

- Para la central.
    `make central`

- Para el Laboratorio
    `make laboratorio`

##### Consideraciones a tomar en cuenta

- Cada maquina virtual tiene su puerto y host correspondiente, donde la central **debe** estar siendo ejecutada en el **dist025** con el **Laboratorio Kampala - Uganda**, el resto debe estar en máquinas virtuales diferentes:
- **dist025** es "*localhost*" con el puerto **50051** &rarr; **Laboratorio Kampala - Uganda**
- **dist026** es "*dist026*" con el puerto **50052** &rarr; **Laboratorio Pohang - Korea**
- **dist027** es "*dist027*" con el puerto **50053** &rarr; **Laboratorio Renca - Chile**
- **dist028** es "*dist028*" con el puerto **50054** &rarr; **Laboratorio Pripiat - Rusia**

- En el caso catastrófico de que se borren los archivos y sea necesario hacer un `pull`, se deben colocar los puertos respectivos en cada maquina correspondiente, **abriendo vim laboratorio.go** en cada máquina, y cambiar el labName por el nombre del laboratorio que corresponda y hostQ por la máquina virtual que corresponda. **(SOLO SI FUÉSE NECESARIO)**

- Para cerrar la ejecución de los servidores, primero se debe apretar `Ctrl+c` en la terminal de la maquina virtual donde se encuentra la central, luego se debe apretar nuevamente `Ctrl+c` en cada terminal de laboratorio.

- En caso de que el archivo txt ya exista, se debe eliminar a través del codigo:
    `rm -f Solicitudes.txt`

- Si la central queda con emergencias en cola de rabbitmq, se debe correr **SOLAMENTE** el archivo make central hasta que queden 0 emergencias. Luego correr con normalidad como se menciona anteriormente.  

*Programado con Go v1.19.1, grpc v1.49.0, RabbitMQ v3.9.13, Protocol Buffer v3.12.4en VSCode , testeado en Ubuntu 22.04*  
☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ.