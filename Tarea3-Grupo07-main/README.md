☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. ･ ｡ﾟ☆ﾟ. * ･

# ( ͡° ͜ʖ ͡°)ﾉ README Laboratorio 3 Sistemas Distribuidos (◕ᴥ◕ʋ)

>Fecha entrega:  03 de Diciembre del 2022  
>**Integrante 1**: Diego Rosales Leon, *201810531-7*, paralelo 200  
>**Integrante 2**: Alan Zapata Silva, *201956567-2*, paralelo 200 

### Observaciones
1. En Guardianes.go especificamos que existen códigos de error, que serán equivalentes a:

   1.  Código -1: **Archivo en el sector X y base Y existe, pero no hay información de los soldados.**
   
   2.  Código -2: **No existe un archivo de registros del sector X.**

   3.  Código -3: **Existe un archivo de registros del sector X pero no existe la base Y.**

2. Si no quieren ingresar soldados, simplemente apretar un enter

3. En guardianes se considerará que si todos los valores del reloj son más recientes (un valor mayor al que se tiene guardado), se actualizará, en cualquier otro caso, nos quedaremos con el reloj existente
   
4. Para borrar los archivos .txt, se tienen que borrar con la manito :c 

5. Si el merge no funciona, porfavor revisar todo el esfuerzo puesto en el (tomó mucho sudor, sangre, lágrimas y tiempo) además si quitan la func merge que está en el main de titan.go, todo funcionará con normalidad exceptuando el merge.

6. El merge es independiente de Monotonic Reads y Read Your Writes. Para este último, se crearan valores cualquiera (no afectan al funcionamiento de la tarea) que se irán actualizando.

7. Si no ingresará soldados a una base, asignar el valor de 0 de manera manual (A veces se cae cuando solo se presiona un enter)

8. **Si se quieren revisar los archivos .txts, hay que ver el contenedor docker o también quitar el comando --rm a los makefile.**

### Compilar proto

- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Proto/message.proto

##### Instrucciones ejecución
- Existe un archivo docker que está compuesto de archivos makefile para ejecutar cada parte de la tarea:

Para broker: 
    `make docker-broker` 

Para vanguardia: 
    `make docker-vanguardia`

Para guardianes: 
    `make docker-guardianes` 

Para tierra: 
    `make docker-tierra`

Para marte: 
    `make docker-marte`

Para titan: 
    `make docker-titan`

##### Consideraciones a tomar en cuenta

- Cada maquina virtual tiene su puerto y host correspondiente ....

- **dist025** es "*localhost*" con el puerto **50051** &rarr; **Broker Rasputin**
- **dist026** es "*dist026*" con el puerto **50052** &rarr; **Tierra** y con el puerto **50055** &rarr; **Vanguardia**
- **dist027** es "*dist027*" con el puerto **50053** &rarr; **Titan**  y con el puerto **50056** &rarr; **Guardianes**
- **dist028** es "*dist028*" con el puerto **50054** &rarr; **Marte**

*Programado con Go v1.19.1, grpc v1.49.0, Protocol Buffer v3.12.4en VSCode , testeado en Ubuntu 22.04*  
☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.☆ﾟ.* ･ ｡ﾟ☆ﾟ. *･ ｡ﾟ☆ﾟ.* ･ ｡ﾟ☆ﾟ.