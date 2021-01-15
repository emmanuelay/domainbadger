# Badger

Badger is a CLI-tool used to find unregistered domains using wildcards & characterset combinations. 

**Example**
```sh
badger -custom aoe -tld se,io,nu h_ll_ w_rld d_min_ti_n
```

This will perform a whois lookup on domain name combinations using selected top-level domains (`se,io` and `nu`) and combinations of the characters `aoe` in three separate search terms `h_ll_`, `w_rld` and `d_min_ti_n`. Underscore `_` is treated as the wildcard character and is replaced with different combinations of `aoe`.

`h_ll_` contains two wildcard characters, which means badger will combine `aoe` in all possible combinations using 2 character slots.(3 characters in 2 slots can be combined in 3^2 = 9 possible combinations (ex. `halla`, `hallo`, `halle` etc).

## Usage
```sh
badger <flags> [searchterms]

-alpha
	Alphabetic search (a-z)

-alphanum
	Alphanumeric search (a-z, 0-9)

-numeric
	Numeric search (0-9)

-custom [characterset]
	Custom characterset search (limited to a-z,0-9 and -)
	
-tld [comma-separated list of top-level domains]
	Top-level domains

-delay
	Delay in milliseconds

-help
```