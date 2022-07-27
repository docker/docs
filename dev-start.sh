#!/bin/bash

if [ "$JEKYLL_ENV" == "production" ]; then
  echo "Starting in production mode"
  exec bundle exec jekyll serve --host=0.0.0.0 -l --config _config.yml,_config_production.yml
else
  exec bundle exec jekyll serve --host=0.0.0.0 -l 
fi
