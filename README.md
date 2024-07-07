# Thots

A Journal for your thots...

## Tool used
- [Golang](https://go.dev/)
- [Templ](https://templ.guide/)
- [Tailwindcss](https://tailwindcss.com/)
- [Htmx](https://htmx.org/)
- [Air](https://github.com/air-verse/air)

## Run locally

### Get started

Clone the repo with the below command
```sh
git clone https://github.com/pratikmane1299/thots.git
```

### Start server

Run the below command to build templ templates in watch mode
```sh
templ generate -w -proxy=https://localhost:4242
```

Build tailwindcss
```sh
npx tailwind -i style.css -o static/styles.css --watch
```

Start the go server with air
```sh
air
```


