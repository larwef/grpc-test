FROM alpine:3.8

RUN apk --no-cache add bash

ENV USER=appuser
ENV HOME=/home/$USER

RUN adduser -S $USER

# Copy files to container
COPY build/docker/toRoot/* $HOME/
COPY target/linux/server $HOME/app

# Configure permissions
RUN chmod 755 $HOME/*.sh && \
    chown -R $USER $HOME

USER $USER
WORKDIR /home/$USER
CMD ["./runapp.sh"]