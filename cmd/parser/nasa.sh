#!/bin/bash
#1995 july and august apache access log dataset
#https://ita.ee.lbl.gov/html/contrib/NASA-HTTP.html
# 3.461.612 log lines, pruned the august dataset to 50.000 records. Took 380 seconds, that's 131 loglines per second
# head -n 50000 NASA_access_log_Aug95 > pruned.txt
# uncompressed size 5,1 Mb, db size 2,5 Mb
time go run . -dbname test11.db -dbpath /home/joris/Downloads/logfileparser_demodata/ -logfilefilter pruned -loglineregex '^(?P<RemoteIP>\S+) (?P<RemoteLogname>\S+) (?P<AuthUser>\S+) \[(?P<Timestamp>[^\]]+)\] "(?P<RequestMethod>\S+)(?:\s+(?P<RequestURI>[^\s"]*)[^"]*)?" (?P<StatusCode>\d+) (?P<BytesSent>\d)' -logpath /home/joris/Downloads/logfileparser_demodata/