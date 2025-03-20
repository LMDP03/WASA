#!/bin/bash

# Esegui i container contemporaneamente
sudo docker run --rm -p 3000:3000 backend:latest &
sudo docker run --rm -p 8080:80 frontend:latest &

# Attendi la terminazione dei processi in background
wait
