#!/bin/sh
app="gql-server"

printf "\nStart running: $app\n"

# Set all ENV vars for the server to run
export $(grep -v '^#' .env | xargs)
time realize start run

printf "\nStopped running: $app\n\n"
