#!/bin/bash
./server \
-l server.log \
-a 0.0.0.0:3000 \
-jwt secret \
-dbn minesweeper \
-dbh localhost \
-dbu minesweeper \
-dbp minesweeper \