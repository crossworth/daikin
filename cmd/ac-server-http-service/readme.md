### Servidor para controlar o ar-condicionado

#### Executando

```shell

ac-server-http-service -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/
2025/01/23 21:03:24 starting http server at :8080

```

Com isso o servidor vai estar sendo executado na porta `8080` de todas as interfaces de rede do computador.

Você pode fazer um `GET` em `/` para conseguir o estado do ar-condicionado e pode fazer um `POST` com o conteúdo em JSON
no body para controlar o dispositivo.

```shell

# Conseguir estado atual
curl http://localhost:8080/
{
  "port1": {
    "power": 0,
    "mode": 3,
    "temperature": 25,
    "fan": 17,
    "h_swing": 0,
    "v_swing": 0,
    "coanda": 0,
    "econo": 0,
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
      "out_temp": 255
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "idu": 1
}

```

```shell

# Liga o ar-condicionado, define a temperatura como 25°C, modo econômico e conforto.
curl -X POST http://localhost:8080/state -d '{"power":1,"temperature":25,"coanda":1,"econo":1}'
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
      "room_temp": 26,
      "out_temp": 255
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "idu": 1
}

```