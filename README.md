# Powerball
Powerball is benchmark tool with [tsenart/vegeta: HTTP load testing tool and library. It's over 9000!](https://github.com/tsenart/vegeta).
Powerball try to test for zore packet loss with binary search.

## Usage
```
$ powerball -debug -rate 100  http://192.0.2.1
Max: 100, Min: 0, Mid: 50, success: 1.000000
Max: 100, Min: 50, Mid: 75, success: 1.000000
Max: 100, Min: 75, Mid: 87, success: 1.000000
Max: 100, Min: 87, Mid: 93, success: 1.000000
Requests [total, rate, throughput]: 930, 93.100463, 93.094404
Durationt [total, attack, wait]: 9.99s, 9.99s, 0.00s
Success [ratio] 100.00%
```
