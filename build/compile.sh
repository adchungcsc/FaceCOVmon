#!/bin/bash
cd ../backend/api/function
rm handler
rm handler.zip
go build -o handler
zip handler.zip handler

