# My attempt at making a small-ish image that is capable
# of running `youtube-dl` and `gsutil` along with the python libs
# for redis and pubsub

# TODO: It would be smart to pin the youtube-dl being used
# the current distributed version is a little old
# and prone to breaking, but I'd like to avoid dev breakages too 

FROM google/cloud-sdk:alpine

RUN apk add -q --no-cache git zip make pandoc && \
    git clone https://github.com/ytdl-org/youtube-dl.git && \
    cd youtube-dl && \
    make install && \
    cd .. && \
    rm -rf youtube-dl && \
    apk del -q git zip make pandoc

WORKDIR /usr/src/app
COPY requirements.txt .
RUN python -m ensurepip --upgrade && \
    python -m pip install --no-cache-dir -r requirements.txt && \
    python -m pip uninstall -y pip

COPY subscriber.py .
ENV PYTHONUNBUFFERED=1
CMD ["python", "subscriber.py"]
