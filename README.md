# Password Critic [![codecov](https://codecov.io/gh/samiam2013/passwordCritic/branch/main/graph/badge.svg?token=GDEPYIjlBw)](https://codecov.io/gh/samiam2013/passwordCritic) [![Build Status](https://app.travis-ci.com/samiam2013/passwordCritic.svg?branch=main)](https://app.travis-ci.com/samiam2013/passwordCritic) [![Go Report Card](https://goreportcard.com/badge/github.com/samiam2013/passwordcritic)](https://goreportcard.com/report/github.com/samiam2013/passwordcritic) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)


###  github.com/samiam2013/passwordCritic
a Go module to

* check the entropy of a password (variety & repitition)
* check if the password shows up in a bloom filter built with a list of common passwords

## Installing & Running
```
git clone https://github.com/samiam2013/passwordCritic.git
cd passwordCritic/
go test ./... # make sure the tests pass
cd cmd/
go build .
./cmd -r -p password123
```
and you can expect some output like
```
Entropy of the password candidate:  3.2776136
2022/01/13 15:25:08 password common, found in list with 100000 elements, but not more common than 1000, minimum set rarity.
{StringVal:password123 Cardinality:10 H:3.2776136}
```

## TODO:
- [] test frequency of false positives for a given BloomFilter.bitSet size
- [] benchmark algorithms used by bloom filter for hashing 
- [x] serialize the filter to something more efficient than the 0|1 mapped JSON
- [] add a backup source for commmon password lists
