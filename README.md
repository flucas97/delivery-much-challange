# Delivery Much Tech Challenge - Recipes API

### Como usar?

```shell
$ git clone github.com/flucas97/delivery-much-challange
$ cd delivery-much-challange
```

Faça uma copia do arquivo .env.development.example e altere o GIPHY_API_KEY com a sua API_KEY do Giphy.
```shell
$ mv .env.development.example env.development
$ vim env.development
```
Caso não tenha uma API_KEY do Giphy ainda, basta seguir este tutorial: https://support.giphy.com/hc/en-us/articles/360020283431-Request-A-GIPHY-API-Key

Após configurar a váriavel de ambiente, inicie o servidor:
```shell
$ make
```

### Rotas Disponíveis:
```
# GET /recipes/?i={ingredient},{ingredient}
```

### Exemplo de pesquisa por receitas, informando dois ingredientes:
```
localhost:9090/recipes/?i=onion,pasta
```

### Regras

1) Máximo de tres (3) e mínimo de um (1) ingredientes por pesquisa;
