##Example config for apis:

### Herodotos:

`ENV=LOCAL;TLS_ENABLED=yes;ELASTIC_SEARCH_PASSWORD=acanbiC4bduS3GDv56CjiL1B;ELASTIC_SEARCH_USER=elastic;ELASTIC_ACCESS=text`

### dionysios:

`ENV=LOCAL;TLS_ENABLED=yes;ELASTIC_SEARCH_PASSWORD=8focJyDd0djZItxhZa;ELASTIC_SEARCH_USER=dionysios;ELASTIC_ACCESS=grammar;ELASTIC_SECONDARY_ACCESS=dictionary`

##Regex

Remove digits:

`\(*[0-9].*[0-9]\)`


Find any non digit with a 2 appended
`[\D]2`