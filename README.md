# Desafio DevOps 2026 — Projeto Korp

Implementação do desafio técnico (Partes 1, 2 e 3): serviço HTTP em Go, Docker Compose com NGINX como proxy reverso, monitoramento com Prometheus e Grafana, e provisionamento automatizado via Ansible.

## Estrutura do projeto

```
.
├── http-server-projeto-korp/   # Serviço Go + Dockerfile
├── nginx/                      # Proxy reverso (volume NGINX)
├── prometheus/                 # Configuração de scrape
├── grafana/                    # Provisionamento datasource + dashboard
├── docker-compose.yml          # Stack completa
└── ansible/                    # Playbook de automação
```

## Parte 1 — Serviço e containers

### Endpoint

- **GET** `http://localhost:80/projeto-korp` (via NGINX)
- Resposta JSON:

```json
{
  "nome": "Projeto Korp",
  "horario": "2026-06-01T12:00:00Z"
}
```

O campo `horario` é gerado em **UTC** a cada requisição.

### Subir manualmente

```bash
docker network create projeto-korp-net 2>/dev/null || true
docker compose up -d --build
curl http://localhost:80/projeto-korp
```

A rede `projeto-korp-net` é externa ao Compose (criada antes do `docker compose up`), conforme o item 3 da Parte 1.

O container `http-server-projeto-korp` **não** expõe porta no host; apenas o NGINX publica a porta **80**.

## Parte 2 — Monitoramento

### Métricas Prometheus (`/metrics`)

| Métrica | Tipo | Descrição |
|---------|------|-----------|
| `projeto_korp_service_available` | gauge | Disponibilidade (1 = up) |
| `projeto_korp_requests_total` | counter | Volume de requisições |

### Acessos

| Serviço | URL |
|---------|-----|
| API (proxy) | http://localhost:80/projeto-korp |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (usuário/senha: `admin` / `admin`) |

O dashboard **http-server-projeto-korp** é provisionado automaticamente em `grafana/dashboards/`.

Para gerar tráfego e visualizar métricas:

```bash
for i in $(seq 1 20); do curl -s http://localhost:80/projeto-korp > /dev/null; done
```

## Parte 3 — Ansible (um único comando)

Pré-requisitos no host de controle: `ansible` e permissão `sudo`.

```bash
cd ansible
ansible-galaxy collection install -r requirements.yml
ansible-playbook -i inventory.ini playbook.yml
```

O playbook:

1. Instala e inicia o Docker (Debian/Ubuntu)
2. Cria a rede bridge `projeto-korp-net`
3. Faz build da imagem `http-server-projeto-korp`
4. Sobe a stack com Docker Compose (NGINX, Prometheus, Grafana)
5. Valida `GET http://localhost:80/projeto-korp` e exibe o JSON no console

## Escolhas técnicas

- **Métricas**: exposição no formato texto do Prometheus em `/metrics`, sem dependências externas no binário Go.
- **Rede**: bridge `projeto-korp-net` compartilhada entre todos os containers.
- **Grafana**: provisionamento via `datasources.yml`, `dashboards.yml` e JSON do dashboard (bônus do desafio).
- **Ansible**: coleção `community.docker` para rede, imagem e integração com o engine Docker.

## Entrega

