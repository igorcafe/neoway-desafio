# Rodando a aplicação

Executando a aplicação com as opções padrão:

```sh
docker compose up --build --abort-on-container-exit
```

Utilizando uma base 7x maior que a fornecida (350 mil linhas):

```sh
BASE_PATH=./resources/base_teste.350K.txt docker compose up --build --abort-on-container-exit
```

Utilizando a base de testes de 10 linhas:

```sh
BASE_PATH=./resources/base_teste.3.txt docker compose up --build --abort-on-container-exit
```

# Estratégias de otimização de performance utilizadas

## Execução de queries em batch

**Resultado 1** - batch com tamanho 50:
- 95% de redução no tempo total de execução (37s para 2s)

**Resultado 2** - batch com tamanho 1000:
- redução de 50% no tempo total de execução (10s para 5s, com dado 7x maior)

**Antes**: mais de 60% do tempo estava sendo gasto na função `Exec` do `pgx.Conn`:

![image](https://user-images.githubusercontent.com/85039990/222931883-3904cbd6-4ee9-4430-b470-4096f995a8e2.png)

**Depois**: só 20% do tempo é gasto com a função `SendBatch`

![image](https://user-images.githubusercontent.com/85039990/222931984-03e7d90d-f60d-473d-b440-0b3158581c46.png)

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

**Antes**: múltiplas alocações concatenando a string
![image](https://user-images.githubusercontent.com/85039990/222934212-ef1818f8-38ac-4d04-b7dd-7903e6460e19.png)

**Depois**: faz somente uma alocação com a função `Grow` do `strings.Builder`
![image](https://user-images.githubusercontent.com/85039990/222934236-68017b97-35f3-4722-9e4c-c7055545ddb6.png)

As duas funções anteriores faziam concatenação de strings, o que gerava várias alocações na heap desnecessárias.
Minha solução foi criar um `strings.Builder` e aumentar a capacidade dele para o tamanho da string.
Dessa forma era garantido que não faria mais de uma alocação porque a string gerada sempre vai ser igual ou menor.

Por exemplo:
```go
func SanitizeTicket(val string) string {
  res := &strings.Builder{}
  res.Grow(len(val))

  for _, r := range val {
    if (condicao) {
      res.WriteRune(r)
    }
  }

  // len(res) <= len(val), sempre
  return res.String()
}
```

## Processar as linhas de forma concorrente

Resultado:
- pouco significativo

Basicamente o que fiz foi processar N linhas concorrentemente, onde N é o número de CPUs lógicas (`runtime.NumCPU()`).
Essas linhas eram então enviadas para uma outra goroutine responsável por enfileirar a query no batch.

Eu imaginei que pudesse economizar alguns segundos porque as validações estavam tomando 20% do tempo total da aplicação, porém mesmo utilizando concorrência não tive melhora significativa.