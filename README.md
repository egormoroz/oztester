## Usage

```bash
$ oztester -h
Usage of oztester:
  -a    print testcase result immediatelly in async mode
  -i string
        path to tests directory (default "./tests/prob1")
  -n int
        concurrency (default 1)
  -nc
        disable colored output
  -t duration
        time limit (default 1s)
  -v    verbose
  -w    normalize whitespace
  -x string
        path to executable (default "./bin/sol1")
```

## Examples

```bash
$ ls ./tests/prob1 
1  10  10.a  11  11.a  12  12.a  13  13.a  16  16.a  17  17.a  18  18.a  19  19.a  1.a  2  20  20.a  21  21.a  22  22.a  23  23.a  24  24.a  2.a  3  3.a  4  4.a  5  5.a  6  6.a  7  7.a  8  8.a  9  9.a

$ oztester -x ./bin/sol1 -i ./tests/prob1       
22 cases found: [1 2 3 4 5 6 7 8 9 10 11 12 13 16 17 18 19 20 21 22 23 24]

1   OK
2   OK
3   OK
4   OK
5   OK
6   OK
7   OK
8   OK
9   OK
10  OK
11  OK
12  OK
13  OK
16  OK
17  OK
18  OK
19  OK
20  OK
21  OK
22  OK
23  OK
24  OK

[22/22] OK 22 WA 0 TL 0 CC 0 ERR 0
```
