# Trabajo practico en Go

Proyecto de GO en el cual consiste en jugar al UNO con arquitectura cliente-servidor.
Se debe levantar primero el servidor para poder conectar los clientes, el limite para
cada partida es de 3 clientes, en donde uno inicializa la partida y luego los demás 
clientes deben conectarse a la partida que desee entre una lista que se mostrará en
el lobby.


**Para lanzar el server ingresar al directorio server y lanzar**

    go run server.go

**Para lanzar el cliente ingresar al directorio client y lanzar**

    go run client.go

**Si hubiese algún package sin instalar se debe lanzar (ver si se requieren permisos)**

    go get -t .


## Ejemplos para ver

## Casos de Uso
Uber
https://eng.uber.com/go-geofence-highest-query-per-second-service/

Framework hecho por Twitch para Go
https://blog.twitch.tv/en/2018/01/16/twirp-a-sweet-new-rpc-framework-for-go-5f2febbf35f/

Manejo de accesos de infraestructura por Dailymotion
https://www.linkedin.com/pulse/golang-dailymotions-open-source-application-asteroid-reemi-shirsath/

Kubernetes
https://github.com/kubernetes/kubernetes
