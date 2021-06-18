#!/bin/bash

for openApi in `ls`;
do
  if [[ -f "./$openApi/$openApi-swagger.yaml" ]]; then
     echo "****** 📗 Getting OpenApi Doc 📗 ******"
     echo $openApi
     cd "./$openApi" && cp "$openApi-swagger.yaml" "../ploutarchos/yaml/"
     echo "****** 📋 Copied OpenApi Doc 📋 ******"
     cd ..
  fi
done

for buildableBlock in `ls`;
do
  if [[ -f "./$buildableBlock/Makefile" ]]; then
     echo "****** 🏗️ Building in process 🏗️  ******"
     echo $buildableBlock
     cd "./$buildableBlock" && make build
     echo "****** 📦 Building complete 📦 ******"
     cd ..
  fi
  if [[ -f "./$buildableBlock/Dockerfile" ]]; then
     echo "****** 🚢 Docker building 🚢 ******"
     echo $buildableBlock
     cd "./$buildableBlock" && make docker-build
     echo "****** 🔱 Container done 🔱 ******"
     cd ..
  fi
done