#!/bin/bash

# Rename the file
mv example-config.json config.json

# Run go mod tidy
go mod tidy
