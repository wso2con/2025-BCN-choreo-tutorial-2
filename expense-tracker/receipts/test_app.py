import unittest
import os
import json
from unittest.mock import patch, MagicMock
from io import BytesIO

from app import app

class BillParserTestCase(unittest.TestCase):
    def setUp(self):
        app.config['TESTING'] = True
        self.client = app.test_client()
        
        # Create a test image (small blank JPG)
        self.test_image = BytesIO(
            b'\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x01\x00H\x00H\x00\x00\xff\xdb\x00C\x00\x08\x06\x06\x07\x06\x05\x08\x07\x07\x07\t\t\x08\n\x0c\x14\r\x0c\x0b\x0b\x0c\x19\x12\x13\x0f\x14\x1d\x1a\x1f\x1e\x1d\x1a\x1c\x1c $.\' ",#\x1c\x1c(7),01444\x1f\'9=82<.342\xff\xdb\x00C\x01\t\t\t\x0c\x0b\x0c\x18\r\r\x182!\x1c!22222222222222222222222222222222222222222222222222\xff\xc0\x00\x11\x08\x00\x01\x00\x01\x03\x01"\x00\x02\x11\x01\x03\x11\x01\xff\xc4\x00\x1f\x00\x00\x01\x05\x01\x01\x01\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\xff\xc4\x00\xb5\x10\x00\x02\x01\x03\x03\x02\x04\x03\x05\x05\x04\x04\x00\x00\x01}\x01\x02\x03\x00\x04\x11\x05\x12!1A\x06\x13Qa\x07"q\x142\x81\x91\xa1\x08#B\xb1\xc1\x15R\xd1\xf0$3br\x82\t\n\x16\x17\x18\x19\x1a%&\'()*456789:CDEFGHIJSTUVWXYZcdefghijstuvwxyz\x83\x84\x85\x86\x87\x88\x89\x8a\x92\x93\x94\x95\x96\x97\x98\x99\x9a\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xff\xc4\x00\x1f\x01\x00\x03\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x00\x00\x00\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\xff\xc4\x00\xb5\x11\x00\x02\x01\x02\x04\x04\x03\x04\x07\x05\x04\x04\x00\x01\x02w\x00\x01\x02\x03\x11\x04\x05!1\x06\x12AQ\x07aq\x13"2\x81\x08\x14B\x91\xa1\xb1\xc1\t#3R\xf0\x15br\xd1\n\x16$4\xe1%\xf1\x17\x18\x19\x1a&\'()*56789:CDEFGHIJSTUVWXYZcdefghijstuvwxyz\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x92\x93\x94\x95\x96\x97\x98\x99\x9a\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xff\xda\x00\x0c\x03\x01\x00\x02\x11\x03\x11\x00?\x00\xfe\xfe(\xa2\x8a\x00\xff\xd9'
        )
        self.test_image.name = 'test_receipt.jpg'
        
        # Sample response from OpenAI
        self.mock_response = MagicMock()
        self.mock_response.choices = [
            MagicMock(
                message=MagicMock(
                    content=json.dumps({
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
                    })
                )
            )
        ]

    @patch('app.client.chat.completions.create')
    def test_parse_bill_success(self, mock_create):
        # Configure the mock to return a predefined response
        mock_create.return_value = self.mock_response
        
        # Send request with test image
        response = self.client.post(
            '/parse-bill',
            data={'image': (self.test_image, 'test_receipt.jpg')},
            content_type='multipart/form-data'
        )
        
        # Assert response
        self.assertEqual(response.status_code, 200)
        data = json.loads(response.data)
        
        # Verify structure and content
        self.assertIn('items', data)
        self.assertEqual(len(data['items']), 2)
        self.assertEqual(data['items'][0]['name'], 'Milk')
        self.assertEqual(data['total'], 8.97)
        self.assertEqual(data['currency'], 'USD')
        
        # Verify OpenAI was called with correct parameters
        mock_create.assert_called_once()
        
    def test_parse_bill_no_image(self):
        # Test with no image
        response = self.client.post(
            '/parse-bill',
            data={},
            content_type='multipart/form-data'
        )
        
        # Assert error response
        self.assertEqual(response.status_code, 400)
        data = json.loads(response.data)
        self.assertIn('error', data)
        self.assertEqual(data['error'], 'No image provided')
    
    def test_parse_bill_invalid_extension(self):
        # Test with invalid file extension
        invalid_file = BytesIO(b'invalid file content')
        invalid_file.name = 'test.txt'
        
        response = self.client.post(
            '/parse-bill',
            data={'image': (invalid_file, 'test.txt')},
            content_type='multipart/form-data'
        )
        
        # Assert error response
        self.assertEqual(response.status_code, 400)
        data = json.loads(response.data)
        self.assertIn('error', data)
        self.assertEqual(data['error'], 'File format not supported. Please upload JPG or PNG')

if __name__ == '__main__':
    unittest.main()