#!/bin/bash

./gator reset
echo "Sleeping through database reset..."
sleep 1
echo "...done"
./gator users
./gator register marty
./gator register fred
./gator register wilma
./gator register velma
./gator register declan
./gator login fred
./gator users
./gator feeds
./gator addfeed "Hacker News RSS" "https://hnrss.org/newest"
./gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
./gator addfeed "Lanes Blog" "https://www.wagslane.dev/index.xml"
./gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
./gator feeds
./gator users
./gator agg 1m

