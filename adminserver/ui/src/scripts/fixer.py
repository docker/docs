#!/usr/bin/env python

import json
import subprocess

def fix(line, column, p, l):
  if p == 'strict':
    return "'use strict';\n\n" + l

  if p == 'quotes' or p == 'react/jsx-quotes': 
    assert l[column] == '"'
    end = l.find('"', column + 1)
    while end != -1 and l[end-1:end] == '\\' and l[end - 2:end] != '\\\\': # the quote we found was escaped
      end = l.find('"', end + 1)
    assert end != -1
    quoteString = l[column + 1:end]
    assert "\\'" not in quoteString
    quoteString = quoteString.replace("'", "\\'")
    quoteString = quoteString.replace('\\"', '"')
    return l[:column] + "'" + quoteString + "'" + l[end + 1:] 

  if p == 'no-trailing-spaces': 
    return l.rstrip() + '\n'

  if p == 'eqeqeq': 
    assert l[column:column + 2] in ('==', '!=') and l[column + 2] != '='
    return l[:column + 1] + '=' + l[column + 1:]

  if p == 'semi': 
    return l[:column] + ';' + l[column:]

  if p == 'no-extra-semi': 
    assert l[column] == ';'
    return l[:column] + l[column + 1:] 

  if p == 'key-spacing':
    assert l[column - 1] == ' '
    return l[:column - 1] + l[column:] 

  if p == 'comma-dangle': 
    assert l[column] == ','
    return l[:column] + l[column + 1:] 

  if p == 'space-infix-ops':
    assert l[column - 1] != ' '
    assert l[column] != ' '
    assert l[column + 1] != ' '
    return l[:column] + ' ' + l[column:column + 1] + ' ' + l[column + 1:] 

  if p == 'space-unary-ops': 
    if l.startswith('! ', column):
      return l[:column + 1] + l[column + 2:] 
    elif l.startswith('delete', column) and l[column + len('delete')] != ' ':
      return l[:column + len('delete')] + ' ' + l[column + len('delete'):]
    else:
      assert False

  if p == 'react/jsx-curly-spacing':
    if l[column] == '{':
      return l[:column + 1] + ' ' + l[column + 1:] 
    elif l[column] == '}':
      return l[:column] + ' ' + l[column:] 
    else:
      assert False

  return l

if __name__ == "__main__":
  '''
  Need to install eslint/babel-eslint globally first.
  One-off script to automatically fix some of the common existing eslint errors.
  It's not smart about how it fixes the errors so you need to verify the output manually.
  Highly destructive! Make a commit before trying it out.
  '''

  while True: # Loop ends when it can't fix anything else
    print 'Linting...'
    output = subprocess.Popen('eslint --format json .', shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE).communicate()[0]
    files = json.loads(output)
    totalErrors = 0
    totalFixed = 0
    for f in files:
      path = f['filePath']
      totalErrors += f['errorCount']

      problems = {} # gather by line and column
      for message in f['messages']:
        line = message['line'] - 1 # zero indexed
        column = message['column'] - 1
        if line not in problems:
            problems[line] = {}
        if column not in problems[line]:
          problems[line][column] = []
        problems[line][column].append(message)

      fixed = 0
      newFile = []
      with open(path, 'r') as readFile:
        for line, text in enumerate(readFile):
          if line in problems:
            currText = str(text)
            for column in problems[line].iterkeys():
              problem = problems[line][column][0]# Skip if there are multiple problems at same column
              fixedText = fix(line, column, problem['ruleId'], currText)
              if currText != fixedText:
                fixed += 1
                print path, '%d:%d,' % (line, column), problem['ruleId'], problem['message']
                print currText[:-1]
                print fixedText[:-1]
                print
              if len(currText) != len(fixedText):
                currText = fixedText
                break # Skip fixing the rest of the problems on this line if the text got shifted around
              else:
                currText = fixedText

            newFile.append(currText)
          else:
            newFile.append(text)
      totalFixed += fixed

      with open(path, 'w') as writeFile:
        writeFile.write(''.join(newFile))
    print totalFixed, 'fixed out of', totalErrors, 'remaining'
    if not totalFixed:
      break
