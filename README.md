# Desafio 5 Pós Go Lang - Full Cycle

## Objetivo

. Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

. Criar um teste que valide a implementação.

## Funcionalidades

- O sistema disponibiliza endpoints para criar e buscar usuário
- O sistema disponibiliza endpoints para criar e listar leilões
- O sistema disponibiliza endpoints para criar e listar lances

## Requisitos

- Docker instalado

## Configuração

1. No diretório raiz do projeto use os comandos abaixo para executar o teste automatizado

   - Use o comando docker compose up -d
   - Use o comando docker compose exec app bash
   - Use o comando go test ./test/e2e

2. No diretório raiz do projeto use os comandos abaixo para executar o projeto

   - Use o comando docker compose up -d
   - Use o comando docker compose exec app bash
   - Use o comando go run cmd/auction/main.go
   - Use o arquivo test.http com a extensão REST Client ou siga os exemplos do arquivo
