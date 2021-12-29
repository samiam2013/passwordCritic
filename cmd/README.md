# passchk

Pass Check is a command line tool for using passwordCritic as a standalone tool.

by default it will give you back an entropy number for a string (generally ~2.0-5.0)
```
/home/chad$ pwcheck HelloWorld
0.0
warning: password found in common passwords bloom filter!
/home/chad$
```

## Ideas
 * `--json` for json formatted output of all the passwordCand struct info
 * `--verbose` for detailed information about the process (benchmark, intermediate values etc.)
 * `--never-too-common` for disabling the bloom filter check
 * `--suppress-warnings`