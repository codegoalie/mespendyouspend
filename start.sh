#!/bin/bash

set -e

. ~/.nvm/nvm.sh

rm -f public/assets/*.js
rm -f public/assets/*.css

nvm use
docker-compose start
buffalo dev
