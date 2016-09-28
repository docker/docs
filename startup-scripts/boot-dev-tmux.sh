#!/bin/bash

set -e

eval $(docker-machine env dev)

###############################################################
# You must have `tmux` installed locally. On OSX, this can be #
# accomplished with `brew install tmux`                       #
###############################################################

SESSION=HubDev
DIR=${PWD##*/}_hub_1
CONTAINER=$(sed s/-//g <<< $DIR)


# Create new tmux session
tmux -2 new-session -d -s $SESSION

# Window 1

## webpack task
tmux split-window
tmux select-pane -t 0
tmux send-keys "DEBUG=* webpack -wd" C-m

## styles

tmux select-pane -t 1
tmux send-keys "DEBUG=* gulp watch::styles::dev" C-m

## Flow

tmux split-window -h
tmux select-pane -t 2
tmux send-keys "flow" C-m

## docker logs

tmux select-pane -t 0
tmux split-window -h
tmux send-keys "docker-compose logs hub" C-m

# Attach to session
tmux -2 attach-session -t $SESSION
