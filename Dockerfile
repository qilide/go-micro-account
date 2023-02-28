FROM alpine
ADD account-service /account-service
ENTRYPOINT [ "/account-service" ]
