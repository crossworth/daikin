### SDK Daikin Split EcoSwing Smart R-32 e Split EcoSwing Smart Gold R-32

![Daikin](assets/running.png)

Controle e veja o status do ar-condicionado
[Daikin Split EcoSwing Smart R-32](https://www.daikin.com.br/produto/ecoswing-r32) e
[Daikin Split EcoSwing Smart Gold R-32](https://www.daikin.com.br/produto/split-ecoswing-smart-gold-r-32).

Para poder comunicar com o ar-condicionado é preciso de uma `secret key` e saber o endereço de IP do ar-condicionado.

Essa `secret key` é gerada pelo aplicativo
[Daikin Smart AC - Brasil](https://play.google.com/store/apps/details?id=in.co.iotalabs.dmb.smartac&hl=pt_BR&gl=US)
durante a configuração do ar-condicionado, ela também é salva no servidor que o aplicativo utiliza, dessa forma é
possível instalar o aplicativo em diferentes dispositivos e controlar o mesmo aparelho.

O ar condicionado possui duas formas de comunicação, por rede local utilizando a `secret key` ou por MQTT falando com um
servidor da AWS tópicos específicos, quando o aplicativo está fora da rede que o ar foi configurado, ele utiliza o
servidor MQTT para comunicação.

### Status

- [x] Suporte a extrair o `secret key` utilizando login e senha.
- [ ] Suporte a configurar o ar-condicionado sem necessidade de aplicativo.
- [x] Consultar estado do ar-condicionado (servidor http).
- [x] Enviar comandos para o ar-condicionado (servidor http).
- [x] Consultar estado do ar-condicionado (mqtt).
- [x] Enviar comandos para o ar-condicionado (mqtt).

### Compatibilidade

Deve ser compatível com todos os aparelhos que utilizam o aplicativo **Daikin Smart AC - Brasil**.

#### Daikin Split EcoSwing Smart R-32:

| Unidade interna | Status         |
|-----------------|----------------|
| FTKP09Q5VL      | Deve funcionar |
| FTKP12Q5VL      | Deve funcionar |
| FTKP18Q5VL      | Deve funcionar |
| FTKP24Q5VL      | Deve funcionar |
| FTHP09Q5VL      | Funcionando    |
| FTHP12Q5VL      | Deve funcionar |
| FTHP18Q5VL      | Funcionando    |
| FTHP24Q5VL      | Deve funcionar |

#### Daikin Split EcoSwing Smart Gold R-32:

| Unidade interna | Status         |
|-----------------|----------------|
| FTKP09S5VL      | Deve funcionar |
| FTKP12S5VL      | Funcionando    |
| FTKP18S5VL      | Deve funcionar |
| FTKP24S5VL      | Deve funcionar |
| FTHP09S5VL      | Deve funcionar |
| FTHP12S5VL      | Funcionando    |
| FTHP18S5VL      | Deve funcionar |
| FTHP24S5VL      | Deve funcionar |

#### Possívelmente funciona (unidades voltadas para o mercado da Índia)

| Modelo  |
|---------|
| FTKR35U |
| FTKR50U |
| FTKR60U |

### Conseguindo uma `secret key`

O primeiro passo é configurar o ar utilizando o aplicativo oficial, durante o processo de configuração é criado
a `secret key`.

Depois disso, você pode conseguir a `secret key` de duas formas, inspecionando as requests que o
aplicativo oficial faz, especificamente para o endpoint `https://dmb.iotalabs.co.in/devices/thinginfo/managething`.

Você também pode utilizar o seguinte
site [https://daikin-extract-secret-key.fly.dev/](https://daikin-extract-secret-key.fly.dev/).
O código do serviço está presente em `cmd/extract-secret-key` e nenhuma informação de login/senha/dispositivo é
coletada.

Você também pode executar o serviço localmente.

### Engenharia reversa

Todo o processo de descoberta e implementação foi baseado na engenharia reversa do aplicativo para Android e da leitura
da biblioteca nativa que o aplicativo utiliza.

#### Aplicativo para Android

O aplicativo é **extremamente** lento considerando que faz algo tão simples, o motivo disso é que ele ~~não guarda
nenhum dado local~~ não guarda todos os dados localmente, toda vez que o aplicativo é aberto ou restaurado (`onResume`)
ele faz requests para o [Amazon Cognito](https://aws.amazon.com/pt/cognito/) para validar o cadastro do usuário,
conseguir chaves de autenticação para depois fazer requests para o servidor de registros dispositivos, só para então
tentar se comunicar com ar localmente. Ele também faz algumas requests para verificação de versão do aplicativo e algum
tipo de logging de comportamento do usuário (fora as requests para o firebase de logging de aplicação).

Os servidores AWS utilizados ficam nos Estados Unidos, tornando as requests ainda mais lentas, considerando que são
diversas e em cascata (uma request aguarda a outra completar para poder continuar).

O servidor de registros dispositivos armazena diversos dados do dispositivo, é feito tracking de versões do
Android, timezone, linguagem, localização (latitude, longitude), modelo dos aparelhos de ar-condicionado e também alguns
dados da rede wireless (como nome).

O aplicativo parece ser uma cópia
do [Daikin AC Manager-India](https://play.google.com/store/apps/details?id=in.co.iotalabs.daikin.smartac&hl=pt_BR)
feito pelo [iota labs](http://iotalabs.co.in/) e adaptado para o mercado brasileiro.

O aplicativo tem uma usabilidade terrível, com controles lentos e capacidades básicas.

#### Ar-condicionado

O ar-condicionado possui um servidor http não conformante aos specs, tornando complicada comunicação (golang `http`
e `curl` não aceitam a resposta inválida), além disso, a forma de comunicação é no minima peculiar, as requests são
compostas por um conjunto de bytes base64 encodados, onde o conteúdo dos bytes é o seguinte:

- [**Initialization vector**](https://en.wikipedia.org/wiki/Initialization_vector): 16 bytes
- **Payload**: N bytes começando da posição 16
- **CRC16**: checksum de 2 bytes, sendo os últimos 2 bytes.

O conteúdo é criptografado utilizando [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) no
modo [CFB](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_feedback_(CFB)), sendo o CRC ignorado no
processo de descriptografia.

Quando mandando uma mensagem criptografada devemos adicionar `BZ` ao final do payload por algum motivo obscuro.

#### Comunicação pelo servidor http (funciona somente dentro da rede local)

Existem diversos exemplos de código de comunicação com servidor local:

- Programa de terminal que exibe o estado do ar em loop [`/cmd/ac-server-read-loop`](/cmd/ac-server-read-loop).
- Serviço http que retorna estado do ar na porta 8080 [`/cmd/ac-server-http-service`](/cmd/ac-server-http-service).
    - `GET /`: Retorna um JSON com o estado do ar.
    - `POST /state`: Aceita um JSON com o estado desejado para o ar.
- Programa controla o estado do ar [`/cmd/daikin`](/cmd/daikin).
    - `daikin get`: Retorna o estado do ar.
    - `daikin set`: Define o estado do ar, para ver os possíveis estados do ar, veja o arquivo [state.md](state.md).

_Todos os exemplos acima trabalham com duas informações `secretKey` e `targetAddress`, que devem ser fornecidas como
flags para os programas, por exemplo: `daikin --secretKey=<SecretKey> --targetAddress=http://192.168.0.70:15914/ get`._

#### Comunicação por MQTT (funciona fora da rede local)

Os tópicos utilizados no MQTT são:

- `$aws/things/{ID_DO_DISPOSITIVO}/shadow/get`: Solicita o estado do dispositivo.
- `$aws/things/{ID_DO_DISPOSITIVO}/shadow/get/accepted`: Retorna o estado do dispositivo.
- `$aws/things/{ID_DO_DISPOSITIVO}/shadow/update`: Atualiza o estado do dispositivo, deve ser enviado como payload um
  JSON com as alterações de estado desejada.
- `$aws/things/{ID_DO_DISPOSITIVO}/shadow/update/accepted`: Retorna o resultado da solicitação de atualziação do estado
  do dispositivo.

Conseguir conectar no servidor MQTT é um processo mais complicado, já que exige a troca de diferentes chaves de
autenticação e tokens de identificação com a AWS.

Mais informações sobre os
tópicos e implementação [nesse link](https://docs.aws.amazon.com/iot/latest/developerguide/iot-device-shadows.html).

Para ver um exemplo de código MQTT consulte a pasta [`/cmd/mqtt-test`](/cmd/mqtt-test).

_O exemplo acima trabalha com login e senha da conta criada no aplicativo móvel, já que é preciso de troca de
informações com servidores da AWS, também é preciso saber o ID do dispositivo (`thingID`) antes, você
pode [https://daikin-extract-secret-key.fly.dev/](https://daikin-extract-secret-key.fly.dev/) para conseguir o ID do
dispositivo._

#### Docker

É possível utilizar os binários `daikin` e `ac-server-http-service` diretamente de um container docker:

```shell

docker run --rm ghcr.io/crossworth/daikin:latest /app/daikin -secretKey=<SecretKey> -targetAddress=http://192.168.0.71:15914/ get
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

```shell

docker run --rm -p 8080:8080 ghcr.io/crossworth/daikin:latest /app/ac-server-http-service -secretKey=<SecretKey> -targetAddress=http://192.168.0.71:15914/
2025/01/24 00:13:06 starting http server at :8080


```
