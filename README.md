# Servidor HTTP com Monitoramento e Automação Ansible

## Arquitetura do Ambiente

A infraestrutura foi projetada para garantir isolamento de rede e simplicidade no gerenciamento:

- **Isolamento de Rede:** O serviço em Go (`http-server-projeto-korp`) roda na porta interna 8080 e não publica suas portas para o host externo. Ele pertence exclusivamente à rede virtual interna `korp-network`.
- **Proxy Reverso:** O Nginx atua como único ponto de entrada público na porta 80, redirecionando as chamadas para a rota `/projeto-korp` no container Go.
- **Segurança de Containers:** O Dockerfile da aplicação utiliza um processo de build em dois estágios (multi-stage) e define a execução sob um usuário comum (`appuser`), evitando privilégios de root.
- **Provisionamento Automático:** O Grafana carrega de forma nativa o Prometheus como datasource padrão e importa o dashboard de monitoramento pré-configurado (`dashboard.json`) ao ser iniciado.

## Estrutura de Arquivos

```
projeto-korp/
├── app/
│   ├── main.go         # Código fonte da API em Go
│   ├── go.mod          # Módulo e dependências do Go
│   └── Dockerfile      # Dockerfile otimizado da aplicação
├── infra/
│   ├── nginx/
│   │   └── http-server-projeto-korp.conf # Configuração do proxy reverso
│   ├── prometheus/
│   │   └── prometheus.yml                # Configuração do coletor de métricas
│   └── grafana/
│       ├── datasources/
│       │   └── datasource.yml            # Definição automática de fonte de dados
│       └── dashboards/
│           ├── dashboard.yml             # Mapeamento do dashboard do Grafana
│           └── dashboard.json            # Estrutura visual do painel de monitoramento
├── docker-compose.yml  # Orquestração local de todos os containers
└── playbook.yml        # Automação Ansible para instalação e deploy
```

## Como Executar e Testar

### Pré-requisitos
Certifique-se de possuir o Docker e o Docker Compose instalados no host de execução.

### Execução local via Docker Compose

1. Navegue até o diretório do projeto:
   ```bash
   cd projeto-korp
   ```
2. Inicialize os serviços em segundo plano:
   ```bash
   docker compose up --build -d
   ```
3. Teste o funcionamento do endpoint HTTP:
   ```bash
   curl http://localhost/projeto-korp
   ```
   **Retorno esperado:**
   ```json
   {"nome":"Projeto Korp","horario":"2026-05-27T19:38:34Z"}
   ```

### Implantação Automatizada com Ansible

O playbook Ansible foi construído para preparar todo o ambiente (instalando o Docker se necessário em ambientes Debian/Ubuntu), configurar a rede, compilar os containers locais e verificar a resposta final da aplicação no console.

Para testar localmente, execute:
```bash
ansible-playbook -i "localhost," -c local playbook.yml
```

## Monitoramento

O monitoramento pode ser acessado localmente através das portas mapeadas:

- **Prometheus UI:** http://localhost:9090
- **Grafana Dashboard:** http://localhost:3000 (Credenciais: admin / admin)

### Métricas Coletadas
- **Disponibilidade (Uptime):** Indicada pela métrica nativa `up` e pela métrica `korp_service_uptime_seconds`, exibindo o tempo de atividade acumulado da API.
- **Volume de Requisições:** Contagem total de chamadas HTTP feitas à aplicação, segmentada no dashboard de RPS (Requisições por Segundo) pelos respectivos códigos de status HTTP (como 200 e 405).
