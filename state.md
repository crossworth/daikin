### Estado possível de ser definido no ar-condicionado

Quando utilizando o programa `daikin` ou fazendo um `POST` para o `ac-server-http-service` você pode informar o estado
desejado para o ar como um JSON com campos opcionais (quando campo é informado, é alterado estado dele, caso contrario
não é alterado).

#### Campos

- `power`: Ligado, **inteiro**, `1` ou `0` define se o ar está ligado ou não (efetivamente liga e desliga ele).
- `mode`: Modo de operação, **inteiro**, consultar a tabela de modos de operação.
- `temperature`: Temperatura, **float**, temperatura desejada, com casa decimal em zero.
- `fan`: Modo de operação do fan, **inteiro**, consultar a tabela de modos de fan.
- `v_swing`: Swing, **inteiro**, `1` ou `0` define se o swing está ligado ou não.
- `coanda`: Coandă effect ou Conforto, **inteiro**, `1` ou `0` define se o modo conforto está ligado ou não.
- `econo`: Econômico, **inteiro**, `1` ou `0` define se modo econômico está ligado ou não.
- `powerchill`: Potente (Powerful), **inteiro**, `1` ou `0` define se modo potente está ligado ou não.

##### Modos de operação

- `0`: Automático
- `2`: Desumidificar
- `3`: Resfriar
- `4`: Aquecer
- `6`: Ventilar

##### Modos de fan

- `3`: Baixa
- `4`: Média-Baixa
- `5`: Média
- `6`: Média-Alta
- `7`: Alta
- `17`: Automático
- `18`: Silencioso

#### Exemplo de comandos:

```shell

# Liga o ar-condicionado e define a temperatura como 25°C.
daikin -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/ set '{"power":1,"temperature":25}'

# Liga o ar-condicionado, define a temperatura como 25°C, modo econômico e conforto.
daikin -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/ set '{"power":1,"temperature":25,"coanda":1,"econo":1}'

# Desliga o ar-condicionado.
daikin -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/ set '{"power":0}'

```

#### Exemplo de retorno de comando

```shell

# Liga o ar-condicionado, define a temperatura como 25°C, modo econômico e conforto.
daikin -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/ set '{"power":1,"temperature":25,"coanda":1,"econo":1}'
{
  "port1": {
    "power": 1,
    "mode": 3,
    "temperature": 25,
    "fan": 17,
    "h_swing": 0,
    "v_swing": 0,
    "coanda": 1,
    "econo": 1,
    "powerchill": 0,
    "good_sleep": 0,
    "streamer": 0,
    "out_quite": 0,
    "on_timer_set": 0,
    "on_timer_value": 0,
    "off_timer_set": 0,
    "off_timer_value": 0,
    "sensors": {
      "room_temp": 25,
      "out_temp": 22
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "idu": 1
}

```
