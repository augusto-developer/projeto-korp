# Desafio Técnico - http-server-projeto-korp

Este repositório contém a resolução do desafio técnico de DevOps/Infraestrutura. O projeto engloba a criação de um serviço HTTP instrumentado em Go, orquestração de containers com Docker Compose, configuração de proxy reverso com Nginx, monitoramento completo com Prometheus/Grafana e automação de implantação via Ansible.

---

## 🏗️ Arquitetura do Ambiente

A arquitetura do projeto foi desenhada visando segurança, isolamento de rede e observabilidade:

```
[ Cliente ] 
     │ (Porta 80)
     ▼
┌──────────────────────────────────────────┐
│              nginx-proxy                 │
└────────────────────┬─────────────────────┘
                     │ (Rede Bridge: korp-network)
                     ▼
┌──────────────────────────────────────────┐
│      http-server-projeto-korp            │ (Porta 8080 - Não exposta ao Host)
└────────────────────┬─────────────────────┘
                     │
                     ├─────────────────────┐
                     ▼                     ▼
┌──────────────────────────┐   ┌──────────────────────────┐
│        Prometheus        │   │         Grafana          │
│      (Porta 9090)        │   │       (Porta 3000)       │
└──────────────────────────┘   └──────────────────────────┘
```

### Decisões de Arquitetura:
1. **Linguagem Go**: Escolhida por sua alta performance, baixo consumo de CPU/Memória, suporte nativo a concorrência e facilidade de integração com bibliotecas nativas de métricas do ecossistema Cloud Native (Prometheus).
2. **Isolamento de Rede (Segurança)**: A aplicação Go roda internamente na rede bridge `korp-network` exposta na porta `8080`. Ela **não expõe portas ao Host**, impossibilitando acesso externo direto. Todo o tráfego externo passa obrigatoriamente pelo Nginx na porta `80` atuando como Proxy Reverso.
3. **Dockerfile Otimizado (Multi-Stage)**:
   - **Estágio de build**: Utiliza uma imagem Go leve baseada em Alpine para compilar o executável estático, removendo símbolos de depuração (`-ldflags="-w -s"`) para reduzir o tamanho.
   - **Estágio de execução**: Utiliza uma imagem mínima do Alpine contendo apenas o binário compilado.
   - **Non-Root**: A aplicação executa sob um usuário sem privilégios (`appuser`), garantindo conformidade com as melhores práticas de segurança de containers (evitando escalada de privilégios).
4. **Grafana Automatizado (Provisionamento)**: O Grafana está configurado para carregar automaticamente o Prometheus como Datasource padrão e renderizar o dashboard pré-configurado (`dashboard.json`) ao ser iniciado, sem necessidade de qualquer clique manual.

---

## 📂 Estrutura do Projeto

```
projeto-korp/
├── app/
│   ├── main.go         # Servidor HTTP Go instrumentado com métricas Prometheus
│   ├── go.mod          # Definição das dependências do Go
│   └── Dockerfile      # Dockerfile multi-stage e seguro (non-root)
├── infra/
│   ├── nginx/
│   │   └── http-server-projeto-korp.conf # Arquivo de Proxy Reverso do Nginx
│   ├── prometheus/
│   │   └── prometheus.yml                # Configuração de scraping do Prometheus
│   └── grafana/
│       ├── datasources/
│       │   └── datasource.yml            # Provisionamento automático da fonte de dados
│       └── dashboards/
│           ├── dashboard.yml            # Configuração de carregamento do dashboard
│           └── dashboard.json            # Layout JSON do painel de monitoramento
├── docker-compose.yml  # Definição e orquestração de todos os containers
└── playbook.yml        # Playbook Ansible para automatizar 100% o provisionamento do zero
```

---

## 🚀 Como Executar e Validar

### Pré-requisitos
* Docker e Docker Compose instalados no host.
* *Opcional:* Ansible (caso queira realizar o provisionamento automatizado por playbook).

---

### Opção A: Execução Direta (Docker Compose)

1. Entre no diretório do projeto:
   ```bash
   cd projeto-korp
   ```
2. Inicie a pilha de serviços em segundo plano:
   ```bash
   docker compose up --build -d
   ```
3. Valide o funcionamento do serviço consumindo o endpoint pelo Nginx:
   ```bash
   curl http://localhost/projeto-korp
   ```
   **Resposta Esperada (JSON Dinâmico):**
   ```json
   {"nome":"Projeto Korp","horario":"2026-05-27T19:38:34Z"}
   ```

---

### Opção B: Implantação Automatizada (Ansible)

O playbook Ansible instala as dependências, cria a rede bridge, compila as imagens locais, orquestra os containers e realiza um teste de fumaça retornando a resposta da requisição HTTP no console em um único comando.

Para testar localmente (localhost):
```bash
ansible-playbook -i "localhost," -c local playbook.yml
```

---

## 📊 Monitoramento e Observabilidade

* **Prometheus UI**: Disponível em [http://localhost:9090](http://localhost:9090)
* **Grafana Dashboard**: Disponível em [http://localhost:3000](http://localhost:3000)
  * **Usuário**: `admin`
  * **Senha**: `admin`

### Painéis do Dashboard Customizado:
1. **Disponibilidade do Serviço**: Exibe `ONLINE` em verde se a aplicação estiver saudável, e `OFFLINE` em vermelho caso o container do Go pare.
2. **Tempo de Uptime (Atividade)**: Exibe em tempo real há quanto tempo a aplicação está ativa.
3. **Volume de Requisições**: Contador incremental que mede o volume total acumulado de requisições.
4. **Taxa de RPS (Requisições por Segundo)**: Gráfico de linha que exibe a taxa de tráfego por segundo segmentada por código de status HTTP (ex: 200 OK, 405 Method Not Allowed, 500 Server Error).
