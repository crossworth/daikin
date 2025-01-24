### MQTT Test

Programa de teste da biblioteca MQTT.

```
λλ go run *.go --username=<EMAIL> --password=<PASSWORD> --thingID=<THING_ID>
ThingState:

{
  "connected": 2,
  "src": 1,
  "port1": {
    "power": 1,
    "mode": 3,
    "temperature": 20,
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
      "out_temp": 24
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "ac_unit_type": "HEAT_PUMP"
}

Set AC to 25°C/Coandă effect/Economy Mode
ThingState:

{
  "connected": 2,
  "src": 1,
  "port1": {
    "power": 1,
    "mode": 3,
    "temperature": 20,
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
      "out_temp": 24
    },
    "rst_r": 12,
    "fw_ver": "p1.0.3.28"
  },
  "ac_unit_type": "HEAT_PUMP"
}

CTRL+C to stop
```