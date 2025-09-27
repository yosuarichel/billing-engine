# Stage 1: runtime
FROM alpine:3.20
WORKDIR /app/billing_engine
# COPY --from=builder /app/billing_engine .
COPY ./output .
# EXPOSE 8080
RUN chmod +x bootstrap.sh
CMD ["sh", "./bootstrap.sh"]
