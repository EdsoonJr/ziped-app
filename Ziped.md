# Desafio de Desenvolvimento — Compactador de Arquivos + Player de Música

Este documento define dois projetos desktop que deverão ser desenvolvidos utilizando obrigatoriamente a stack:

- **Wails**
- **Go**
- **Svelte**
- **TypeScript**, quando necessário no frontend

O objetivo é produzir aplicações funcionais, testáveis e com escopo compatível com os prazos definidos.

O uso de IA como ferramenta de apoio ao desenvolvimento é permitido e incentivado, incluindo:

- geração de código;
- explicação de bibliotecas;
- debugging;
- testes;
- arquitetura;
- refatoração;
- documentação.

Porém, o desenvolvedor deverá compreender e conseguir explicar as principais decisões técnicas adotadas.

---

# Regras gerais de arquitetura

As duas aplicações devem seguir, preferencialmente, a seguinte separação:

## Backend — Go

Responsável por:

- regras de negócio;
- acesso a arquivos;
- compactação/descompactação;
- reprodução/processamento de áudio;
- leitura e escrita de metadados;
- persistência;
- integrações com bibliotecas nativas;
- operações pesadas;
- comunicação com o sistema operacional.

## Desktop Runtime — Wails

Responsável por:

- empacotamento desktop;
- integração entre Go e frontend;
- chamadas entre frontend e backend;
- acesso às APIs nativas;
- geração do executável.

## Frontend — Svelte

Responsável por:

- interface;
- navegação;
- componentes;
- estados visuais;
- feedback de progresso;
- listas;
- controles.

## TypeScript

Utilizar quando trouxer benefício real para:

- tipagem de dados recebidos do backend;
- contratos de API;
- stores;
- modelos;
- componentes;
- validações;
- utilitários frontend.

Evitar colocar regras de negócio críticas somente no frontend.

---

# PROJETO 1 — Compactador e Gerenciador de Arquivos

## Prazo

**Máximo: 4 dias**

---

# Objetivo

Desenvolver um aplicativo desktop inspirado na experiência de uso de ferramentas como WinRAR e 7-Zip.

O aplicativo deverá possuir:

- formato próprio;
- criação e extração de ZIP;
- leitura e extração de RAR;
- leitura e extração de 7z;
- interface gráfica desktop;
- gerenciamento de arquivos internos;
- tratamento adequado de erros.

O projeto não precisa copiar visualmente o WinRAR.

O objetivo é reproduzir o conceito de um gerenciador de arquivos compactados.

---

# Formatos obrigatórios

| Formato | Abrir/Listar | Extrair | Criar | Alterar |
|---|---:|---:|---:|---:|
| Formato próprio | Sim | Sim | Sim | Sim |
| ZIP | Sim | Sim | Sim | Sim |
| RAR | Sim | Sim | Não | Não |
| 7z | Sim | Sim | Não | Não |

RAR e 7z serão tratados como formatos somente leitura.

---

# TASK 01 — Estrutura inicial do projeto

Criar o projeto utilizando:

- Wails;
- Go;
- Svelte;
- TypeScript quando necessário.

Separar minimamente:

```text
/backend
    /archive
    /formats
    /services
    /models
    /storage

/frontend
    /src
        /components
        /stores
        /lib
        /types
        /views
```

Não é obrigatório seguir exatamente essa estrutura, mas deverá existir separação clara entre:

- UI;
- manipulação de arquivos;
- compactação;
- formatos;
- serviços;
- modelos;
- persistência;
- tratamento de erros.

## Critérios de aceite

- projeto compila;
- Wails inicia corretamente;
- frontend Svelte comunica com Go;
- README explica como executar;
- README explica como gerar build.

---

# TASK 02 — Interface principal

Criar uma interface de gerenciamento de arquivos compactados.

A tela principal deverá possuir pelo menos:

- Novo arquivo;
- Abrir arquivo;
- Adicionar arquivos;
- Adicionar pasta;
- Extrair;
- Excluir item;
- Testar integridade.

Ao abrir um pacote, exibir:

- nome;
- caminho interno;
- tamanho original;
- tamanho compactado, quando disponível;
- tipo;
- data de modificação.

Permitir navegar por diretórios existentes dentro do arquivo compactado.

## Critério de aceite

Todas as funções principais devem poder ser utilizadas pela interface gráfica sem necessidade de terminal.

---

# TASK 03 — Integração Wails entre frontend e backend

Criar uma camada de serviços Go exposta ao frontend através do Wails.

Exemplo conceitual:

