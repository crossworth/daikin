### Programa para controlar o ar-condicionado

#### get

```shell

# Retorna os dados do ar-condicionado
daikin -secretKey=<SECRET_KEY> -targetAddress=http://192.168.0.71:15914/ get
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
      "out_temp": 22
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "idu": 1
}

```

#### set

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