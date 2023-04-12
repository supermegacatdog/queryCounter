# Query Counter

Counts queries from a big input file and writes them and their frequencies to an output file


```bash
./queryCounter input.txt output.txt N M
```

> N - max size of elements in map
>
>M - mode

### Modes:

> 1 - reading using bufio.Reader
> 
> 2 - reading using bufio.Scanner
