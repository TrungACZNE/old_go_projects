#!/usr/bin/env python
import os
import sys
import time

for n in range(6):
    time.sleep(1)
    if n % 2 == 0:
        sys.stdout.write("STDOUT " + str(n) + "\n")
        sys.stdout.flush()
    else:
        sys.stderr.write("STDERR " + str(n) + "\n")
        sys.stderr.flush()
