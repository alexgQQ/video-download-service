FROM python:3.10-slim

WORKDIR /usr/src/app

RUN pip install --no-cache-dir google-cloud-pubsub
COPY create.py .

ENV PYTHONUNBUFFERED=1

CMD ["python", "create.py"]
