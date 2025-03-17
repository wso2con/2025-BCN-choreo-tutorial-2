import requests
import argparse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

def parse_receipt(image_path, url=None):
    """
    Send a receipt image to the bill parser API and display the result.
    
    Args:
        image_path: Path to the receipt image file (jpg or png)
        url: API URL (defaults to localhost)
    """
    if not os.path.exists(image_path):
        print(f"Error: File {image_path} not found")
        return

    if url is None:
        url = f"http://localhost:{os.environ.get('PORT', 8080)}/parse-bill"
    
    # Prepare the file for upload
    with open(image_path, 'rb') as f:
        files = {'image': (os.path.basename(image_path), f, 'image/jpeg' if image_path.endswith(('.jpg', '.jpeg')) else 'image/png')}
        
        # Send the request
        print(f"Sending image {image_path} to {url}...")
        response = requests.post(url, files=files)
    
    # Process the response
    if response.status_code == 200:
        print("Success! Parsed receipt details:")
        result = response.json()
        
        # Display items
        print("\nItems:")
        for item in result.get('items', []):
            print(f"  - {item.get('name', 'Unknown')}: {item.get('quantity', 0)} x ${item.get('price', 0):.2f}")
        
        # Display total and other information
        print(f"\nTotal: ${result.get('total', 0):.2f} {result.get('currency', '')}")
        print(f"Date: {result.get('date', 'Unknown')}")
        print(f"Merchant: {result.get('merchant', 'Unknown')}")
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Test the Bill Parser API with a receipt image")
    parser.add_argument("image_path", help="Path to the receipt image file (jpg or png)")
    parser.add_argument("--url", help="API URL (defaults to localhost)")
    
    args = parser.parse_args()
    parse_receipt(args.image_path, args.url)