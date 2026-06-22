FROM debian:stable-slim

# COPY source destination
COPY books /bin/books
CMD ["/bin/books"]
