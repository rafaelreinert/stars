# Stars

[![Go Report Card](https://goreportcard.com/badge/github.com/rafaelreinert/stars)](https://goreportcard.com/report/github.com/rafaelreinert/stars)
[![<rafaelreinert>](https://circleci.com/gh/rafaelreinert/stars.svg?style=svg)](https://github.com/rafaelreinert/stars)



## Instruções para Execução

Este projeto esta configurando com docker-compose e make, para executar o sistema basta usar as opcoes abaixo:

## make

- `make run` - Executa o `docker-compose build` executa o docker-compose com a imagem e uma instancia do MongoDB
- `make test` - Inicia o mongo com o docker-compose e executa os testes.


## API exemplos

Criacao do planeta:
``` curl
curl --location --request POST 'http://localhost:8080/planets' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Tund",
    "climate": "Arid",
    "terrain": "Dessert"
}'
```

Listagem de todos os planetas:
``` curl
curl --location --request GET 'http://localhost:8080/planets'
```

Busca por nome:
``` curl
curl --location --request GET 'http://localhost:8080/planets?name=Tund'
```

Busca por Id:
``` curl
curl --location --request GET 'http://localhost:8080/planets/5ef8c2d1c38c14ecf5ee6d75'
```

Remoçäo do planeta:
``` curl
curl --location --request DELETE 'http://localhost:8080/planets/5ef953f950d25d0f6f81b195'
```

Remoçäo do planeta:
``` curl
curl --location --request DELETE 'http://localhost:8080/planets/5ef953f950d25d0f6f81b195'
```

Update do planeta:
``` curl
curl --location --request PUT 'http://localhost:8080/planets/5ef9549050d25d0f6f81b196' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Mars",
    "climate": "Arid",
    "terrain": "Dessert"
}'
```
