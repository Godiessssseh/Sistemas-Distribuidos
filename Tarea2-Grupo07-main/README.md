☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. * ･

# ( ͡° ͜ʖ ͡°)ﾉ README Laboratorio 2 Sistemas Distribuidos (◕ᴥ◕ʋ)

>Fecha entrega:  02 de Noviembre del 2022  
>**Integrante 1**: Diego Rosales Leon, *201810531-7*, paralelo 200  
>**Integrante 2**: Alan Zapata Silva, *201956567-2*, paralelo 200  

### Compilar proto

- De llegar a ser necesario, correr el proto dentro de cada máquina virtual para que los servidores y toda la tarea funcionen de manera correcta
    `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Proto/message.proto`

##### Instrucciones ejecución

- Existe un makefile para ejecutar cada parte de la tarea:

- Para combine:
    `make combine`

- Para rebeldes:
    `make rebeldes`

- Para namenode:
    `make namenode`

- Para grunt:
    `make grunt`

- Para synth:
    `make synth`

- Para el cremator:
    `make cremator`


##### Consideraciones a tomar en cuenta

- Cada maquina virtual tiene su puerto y host correspondiente, donde el Combine **debe** estar siendo ejecutado en el **dist025** con el **NameNode**, el resto debe estar en máquinas virtuales diferentes:
- **dist025** es "*localhost*" con el puerto **50051** &rarr; **Namenode y Combine**
- **dist026** es "*dist026*" con el puerto **50052** &rarr; **Grunt**
- **dist027** es "*dist027*" con el puerto **50053** &rarr; **Synth**
- **dist028** es "*dist028*" con el puerto **50054** &rarr; **Cremator**

- **Cabe decir que Rebeldes puede ejecutarse en cualquier otra consola :D, en este escenario se definió en la dist026**

- **Libreria String** -> Se utilizo para usar un split de strings para separar la informacion eficazmente.
- **Libreria Errors** -> Cuando googlié, encontré que era una buena fuente para revisar errores en la creacion de archivos

- Para cerrar la ejecución de los servidores, primero se debe apretar `Ctrl+c` en la terminal de la maquina virtual donde se encuentra la central, luego se debe apretar nuevamente `Ctrl+c` en cada terminal de laboratorio.

- En caso de que el archivo txt ya exista y se quiera eliminar, se debe hacer de manera manual en cada carpeta (NameNode, Datanode), con el comando:
    `rm -f DATA.txt`
    
*Programado con Go v1.19.1, grpc v1.49.0, Protocol Buffer v3.12.4en VSCode , testeado en Ubuntu 22.04*  
☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ.
