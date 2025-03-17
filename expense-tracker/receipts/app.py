import os
import base64
from flask import Flask, request, jsonify
from openai import OpenAI
import tempfile
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

app = Flask(__name__)

# Initialize OpenAI client
api_key = os.environ.get("CHOREO_OPEN_AI_CONNECTION_OPENAI_API_KEY")
if not api_key:
    raise ValueError("CHOREO_OPEN_AI_CONNECTION_OPENAI_API_KEY environment variable is not set")
client = OpenAI(api_key=api_key)

@app.route('/parse-bill', methods=['POST'])
def parse_bill():
    if 'image' not in request.files:
        return jsonify({"error": "No image provided"}), 400

    file = request.files['image']

    # Check if the file is allowed
    allowed_extensions = {'jpg', 'jpeg', 'png'}
    if '.' not in file.filename or file.filename.rsplit('.', 1)[1].lower() not in allowed_extensions:
        return jsonify({"error": "File format not supported. Please upload JPG or PNG"}), 400

    # Save the uploaded file temporarily
    temp_file = tempfile.NamedTemporaryFile(delete=False)
    file.save(temp_file.name)
    temp_file.close()

    try:
        # Convert image to base64 for OpenAI
        with open(temp_file.name, "rb") as image_file:
            base64_image = base64.b64encode(image_file.read()).decode('utf-8')

        # Call OpenAI Vision API with prompt
        response = client.chat.completions.create(
            model="gpt-4o",  # Using GPT-4o which supports vision
            messages=[
                {
                    "role": "system",
                    "content": "You are a specialized receipt parser. Extract all item details from receipts including item names, quantities, individual prices, and the total amount. You must respond ONLY with a valid, parseable JSON object with no additional text or formatting. Use this exact structure: {\"items\": [{\"name\": string, \"quantity\": number, \"price\": number}], \"total\": number, \"currency\": string, \"date\": string, \"merchant\": string}"
                },
                {
                    "role": "user",
                    "content": [
                        {
                            "type": "text",
                            "text": "Parse this receipt image and extract all item details (name, quantity, price), the total amount, currency, date, and merchant name. Return ONLY valid JSON with no markdown formatting or additional text."
                        },
                        {
                            "type": "image_url",
                            "image_url": {
                                "url": f"data:image/jpeg;base64,{base64_image}"
                            }
                        }
                    ]
                }
            ],
            max_tokens=1000,
            response_format={"type": "json_object"}  # Force JSON response format
        )

        # Extract the parsed result
        parsed_result = response.choices[0].message.content.strip()

        # Parse the string result into a Python dictionary
        import json
        try:
            parsed_json = json.loads(parsed_result)
            # Return the properly parsed JSON
            return jsonify(parsed_json), 200
        except json.JSONDecodeError as json_error:
            return jsonify({
                "error": "Failed to parse OpenAI response as JSON",
                "raw_response": parsed_result,
                "json_error": str(json_error)
            }), 500

    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        # Clean up the temporary file
        os.unlink(temp_file.name)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))