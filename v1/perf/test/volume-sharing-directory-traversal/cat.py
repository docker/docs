#!//usr/bin/env python

from optparse import OptionParser

parser = OptionParser()
parser.add_option("-d", "--depth", dest="depth", type='int')
parser.add_option("-w", "--width", dest="width", type='int')

def parse_args():
    (options, _) = parser.parse_args()
    w, d = options.width, options.depth
    return (int(w) if w is not None else 1,
            int(d) if d is not None else 1)

def generate(width, depth):
    sexp = '(base ' + ' '.join('w-%05d' % i for i in range(1, width))  + ')'
    prefix = ' '.join('(s-%d ' % i for i in range(1, depth))
    suffix = ' '.join(')' for _ in range(1, depth))
    return prefix + sexp + suffix

if __name__ == '__main__':
    width, depth = parse_args()
    print generate(width, depth)

