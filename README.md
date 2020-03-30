# backend
Das Backend zu unserem Hackathon WirVsVirus Projekt DeMed.

## Desciption
DeMed matcht Nachfrage und Angebot von Medikamenten und medizinischem Equipment, in dem es auf einer Plattform einen schnellen Austausch zwischen (Krankenhaus-)Apotheken zulässt, um Medikamenten-Engpässe zu überbrücken. Die Abwicklung findet dann über Apotheken und den Pharmagroßhandel statt. 

## Get Started
[Install Go](https://golang.org/doc/install)

Das Projekt Klonen:
 ```bash 
git clone https://github.com/WirvsVirus-DeMed/backend 
```
``` bash
cd backend
```
Die Submodule des Projektes installieren:
``` bash 
git submodule update --init
```
``` bash 
git submodule update --remote 
```
Den Code ausführen (und golang Bibliotheken downloaden):
``` bash
go run ./main.go
```
ODER mit Docker:
``` bash
docker build -t demed .
````
``` bash
docker run -it -p 8080:8080 demed
```