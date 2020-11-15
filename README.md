# Cliente-servidor procesos Golang
## Protocolo TCP

Notas:

- Cuando se conecta un cliente al servidor, el cliente se desconecta del servidor después de haber recibido el proceso asignado por el servidor.
- Con lo anterior, pudiera ser necesario construir un protocolo, en el cual se le haga saber al servidor si es un nuevo cliente que requiere un proceso o, un cliente que retornará un proceso al servidor.
- Los clientes se conectan dos veces al servidor:
    - La primera vez para solicitar un proceso.
    - La segunda vez para retornar el proceso asignado, con el objetivo de que el servidor lo siga ejecutando.
- Cuando el servidor asigna un proceso a un cliente, el servidor lo deja de ejecutar.
- Los procesos nunca mueren cuando el cliente termina su ejecución, el proceso lo debe de seguir ejecutando el servidor.