```text
OpenArchive()
CreateArchive()
AddFiles()
AddFolder()
ExtractFiles()
DeleteEntry()
TestArchive()
CancelOperation()
```

O frontend deverá receber do backend:

- resultados;
- mensagens de erro;
- progresso;
- informações dos arquivos.

## Critério de aceite

A UI não deve implementar diretamente operações de compactação ou acesso pesado ao sistema de arquivos.

---

# TASK 04 — Formato próprio

Criar um formato próprio para a aplicação.

O desenvolvedor deverá definir uma extensão própria.

Exemplo:

```text
.garc
```

A extensão pode ser diferente.

O formato deverá suportar:

- múltiplos arquivos;
- diretórios;
- nomes Unicode;
- arquivos vazios;
- diretórios vazios;
- timestamps;
- tamanho original;
- tamanho armazenado;
- checksum;
- versão do formato.

Estrutura conceitual possível:

```text
MAGIC BYTES
VERSION
HEADER
METADATA
FILE TABLE
COMPRESSED DATA
CHECKSUM
```

Não é obrigatório criar um algoritmo de compressão próprio.

Podem ser utilizados algoritmos existentes, por exemplo:

- Deflate;
- LZMA;
- Zstandard;
- outro compatível com o projeto.

## Critérios de aceite

Um arquivo criado pelo aplicativo deverá:

1. ser salvo;
2. ser fechado;
3. ser aberto novamente;
4. listar corretamente seu conteúdo;
5. ser extraído;
6. manter a estrutura original dos arquivos.

---

# TASK 05 — Compactação do formato próprio

Implementar criação do formato próprio.

Permitir:

- selecionar arquivos;
- selecionar pastas;
- escolher destino;
- iniciar compactação;
- acompanhar progresso;
- cancelar operação.

## Requisito técnico

Arquivos grandes não devem necessariamente ser carregados completamente em memória.

Dar preferência para:

- streams;
- buffers;
- processamento em blocos.

## Critério de aceite

Os arquivos extraídos devem ser binariamente equivalentes aos arquivos originais.

---

# TASK 06 — Suporte ZIP

Implementar suporte completo para ZIP.

Obrigatório:

- criar;
- abrir;
- listar;
- extrair;
- adicionar arquivos;
- remover arquivos.

Quando a biblioteca escolhida exigir recriação do ZIP para alteração, isso é aceitável.

## Critérios de aceite

Um ZIP criado pela aplicação deverá abrir normalmente em softwares externos.

Um ZIP criado externamente deverá abrir corretamente na aplicação.

---

# TASK 07 — Suporte RAR somente leitura

Implementar:

- abrir;
- listar;
- navegar;
- testar;
- extrair.

Não implementar criação de RAR.

Não implementar alteração de RAR.

## Critério de aceite

Ao abrir um RAR, operações de escrita devem estar indisponíveis ou apresentar uma mensagem explicando que o formato é somente leitura.

---

# TASK 08 — Suporte 7z somente leitura

Implementar:

- abrir;
- listar;
- navegar;
- testar;
- extrair.

Não implementar criação ou alteração de arquivos 7z.

## Critério de aceite

Arquivos 7z válidos devem ser listados e extraídos corretamente.

---

# TASK 09 — Drag and Drop

Adicionar suporte a drag and drop.

Quando nenhum pacote estiver aberto:

- permitir iniciar criação de novo arquivo.

Quando um pacote editável estiver aberto:

- adicionar arquivos ao pacote.

## Critério de aceite

Deve ser possível arrastar múltiplos arquivos do sistema operacional para a aplicação.

---

# TASK 10 — Progresso de operações

Operações demoradas deverão enviar progresso do Go para o frontend.

Exemplos:

- compactação;
- extração;
- análise de integridade.

A interface deverá apresentar:

- percentual;
- arquivo atual;
- estado da operação.

Quando possível, implementar cancelamento.

---

# TASK 11 — Tratamento de erros

Tratar pelo menos:

- arquivo inexistente;
- arquivo corrompido;
- formato desconhecido;
- falta de permissão;
- arquivo bloqueado;
- destino indisponível;
- disco sem espaço;
- tentativa de escrita em RAR;
- tentativa de escrita em 7z.

## Critério de aceite

Nenhuma dessas situações deverá encerrar inesperadamente a aplicação.

---

# TASK 12 — Build

Gerar executável da aplicação usando Wails.

O usuário final não deverá precisar instalar:

- Go;
- Node.js;
- npm;
- ambiente de desenvolvimento.

---

# Casos de teste — Compactador

