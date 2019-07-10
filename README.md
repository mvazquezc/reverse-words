**Get reverse word**

```sh
curl http://127.0.0.1:8080/ -X POST -d '{"word":"PALC"}'
{"reverse_word":"CLAP"}
```

**Get release**

```sh
curl http://127.0.0.1:8080/ -X GET
```

**Get Health**

```sh
curl http://127.0.0.1:8080/health -X GET
```
