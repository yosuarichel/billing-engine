# Stage 1: runtime
FROM alpine:3.20
WORKDIR /app/billing-engine
COPY output/ .
RUN ls -lah .
RUN chmod +x bootstrap.sh
CMD ["sh", "./bootstrap.sh"]
