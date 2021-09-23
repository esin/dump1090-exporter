# dump1090-exporter
Prometheus exporter for [antirez/dump1090](https://github.com/antirez/dump1090)

### Exported information

This exporter can export flight information:

```
...
dump1090_ac_tracking_flights{altitude="12925",flight="AZO435",geohash="ucfdvxww1rwe",lat="55.369418",long="37.516313",squawk="2115"} 1
dump1090_ac_tracking_flights{altitude="3375",flight="SBI3046",geohash="ucfu6bps42j8",lat="55.591507",long="37.748441",squawk="3332"} 1
...
```

Count of tracking AC:

```
...
dump1090_ac_tracking_now_count 2
...
```

Total catched messages

```
...
dump1090_messages_catched_count 1.006282e+06
...
```

### Grafana

![Grafana](https://github.com/esin/dump1090-exporter/blob/imgs/imgs/grafana.png?raw=true)
