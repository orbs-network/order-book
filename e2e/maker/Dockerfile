FROM python:3.9-slim

WORKDIR /app

COPY /e2e/maker /app

RUN pip install --no-cache-dir -r ./requirements.txt

EXPOSE 80

CMD ["python", "main.py"]
