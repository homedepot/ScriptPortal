#!/usr/bin/env python3
import sys 
def openAndPrint(filename):
    with open(filename, 'r') as filehandler:
        for line in filehandler:
            print(line,file=sys.stdout,flush=True)

for x in range(len(sys.argv)):
    if sys.argv[x] == "-f":
        openAndPrint(sys.argv[x+1])
