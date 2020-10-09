#!/bin/bash
protoc --go_out=plugins=grpc:./pb  pb/oauth.proto pb/user.proto