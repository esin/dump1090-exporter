#!/bin/bash
dt=$(date +%y%m%d%H%M)
docker build -t dump1090-exporter:latest -t dump1090-exporter:$dt -f Dockerfile --load .
if [ $? -ne 0 ]; then
  exit 1;
fi

docker tag dump1090-exporter:latest es1n/dump1090-exporter:latest
docker tag dump1090-exporter:$dt es1n/dump1090-exporter:$dt
docker push es1n/dump1090-exporter:$dt
docker push es1n/dump1090-exporter:latest
