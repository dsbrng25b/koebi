# Installation
Hier wird die Installation des folgenden Setups beschrieben:

<img src="img/koebi_overview.png" />

## Influx DB
### Installation
Die Installation der Influx Datenbank ist hier beschrieben:
https://docs.influxdata.com/influxdb/v1.1/introduction/installation/

```shell
curl -sL https://repos.influxdata.com/influxdb.key | sudo apt-key add -
source /etc/os-release
echo "deb https://repos.influxdata.com/debian jessie stable" | sudo tee /etc/apt/sources.list.d/influxdb.list
```
### Starten
```shell
sudo apt-get update && sudo apt-get install influxdb
sudo systemctl start influxdb
```
### DB erstellen
```shell
$ influx -precision rfc3339
Connected to http://localhost:8086 version 1.1.x
InfluxDB shell 1.1.x
> CREATE DATABASE temperature
```

## Grafana
Grafana kann grundsätzlich nach der folgenden Anleitung installiert werden: http://docs.grafana.org/installation/debian/. 
Da aber Grafana kein offizielles Paket für ARM zur Verfügung stellt, muss man auf dem Raspberry Pi ein wenig anders vorgehen.
Das Projekt [grafana-on-raspberry](https://github.com/fg2it/grafana-on-raspberry) bietet Pakete für ARM an.
```shell
curl -O https://github.com/fg2it/grafana-on-raspberry/releases/download/v4.1.0/grafana_4.1.0-1484083445_armhf.deb
dpkg -i grafana_4.1.0-1484083445_armhf.deb
/bin/systemctl enable grafana-server
/bin/systemctl start grafana-server
```
Nun kann man sich auf http://*raspberry_ip*:3000 verbinden und mit dem User *admin* und dem Passwort *admin* einloggen.

## Köbi
### Installation go
Zunächst muss man sich die Go Entwicklungsumgebung einrichten, damit man den Köbi bauen kann. Wie das funktioniert ist im Detail unter https://golang.org/doc/install beschrieben.

```shell
curl -O https://storage.googleapis.com/golang/go1.7.4.linux-armv6l.tar.gz
tar xzvf go1.7.4.linux-armv6l.tar.gz -C /usr/local/
PATH=$PATH:/usr/local/go/bin
mkdir $HOME/go
export GOPATH=$HOME/go
```
### Köbi bauen
```shell
go get github.com/dvob/koebi
go install github.com/dvob/koebi 
```
### Köbi installieren
```shell
mkdir /opt/koebi
cp $GOPATH/bin/koebi /opt/koebi/
cp $GOPATH/src/github.com/dvob/koebi/linux/koebi.service /lib/systemd/system/
systemctl enable koebi.service
systemctl start koebi.service
```
