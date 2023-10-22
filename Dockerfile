# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.21-alpine AS build

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .

# compile application
RUN go build -o /shop

##
## STEP 2 - DEPLOY
##
FROM scratch

WORKDIR /

COPY --from=build /shop /shop

EXPOSE 3333

ENTRYPOINT ["/shop"]