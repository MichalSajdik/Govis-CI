FROM golang:alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN cd /etc/ssl/certs/ \
 && wget "<sap certificate server>/aia/SAP%20Global%20Root%20CA.crt" -O SAP_Global_Root_CA.crt \
 && cat SAP_Global_Root_CA.crt >> ca-certificates.crt \
 && wget "<sap certificate server>/aia/SAP%20Global%20Sub%20CA%2002.crt" -O SAP_Global_Sub_CA_02.crt \
 && cat SAP_Global_Sub_CA_02.crt >> ca-certificates.crt \
 && wget "<sap certificate server>/aia/SAP%20Global%20Sub%20CA%2004.crt" -O SAP_Global_Sub_CA_04.crt \
 && cat SAP_Global_Sub_CA_04.crt >> ca-certificates.crt \
 && wget "<sap certificate server>/aia/SAP%20Global%20Sub%20CA%2005.crt" -O SAP_Global_Sub_CA_05.crt \
 && cat SAP_Global_Sub_CA_05.crt >> ca-certificates.crt \
 && wget "<sap certificate server>/aia/SAPNetCA_G2.crt" -O SAPNetCA_G2.crt \
 && cat SAPNetCA_G2.crt >> ca-certificates.crt


RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
CMD ["./Govis-CI"]

FROM alpine:latest AS alpine
RUN mkdir /app/results -p
RUN apk update && apk add --no-cache git
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/Govis-CI /app/Govis-CI
WORKDIR /app/
EXPOSE 8000
CMD ["./Govis-CI"]