###############################################################################
# Step 1 : Builder image
#
FROM golang:1.17.6-alpine3.15 AS builder

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
#ENV CGO_ENABLED=1

WORKDIR /app
COPY . .

RUN go build -o="./out/app_program"

############################
# STEP 2 build a production image
############################
FROM alpine:3.13.5

#Install dependencies
RUN apk --no-cache add ghostscript 

COPY --from=builder --chown=app:app /app/out/app_program /app/app_program
COPY --from=builder --chown=app:app /app/scripts/shrinkpdf.sh /app/scripts/shrinkpdf.sh

WORKDIR /app
RUN addgroup -S app && adduser -S app -G app && chown -R app:app /app
RUN chown -R app:app /app
USER app

EXPOSE 3000

# Command to run the executable
CMD ["./app_program"]