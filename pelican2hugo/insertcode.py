import sys
with open(sys.argv[1]) as f:
    newlines = ""
    lines = f.readlines()
    for line in lines:
        if line.startswith('.. code-include::'):
            filepath = line.split('::')[1].lstrip().rstrip("\n")
            filedata = open(filepath).readlines()
            newlines += ".. code::\n\n"
            for l in filedata:
                newlines += "   " + l
        elif ':lexer:' in line.strip():
            continue
        else:
            newlines += line
print(newlines)
