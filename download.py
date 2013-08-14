#!/usr/bin/env python

import os
import sys

import settings

def fetch(url):
    filename = os.path.basename(url)
    path = os.path.join(settings.data_dir, filename)
    if not os.path.isfile(path):
        print "saving %s as %s" % (url, path)
        os.system("wget --quiet --output-document %s %s" % (path, url))

def fetch_all():
    if not os.path.isdir(settings.data_dir):
        os.mkdir(settings.data_dir)
    for lang in settings.languages:
        fetch("http://dumps.wikimedia.org/%swiki/latest/%swiki-latest-externallinks.sql.gz" % (lang, lang))

if __name__ == "__main__":
    fetch_all()
