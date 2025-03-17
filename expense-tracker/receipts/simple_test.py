import requests
import sys
import os

def test_api(image_path):
    """Simple test script to check if the API is working"""
    url = f"http://localhost:8080/parse-bill"
    
    if not os.path.exists(image_path):
        print(f"Error: File {image_path} not found")
        return
    
    # Prepare the file for upload
    with open(image_path, 'rb') as f:
        files = {'image': (os.path.basename(image_path), f, 'image/jpeg')}
        
        # Send the request
        print(f"Sending image {image_path} to {url}...")
        response = requests.post(url, files=files)
    
    # Print the raw response
    print(f"Status code: {response.status_code}")
    print(f"Response content:")
    print(response.text)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python simple_test.py <image_path>")
        sys.exit(1)
    
    test_api(sys.argv[1])