# Domainbadger

Domainbadger is a fast and simple CLI-tool used to find unregistered domains using wildcards & characterset combinations.

## Installation

### Homebrew

```sh
brew install emmanuelay/tap/domainbadger
```

## Demo

[![asciicast](https://asciinema.org/a/525969.svg)](https://asciinema.org/a/525969)

**Example**
```sh
domainbadger -c aoe -t se,io,nu h_ll_ w_rld d_min_ti_n
```

This will perform a WHOIS lookup on domain name combinations using selected top-level domains (`se`, `io`, and `nu`) and combinations of the characters `aoe` in three separate search terms `h_ll_`, `w_rld`, and `d_min_ti_n`. Underscore `_` is treated as the wildcard character and is replaced with different combinations of `aoe`.

`h_ll_` contains two wildcard characters, which means domainbadger will combine `aoe` in all possible combinations using 2 character slots (3 characters in 2 slots can be combined in 3^2 = 9 possible combinations: `halla`, `hallo`, `halle`, etc.).

## Usage
```
domainbadger [flags] <searchterms>

Flags:
  -A, --all             Use all possible characters (a-z, 0-9, -) (default true)
  -a, --alpha           Use alphabetic range (a-z)
  -n, --alphanum        Use alphanumeric range (a-z, 0-9)
  -c, --custom string   Use a custom character range (ex. abc123)
  -d, --delay int       Delay between lookup attempts, in milliseconds (default 500)
  -h, --help            Help for domainbadger
  -N, --numeric         Use numeric range (0-9)
  -t, --tld string      TLDs to search. Use comma to add multiple (ex. com,org,net) (default "com")
  -v, --version         Show version information
```

### Examples
```sh
# Search with custom characters
domainbadger -c aoe -t se,io h_ll_ w_rld

# Search with alphabetic characters
domainbadger -a -t com,net c_t d_g

# Search with alphanumeric characters
domainbadger -n -t io t_st

# Search with numeric characters
domainbadger -N -t com ap_ _pp

# Show version
domainbadger -v
```