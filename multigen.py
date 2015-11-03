#!/usr/bin/python

import sys
import os
import random

HIGH = 0
MED = 1
LOW = 2

keyString = sys.argv[1]
numStudents = int(sys.argv[2])

# Turn the key string into an array of ints
key = map(int, keyString.split(","))

for sid in range(1, numStudents+1):
    performs = random.randint(0, 2)
    result = []
    for q in range(0, len(key)):
        willGuess = False

        r = random.random()
        if performs == HIGH:
            if r > 0.9:
                willGuess = True
        elif performs == MED:
            if r > 0.7:
                willGuess = True
        elif performs == LOW:
            if r > 0.5:
                willGuess = True

        if willGuess:
            gint = random.randint(0,4)
            g = 1 << gint 
            result.append(g)
        else:
            result.append(key[q])

    sidFlag = "-sid "+`sid`
    dataFlag = "-data "+(",".join(map(str, result)))
    outFlag = "-out ./s"+`sid`+".png"
    flags = sidFlag+" "+dataFlag+" "+outFlag
    print(flags)
    os.system("./gdformgen "+flags)
