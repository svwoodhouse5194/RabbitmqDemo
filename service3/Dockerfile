# Use an official Python runtime as a parent image
FROM python:3.11

# This is to enable logs
ENV PYTHONUNBUFFERED=1 

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY service3.py /app

# Install any needed packages specified in requirements.txt
RUN pip install pika --upgrade

# Run app.py when the container launches
CMD ["python", "service3.py"]
