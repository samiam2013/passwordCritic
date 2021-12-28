# Password Critic [![codecov](https://codecov.io/gh/samiam2013/passwordCritic/branch/main/graph/badge.svg?token=GDEPYIjlBw)](https://codecov.io/gh/samiam2013/passwordCritic) [![Build Status](https://app.travis-ci.com/samiam2013/passwordCritic.svg?branch=main)](https://app.travis-ci.com/samiam2013/passwordCritic) [![Go Report Card](https://goreportcard.com/badge/github.com/samiam2013/passwordcritic)](https://goreportcard.com/report/github.com/samiam2013/passwordcritic) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)


###  github.com/samiam2013/passwordCritic
a Go module for critiquing your selection of passwords

first using the entropy of your chosen password, then a set of bloom filters for common passwords. 

if the password is bad, the code will explain why.

## the current plan

check the entropy of the password, stopping things like "12341234"

if it passes an entropy check, look for it in a bloom filter of the top 10,000,000 most common passwords, 

if it's found in the 10MM list, look it up in a 10k and then 100 item bloom filter, 

if you find it in a 10k list from that 10MM, you get back an error

## TODO:
 * test frequency of false positives for a given BloomFilter.bitSet size
 * benchmark algorithms used by bloom filter for hashing 
 * serialize the filter to a file included in the library so the filter does not have to be computed from a list when included as a go module
 * look up the most ubiquitous way to get a list of things from the internet as a go module, ?maybe consider forking the list and using git from inside go for updates?
