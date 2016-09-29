#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Constants for colour output"""

# codes for colour output
COLOUR_RESET = "\033[0m"
COLOUR_SEQ = "\033[1;%dm"
COLOUR_BOLD = "\033[1m"

COLOUR_BLACK = COLOUR_SEQ % (30)
COLOUR_RED = COLOUR_SEQ % (31)
COLOUR_GREEN = COLOUR_SEQ % (32)
COLOUR_YELLOW = COLOUR_SEQ % (33)
COLOUR_MAGENTA = COLOUR_SEQ % (35)
COLOUR_GREY = COLOUR_SEQ % (37)
COLOUR_FOREGROUND = COLOUR_SEQ % (39)
