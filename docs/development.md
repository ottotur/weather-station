# Kehitysympäristö

Käytä projektia WSL:n Linux-tiedostojärjestelmässä, ei Windowsin `/mnt/c`-polulla.
Toteutuksen alussa todettu ympäristö: Windows 11, WSL 2.7.14.0, Ubuntu
26.04.1 LTS, Linux x86_64; `wsl.exe --version`, `cat /etc/os-release` ja
`uname -m` näyttävät paikalliset tiedot.

Työkalut: Go 1.27.4 (Go-jakelun uusin vakaa versio tarkistushetkellä;
`go version`, tarvittaessa `GOTOOLCHAIN=go1.27.4 go version`), Node 24.21.0
ja npm 11.19.0 (`node --version`, `npm --version`). Käyttöliittymän lukitut
versiot ovat Svelte 5.57.3, Vite 8.3.8, Svelte Vite -liitännäinen 7.3.1,
TypeScript 5.9.3, Node-tyypit 24.19.2 ja svelte-check 4.7.8; versiot voi tarkistaa
komennolla `npm ls --depth=0` hakemistossa `web`.

Tämä on vasta keskeneräinen runko, ei toimiva vedenkorkeussovellus.

`npm audit` raportoi kehitysaikaisen Svelte Vite -liitännäisen `deepmerge`-
riippuvuudesta matalan vakavuuden haavoittuvuuden GHSA-gcqx-w48m-w2h4.
Tarkistushetkellä npm ei tarjonnut korjausta; riippuvuus ei ole selaimen
ajonaikainen riippuvuus.
