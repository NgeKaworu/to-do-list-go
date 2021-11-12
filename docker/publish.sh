#!/bin/bash
set -e

tag=ngekaworu/todo-list-go

docker build --file ./Dockerfile --tag ${tag} ..;
docker push ${tag};
