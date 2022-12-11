# https://github.com/markdownlint/markdownlint/blob/master/docs/RULES.md
#
# When updating rules in this file, ensure the corresponding rule list in
# .markdownlint.json is also updated.

# style
rule 'header-style'
rule 'hr-style'

# whitespace rules
rule 'no-missing-space-atx'
rule 'no-multiple-space-atx'
rule 'no-missing-space-closed-atx'
rule 'no-multiple-space-closed-atx'
rule 'no-space-in-emphasis'
rule 'no-space-in-code'
rule 'no-space-in-links'

# miscellaneous
rule 'ol-prefix', :style => :ordered
rule 'no-reversed-links'