## CT-ARQ-001 — Criar pacote próprio

### Pré-requisitos

Possuir uma pasta contendo arquivos de diferentes tipos.

### Passos

1. Abrir o aplicativo.
2. Criar novo pacote.
3. Adicionar uma pasta.
4. Salvar.
5. Fechar o aplicativo.
6. Abrir novamente.
7. Abrir o pacote criado.

### Resultado esperado

Todos os arquivos e diretórios devem aparecer corretamente.

---

## CT-ARQ-002 — Extração do formato próprio

### Passos

1. Criar pacote próprio.
2. Adicionar arquivos.
3. Extrair para outro diretório.
4. Comparar com os originais.

### Resultado esperado

Os arquivos extraídos devem possuir o mesmo conteúdo dos arquivos originais.

---

## CT-ARQ-003 — Estrutura de diretórios

Utilizar:

```text
teste/
├── arquivo.txt
└── pasta1/
    ├── imagem.png
    └── pasta2/
        └── dados.json
```

Compactar e extrair.

### Resultado esperado

A hierarquia deve permanecer intacta.

---

## CT-ARQ-004 — Arquivo vazio

Adicionar arquivo de 0 bytes.

### Resultado esperado

O arquivo deve ser armazenado e restaurado corretamente.

---

## CT-ARQ-005 — Diretório vazio

Adicionar um diretório vazio.

### Resultado esperado

O diretório deve existir depois da extração.

---

## CT-ARQ-006 — Unicode

Compactar arquivos como:

```text
ação.txt
日本語.txt
🎵 musica.txt
```

### Resultado esperado

Os nomes devem permanecer intactos.

---

## CT-ARQ-007 — Arquivo grande

Compactar arquivo de pelo menos 1 GB.

### Resultado esperado

A aplicação deverá continuar responsiva e apresentar progresso.

---

## CT-ARQ-008 — ZIP externo

Criar ZIP utilizando software externo.

Abrir na aplicação.

### Resultado esperado

O conteúdo deve ser listado e extraído corretamente.

---

## CT-ARQ-009 — ZIP criado pela aplicação

Criar ZIP pela aplicação.

Abrir em outro software compatível.

### Resultado esperado

O ZIP deverá ser reconhecido e extraído corretamente.

---

## CT-ARQ-010 — RAR

Abrir um arquivo RAR válido.

### Resultado esperado

O conteúdo deve ser listado e extraído.

---

## CT-ARQ-011 — Alteração de RAR

Abrir RAR e tentar adicionar arquivo.

### Resultado esperado

A funcionalidade deverá estar bloqueada.

---

## CT-ARQ-012 — 7z

Abrir arquivo 7z válido.

### Resultado esperado

O conteúdo deve ser listado e extraído.

---

## CT-ARQ-013 — Arquivo corrompido

Modificar manualmente alguns bytes de um pacote próprio.

### Resultado esperado

A aplicação deve detectar falha de integridade.

Não deve fechar inesperadamente.

---

## CT-ARQ-014 — Extensão falsa

Renomear um arquivo TXT para a extensão própria da aplicação.

### Resultado esperado

O aplicativo deverá identificar o arquivo como inválido.

---

## CT-ARQ-015 — Cancelamento

Iniciar compactação de arquivo grande.

Cancelar durante a operação.

### Resultado esperado

A operação deverá parar sem travar a interface.

Um arquivo incompleto não deve ser considerado um pacote válido.

---

# Critérios de aceite gerais — Compactador

O projeto será considerado aprovado quando:

- possuir formato próprio;
- criar e extrair ZIP;
- ler e extrair RAR;
- ler e extrair 7z;
- preservar diretórios;
- preservar nomes Unicode;
- detectar corrupção no formato próprio;
- utilizar processamento adequado para arquivos grandes;
- apresentar progresso;
- possuir tratamento de erros;
- possuir interface gráfica;
- gerar executável;
- possuir README;
- frontend Svelte estiver devidamente integrado ao backend Go;
- operações pesadas não bloquearem a UI.

---

# Cronograma sugerido — Compactador

## Dia 1

- criar projeto Wails;
- configurar Svelte;
- arquitetura;
- comunicação frontend/backend;
- tela principal;
- iniciar formato próprio.

## Dia 2

- finalizar formato próprio;
- compactação;
- extração;
- suporte ZIP.

## Dia 3

- RAR;
- 7z;
- progresso;
- drag and drop;
- cancelamento.

## Dia 4

- tratamento de erros;
- testes;
- correções;
- build;
- documentação.

---