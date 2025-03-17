FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Copy requirements first for better caching
COPY requirements.txt .

# Install dependencies with pinned versions
RUN pip install --no-cache-dir -r requirements.txt

# Create a non-root user for security
RUN groupadd -g 10014 choreouser && \
    useradd -u 10014 -g choreouser -s /bin/bash -m choreouser

# Copy application code
COPY . .

# Set environment variables
ENV PORT=8080
ENV PYTHONUNBUFFERED=1

# Set permissions for the non-root user
RUN chown -R choreouser:choreouser /app

# Switch to non-root user
USER 10014

# Expose the port
EXPOSE 8080

# Use gunicorn as the production server
CMD ["gunicorn", "--bind", "0.0.0.0:8080", "app:app"]
