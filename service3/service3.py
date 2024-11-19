#!/usr/bin/env python
import pika
import os 
import sys
import json

# This function sends messages to the next service
def sendMessage(message):
	# Connects to the RabbitMQ server
    host = os.getenv('RABBITMQ_SVC_NAME', 'localhost')
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host=host))
    channel = connection.channel()

    #  Declares the queue
    channel.queue_declare(queue='queue1')

    # Parse JSON string to Python object
    json_string = message.decode('utf-8')
    data = json.loads(json_string)

    # Edits the message
    data['History'].append("Message was received to Service 3")
    data["History"].append("Message was sent to Service 4")
    data["Body"] = "Sent from Python Service 3"

    # Sends the message to next service
    channel.basic_publish(exchange='', routing_key='queue1', body=json.dumps(data))
    print(" [x] Sent JSON payload to Service 4'")

    connection.close()

def main():
    # Connects to the RabbitMQ server
    host = os.getenv('RABBITMQ_SVC_NAME', 'localhost')
    connection = pika.BlockingConnection(pika.ConnectionParameters(host=host))
    channel = connection.channel()

    # Declares the queue
    channel.queue_declare(queue='pythonQueue')

    # Callback method for receiving messages
    def callback(ch, method, properties, body):
        print("Received a message from Service 2")
        sendMessage(body)

    # Consumes the messages
    channel.basic_consume(queue='pythonQueue', on_message_callback=callback, auto_ack=True)

    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)