#!/bin/bash

echo "Validating env" 
# verify user has OPENAI_TOKEN set
if [ -z "$OPENAI_API_KEY" ]
then
	  echo "\$OPENAI_API_KEY is empty. Please set this to your OpenAI API token."
	  exit 1
fi

# verify user has Go installed and that it is at least version 1.18
if ! command -v go &> /dev/null
then
	  echo "Go could not be found"
	  exit 1
fi


echo "Installing rft"
go install github.com/richardfaulkner-qz/rft

# install the auto complete script, based on if you're using bash or zsh

echo "rft successfully installed"

