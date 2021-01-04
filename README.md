# Domain search - badger

This is an experimental tool used to find unregistered domains using wildcards & characterset combinations.


# Usage
```
	badger <flags> [searchterms]

	-alpha
		Alphabetic search (a-z)

	-alphanum
		Alphanumeric search (a-z, 0-9)

	-numeric
		Numeric search (0-9)

	-custom [character set]
		Custom character set search (limited to a-z,0-9 and -)
		
	-tld [comma-separated list of top-level domains]
		Top-level domains

	-delay
		Delay in milliseconds
```


### Example
```
	badger -custom aoe -tld se,io,nu h_ll_ w_rld d_min_ti_n
```

This will perform a whois lookup on domain name combinations using selected tlds and combinations of the characters `aoe` in three separate search terms `h_ll_`, `w_rld` and `d_min_ti_n`, where underscore is the wildcard character.

`h_ll_` contains two wildcard characters, which means badger will combine `aoe` in all possible combinations using 2 character slots. 3 characters in 2 slots can be combined in 3^2 = 9 ways which isnt a big deal to manually.

But, when you're looking at multiple wildcards using a wider range of characters it quickly gets tedious.


### Parameters


|Parameter|Description|Example|
|-|-|-|
|alpha|Use alphabetic characters (a-z)|```badger -alpha```|
|alphanum|Use alphanumeric characters (a-z, 0-9)|```badger -alphanum```|
|numeric|Use numeric characters (0-9)|```badger -numeric```|
|custom|Use custom characters|```badger -custom abc123```|
|tld|List with top-level domains should be searched|```badger -tld se,dk,no```|
|delay|Delay in milliseconds between whois lookups|```badger -delay 500```|