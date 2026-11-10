# demo-renovate-basic-service-plan

Service Go qui sert le plan SVG d'un entrepôt. Dépôt de démo de la série de blog sur Renovate de Clever Tech Ware : les
versions des dépendances, du Dockerfile et des actions GitHub sont volontairement anciennes pour que Renovate ait des
mises à jour à proposer. Elles ne doivent pas être relevées à la main.

## Ce que montre l'exemple

- `go.mod` : gin `v1.9.1`, `bun`, `pgdialect` et `pgdriver` en `v1.1.17`, directive `go 1.23.0`.
- `Dockerfile` : `golang:1.23-alpine` puis `alpine:3.20`.
- `.github/workflows/ci.yml` : `actions/checkout@v4`, `actions/setup-go@v5`, `go-version: '1.23'`.
- `renovate.json` et `.github/workflows/renovate.yml` : Renovate auto-hébergé, une instance par dépôt.

## Prérequis

- Go 1.23 ou plus récent
- Docker (les tests d'intégration démarrent un Postgres avec testcontainers)
- Une base Postgres `plan` pour lancer le service (voir le dépôt infra de la démo)

## Commandes

```bash
go build ./...        # compilation
go test ./...         # tests unitaires et d'intégration (Docker requis)
go run .              # lance le service sur le port 8081
docker build -t service-plan .
```

Variable d'environnement : `DATABASE_URL`, par défaut `postgres://plan:plan@localhost:5432/plan?sslmode=disable`.
Les tables `warehouse_plan` et `zone` sont créées et alimentées au démarrage, de façon idempotente. Les identifiants
Postgres de la démo (`plan` / `plan`) ne conviennent qu'à un usage local.

## Endpoints

| Route | Réponse |
|---|---|
| `GET /plans/{code}.svg` | plan SVG de l'entrepôt (`LYS-01`, `NTE-02`, `LIL-03`), 404 si le code est inconnu |
| `GET /healthz` | `{"status":"ok"}` |

## Série

Série de blog Renovate : <https://www.clevertechware.fr/blog>
