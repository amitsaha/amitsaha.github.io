import sys
with open(sys.argv[1]) as f:
    newlines = ""
    lines = f.readlines()
    for line in lines:
        if not line.startswith('.. code-include::'):
            newlines += line
        else:
            filepath = line.split('::')[1].lstrip().rstrip("\n")
            filedata = open(filepath).readlines()

            newlines += ".. code::\n\n"
            for line in filedata:
                newlines += "   " + line

print(newlines)
