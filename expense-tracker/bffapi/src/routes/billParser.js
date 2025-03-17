const express = require('express');
const multer = require('multer');
const FormData = require('form-data');
const { billParserClient, accountsClient } = require('../utils/serviceClient');

const parserRouter = express.Router();

// Configure multer for file uploads
const storage = multer.memoryStorage();
const upload = multer({ 
  storage,
  limits: { fileSize: 5 * 1024 * 1024 } // Limit file size to 5MB
});

// Combined parser endpoint that can both parse and optionally create a bill
parserRouter.post('/parse', upload.single('image'), async (req, res) => {
  try {
    if (!req.file) {
      return res.status(400).json({ error: 'No image file provided' });
    }

    // Create form data with the image file
    const formData = new FormData();
    formData.append('image', req.file.buffer, {
      filename: req.file.originalname,
      contentType: req.file.mimetype
    });

    // Parse the receipt using the bill parser service
    const parseResponse = await billParserClient.post('/parse-bill', formData, {
      headers: {
        ...formData.getHeaders(),
        'Content-Type': 'multipart/form-data'
      }
    });

    const parsedData = parseResponse.data;

    // If create_bill parameter is true, create a bill from the parsed data
    if (req.body.create_bill === 'true') {
      // Generate a title if not provided
      const title = req.body.title || 
                    (parsedData.merchant ? 
                      `Bill from ${parsedData.merchant}` : 
                      `Bill for $${parsedData.total}`);

      // Convert parsed receipt data to bill input format
      const billData = {
        title,
        total: parsedData.total,
        due_date: parsedData.date || new Date().toISOString().split('T')[0],
        paid: false,
        items: parsedData.items.map(item => ({
          name: item.name,
          amount: item.price,
          quantity: item.quantity
        }))
      };

      // Create the bill using the accounts service
      const billResponse = await accountsClient.post('/bills', billData);

      // Return combined response
      return res.json({
        message: 'Bill created successfully',
        billId: billResponse.data.id,
        parsedData
      });
    }

    // If not creating a bill, just return the parsed data
    return res.json(parsedData);
  } catch (error) {
    console.error('Error parsing bill:', error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to parse bill' }
    );
  }
});

module.exports = { parserRouter };
