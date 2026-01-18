# HomePiggyBank

HomePiggyBank to lekka aplikacja webowa działająca w przeglądarce.  
Pozwala na prywatne zarządzanie finansami bez integracji z bankami — wszystkie dane są przechowywane lokalnie i szyfrowane.

## Technologie
Aplikacja typu SPA z serwerowym renderowaniem (SSR). Została zbudowana w oparciu o:

- Go
- Templ
- HTMX
- Alpine.js
- Tailwind CSS
- Chart.js

## Cechy
- szybkie i responsywne działanie
- bezpieczne przechowywanie danych lokalnie
- brak konieczności logowania do banku
- łatwe uruchomienie lokalne

## Development

### Dostępne cele (targets):
```
make tailwind-watch
```
Obserwuje plik ./web/input.css i automatycznie przebudowuje style Tailwind CSS przy każdej zmianie.

```
make tailwind-build
```
Minifikuje style Tailwind CSS.

```
make templ-watch
```
Obserwuje zmiany w plikach *.templ i automatycznie generuje szablony.

```
make templ-generate
```
Generuje szablony przy użyciu polecenia templ.

```
make dev
```
Uruchamia serwer deweloperski z użyciem Air, który wspomaga hot-reload aplikacji Go podczas developmentu.

```
make build
```
Buduje aplikację produkcyjną:
- kompiluje Tailwind,
- generuje templaty,
- kompiluje aplikację Go do ./bin/HomePiggyBank.

```
make docker-build
```
Buduje obraz Docker używając docker-compose z pliku ./dev/docker-compose.yml.

```
make docker-up
```
Uruchamia kontener z docker-compose.

```
make docker-down
```
Zatrzymuje kontener.

```
make docker-clean
```
Zatrzymuje kontener i usuwa wolumeny oraz obrazy.