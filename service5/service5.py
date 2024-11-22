#!/usr/bin/env python
import pika
import os

# Connects to the RabbitMQ server
username = os.getenv('USER', 'guest')
pw = os.getenv('PASS', 'guest')
host = os.getenv('RABBITMQ_SVC_NAME', 'localhost')

credentials = pika.PlainCredentials(username, pw)
connection = pika.BlockingConnection(
    pika.ConnectionParameters(host,
                                5672,
                                '/',
                                credentials))
channel = connection.channel()

channel.exchange_declare(exchange='siblingExchange', exchange_type='fanout')

result = channel.queue_declare(queue='', exclusive=True)
queue_name = result.method.queue

channel.queue_bind(exchange='siblingExchange', queue=queue_name)

print(' [*] Waiting for logs. To exit press CTRL+C')

def callback(ch, method, properties, body):
    print("Received a message from Service 4")
    print(f" [x] {body}")

channel.basic_consume(
    queue=queue_name, on_message_callback=callback, auto_ack=True)

channel.start_consuming()