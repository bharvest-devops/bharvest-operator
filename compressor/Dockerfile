FROM busybox:1.36.1-musl AS busybox-full

FROM alpine:3.19.1

LABEL org.opencontainers.image.source="https://github.com/bharvest-devops"
WORKDIR /bin

# Install ln (for making hard links) and rm (for cleanup) from full busybox image (will be deleted, only needed for image assembly)
COPY --from=busybox-full /bin/ln /bin/rm ./

# Install minimal busybox image as shell binary (will create hardlinks for the rest of the binaries to this data)
COPY --from=busybox-full /bin/busybox /bin/sh

# Add hard links for read-only utils
# Will then only have one copy of the busybox minimal binary file with all utils pointing to the same underlying inode
RUN for b in \
  cat \
  date \
  df \
  du \
  env \
  grep \
  head \
  less \
  ls \
  md5sum \
  pwd \
  sha1sum \
  sha256sum \
  sha3sum \
  sha512sum \
  sleep \
  stty \
  tail \
  tar \
  tee \
  tr \
  watch \
  which \
  ; do ln sh $b; done

#  Remove write utils
RUN rm ln rm

RUN apk add --no-cache lz4

# Create compressor user
RUN addgroup --gid 1025 -S compressor && adduser --uid 1025 -S compressor -G compressor

USER compressor