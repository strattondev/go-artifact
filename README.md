# go-artifact
Install
```
github.com/strattonw/go-artifact
```

## Usage
### DeckCode
#### Decoder
```go
//Blue/Red Example
d, err := deckcode.ParseDeck("ADCJQUQI30zuwEYg2ABeF1Bu94BmWIBTEkLtAKlAZakAYmHh0JsdWUvUmVkIEV4YW1wbGU_")

if err != nil {
	log.Fatal(err)
}

fmt.Println(d)
```
#### Encoder
Possibly coming soon
### CardSet API
```go
csr := cardset.Receiver{&http.Client{}}
cs, err := csr.RetrieveCardSet("00")

if err != nil {
	log.Fatal(err)
}

fmt.Println(cs)
```
