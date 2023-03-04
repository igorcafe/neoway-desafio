# Rodando a aplicação
```sh
docker compose up --build
```

Utilizará por padrão a base de dados fornecida de 50 mil linhas.

# Estratégias de otimização de performance utilizadas

## Execução de queries em batch
**Resultado 1** - batch com tamanho 50:
- 95% de redução no tempo total de execução (37s para 2s)

**Resultado 2** - batch com tamanho 1000:
- redução de 50% no tempo total de execução (10s para 5s, com dado 7x maior)

Em vez enviar os comandos de `INSERT` um por um para o banco, a execução em batch envia vários inserts de uma vez.

Por mais que o banco de dados esteja rodando na mesma máquina, existe um custo de preparar o pacote TCP, enviar, esperar o ACK, fora a parte do protocolo do postgres em si.

Ao enviar em batch múltiplas queries são agrupadas em um mesmo pacote TCP reduzindo assim _round trips_ de rede, o que é uma possível explicação para a melhoria.

## Utilizar `strings.Builder` em vez de usar concatenação

Resultado do benchmark da função SanitizeCpfOrCnpj (executada 50 mil vezes):
- menos 1,3 milhões de alocações
- 7x mais rápida
- usando 9x menos memória

Resultado do benchmark da função SanitizeTicket (executada 50 mil vezes):
- menos 55 mil alocações
- 5x mais rápida
- usando 8x menos memória

As duas funções anteriores faziam concatenação de strings, o que gerava várias alocações na heap desnecessárias