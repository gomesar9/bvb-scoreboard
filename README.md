## Para configurar a dependência do repositório privado

Configure o git para usar sempre ssh, ao invés de HTTPS.

Em `~/.gitconfig` adicione:

```
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
```

Depois execute os comando na raíz do projeto

```
go env -w GOPRIVATE=github.com/gomesar9
go get github.com/gomesar9/bvb-core@main
```
