import sys
print("the file was called as the following:")
print(*sys.argv, file=sys.stderr, flush=True)
print(*sys.argv, file=sys.stdout, flush=True)
