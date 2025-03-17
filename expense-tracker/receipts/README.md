# Bill Parser API

A RESTful API that extracts item details and totals from receipt images using OpenAI's Vision API.

## Features

- Accepts JPG and PNG images of receipts/bills
- Extracts item names, quantities, prices, and totals
- Returns data in structured JSON format
- Includes currency, date, and merchant information

## Requirements

- Python 3.10+
- OpenAI API key

## Setup

1. Clone the repository
2. Install dependencies:
   ```
   pip install -r requirements.txt
   ```
3. Create a `.env` file based on `.env.example` and add your OpenAI API key:
   ```
   CHOREO_OPEN_AI_CONNECTION_OPENAI_API_KEY=your_openai_api_key_here
   PORT=8080
   ```

## Running the API

### Local Development

```bash
python app.py
```

### Using Docker

```bash
docker build -t bill-parser .
docker run -p 8080:8080 --env-file .env bill-parser
```

## API Usage

Send a POST request to `/parse-bill` with a form-data body containing the receipt image:

```bash
curl -X POST http://localhost:8080/parse-bill \
  -F "image=@/path/to/receipt.jpg" \
  -H "Content-Type: multipart/form-data"
```

### Example Response

```json
{
  "items": [
    {
      "name": "Milk",
      "quantity": 1,
      "price": 3.99
    },
    {
      "name": "Bread",
      "quantity": 2,
      "price": 2.49
    }
  ],
  "total": 8.97,
  "currency": "USD",
  "date": "2023-05-15",
  "merchant": "Grocery Store Inc."
}
```

## Testing

### Automated Tests

Run the unit tests with:

```bash
python -m unittest test_app.py
```

The tests use mocked OpenAI responses to verify the API functionality without making actual API calls.

### Manual Testing

We've included a convenient script for manual testing with real receipt images:

1. Make sure the API is running (either locally or in Docker)
2. Use the manual_test.py script to test with a real receipt image:

```bash
# Install the required package for the test script
pip install requests

# Test with a local receipt image
python manual_test.py /path/to/your/receipt.jpg

# If testing with a remote API endpoint
python manual_test.py /path/to/your/receipt.jpg --url https://your-api-url.com/parse-bill
```

#### Sample Test Images

For testing purposes, you can use receipts from:
- Your own shopping receipts
- Sample receipt images from the internet
- Create a synthetic test receipt using the included generator:

```bash
# Install Pillow if not already installed
pip install pillow

# Generate a sample receipt image
python create_test_image.py

# The image will be saved as sample_receipt.jpg
# Now test it with the manual test script
python manual_test.py sample_receipt.jpg
```

### Test Image Preparation Tips

For best results with the parser:
1. Ensure the receipt image is well-lit and clear
2. The image should capture the entire receipt including all items and the total
3. Avoid glare, shadows, and crumpled receipts when possible

## OpenAPI Specification

The API is documented using OpenAPI 3.0. View the specification in the `openapi.yaml` file.