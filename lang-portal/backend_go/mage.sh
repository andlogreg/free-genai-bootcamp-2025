#!/bin/bash

# This script is a wrapper for Mage
# It allows running Mage targets without having to type the full command

# Check if the command is provided
if [ $# -eq 0 ]; then
    # No arguments, run the default target
    go run github.com/magefile/mage
else
    # Run the specified target
    go run github.com/magefile/mage "$@"
fi
