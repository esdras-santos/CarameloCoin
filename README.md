# CarameloCoin

CarameloCoin is a MemeCoin that I (Esdras Santos) wrote with Golang just to put into practice all the 
theoretical knowledge I acquired in researches about the blockchain core in books and papers (this project has 
no financial purpose, so don't put your money on it).

## Running

Clone this repo:

```shell
git clone https://github.com/esdras-santos/CarameloCoin
```

Make two copies with different names like `CarameloCoin` and `CarameloCoin2`, each copy will represent a full-node.

Then `cd` into the `CarameloCoin` directory:

```shell
cd CarameloCoin
```

Then in a different terminal `cd` into the `CarameloCoin2` directory:

```shell
cd CarameloCoin2
```

In `CarameloCoin` directory type:

```shell
go run main.go
```
this will initiate the blockchain with the genesis block and create the wallet with `50 caramel` as a reward of the genesis block.


Copy the full address that is returned:

```shell
    gen node
5ed004db888827562720cbcc54667b6f35971fb6e04a142752cdccb4c6f0434c
2021/11/05 09:31:30 listening for connections
2021/11/05 09:31:33 I am `/ip4/127.0.0.1/tcp/61617/p2p/Qme4LvaqiPDBAwcGaCsyU1UM4gyPLo5r7CpAhXjgieFq8K`
2021/11/05 09:31:33 peerid: Qme4LvaqiPDBAwcGaCsyU1UM4gyPLo5r7CpAhXjgieFq8K

address:  1J16Pki7mgBiP9iorFGhPJEZxDRP1UhQ5T

type "command" to see the list of commands

>>>
```

Now initiate the second node into the `CarameloCoin2` directory with the full address of the genesis node:

```shell
go run main.go `/ip4/127.0.0.1/tcp/61617/p2p/Qme4LvaqiPDBAwcGaCsyU1UM4gyPLo5r7CpAhXjgieFq8K`
```
the output will be something like this:

```shell
2021/11/05 09:42:30 I am /ip4/127.0.0.1/tcp/63650/p2p/QmaT6yQ6nsAYbsUUWX4iq3AUFvWtqaeP9z7grvUqbsBzuz
2021/11/05 09:42:30 peerid: QmaT6yQ6nsAYbsUUWX4iq3AUFvWtqaeP9z7grvUqbsBzuz

address:  1MsKpWvc9ipxGm8ZTa6K6qZo1SuwRD7uf9

type "command" to see the list of commands

>>>
```

Now the two nodes are connected in a p2p network.

## Commands

Transact `10 caramels` from the `CarameloCoin` node to the `CarameloCoin2` node with the `send` command:

```shell
>>> send

        To: 1MsKpWvc9ipxGm8ZTa6K6qZo1SuwRD7uf9

        Amount: 10
```
if the it's a valid transaction it will spread through the network.


Mine the previous transaction with the `mine` command:

```shell
>>> mine
5d05d37f9db658528520c226f16da353a965bfabd2c58eef26a33311d32e2625
```
it will return the hash of the mined block.

Check The balances with the `balance` command:

check the `CarameloCoin` node balance:
```shell
>>> balance

        Of: 1J16Pki7mgBiP9iorFGhPJEZxDRP1UhQ5T

        balance: 90
```
`50 caramels` as reward of the mined block plus the previous balance of `40 caramels`.

check the `CarameloCoin2` node balance:
```shell
>>> balance

        Of: 1MsKpWvc9ipxGm8ZTa6K6qZo1SuwRD7uf9

        balance: 10
```
`10 caramels` that was sended by the `CarameloCoin` node.




