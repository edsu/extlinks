extlinks
========

extlinks is a small utility for downloading all the [wikipedia external link dump files](http://dumps.wikimedia.org/), and extracting the wikipedia page ids and external links from them as a stream suitable for other processing.

Usage
-----

First download all the wikipedia dumps. If you only want specific wikipedia
languages modify the `languages` variable in `settings.py`. You can also
configure what directory the data files will be stored in.

```
./download.py
```

Once the downloads are done you can parse the gzipped mysql dumps, and write 
out the language code, article id and url as tab delimited rows to stdout.

```
./parse.py
```

License
-------

* CC0
