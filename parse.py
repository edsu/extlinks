#!/usr/bin/env python

import os
import re
import gzip
import codecs

import settings

def load_all():
    for lang in settings.languages:
        filename = "%swiki-latest-externallinks.sql.gz" % lang
        path = os.path.join(settings.data_dir, filename)
        load_links_dump(path, lang)

def load_links_dump(filename, lang):
    pattern = r"\((\d+),'(.+?)','(.+?)'\)"
    parse_sql(filename, pattern, process_externallink_row, lang)

def process_externallink_row(row, lang):
    article_id, url, reversed_url = row
    print "\t".join([lang, article_id, url])

def parse_sql(filename, pattern, func, lang):
    "Rips columns specified with the regex out of the sql"
    if filename.endswith('.gz'):
        fh = codecs.EncodedFile(gzip.open(filename), data_encoding="utf-8")
    else:
        fh = codecs.open(filename, encoding="utf-8")

    line = ""
    count = 0
    while True:
        buff = fh.read(1024)
        if not buff:
            break

        line += buff

        rows = list(re.finditer(pattern, line))
        for row in rows:
            try:
                func(row.groups(), lang)
            except Exception, e:
                print "uhoh: %s: %s w/ %s" % (func, e, row.groups())

        # don't forget unconsumed bytes, need them for the next match
        if len(rows) > 0:
            line = line[rows[-1].end():]

if __name__ == "__main__":
    load_all()
