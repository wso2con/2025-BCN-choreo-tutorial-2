const express = require('express');
const { accountsClient } = require('../utils/serviceClient');

const billsRouter = express.Router();

// Get all bills
billsRouter.get('/', async (req, res) => {
  try {
    const response = await accountsClient.get('/bills');
    // If response is null or undefined, return an empty array
    if (!response || !response.data) {
      return res.json([]);
    }
    res.json(response.data);
  } catch (error) {
    console.error('Error fetching bills:', error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to fetch bills' }
    );
  }
});

// Get a bill by ID
billsRouter.get('/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const response = await accountsClient.get(`/bills/${id}`);
    res.json(response.data);
  } catch (error) {
    console.error(`Error fetching bill ${req.params.id}:`, error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to fetch bill' }
    );
  }
});

// Create a new bill
billsRouter.post('/', async (req, res) => {
  try {
    const billData = req.body;
    const response = await accountsClient.post('/bills', billData);
    res.status(201).json(response.data);
  } catch (error) {
    console.error('Error creating bill:', error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to create bill' }
    );
  }
});

// Update a bill
billsRouter.put('/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const billData = req.body;
    const response = await accountsClient.put(`/bills/${id}`, billData);
    res.json(response.data);
  } catch (error) {
    console.error(`Error updating bill ${req.params.id}:`, error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to update bill' }
    );
  }
});

// Delete a bill
billsRouter.delete('/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const response = await accountsClient.delete(`/bills/${id}`);
    res.json(response.data);
  } catch (error) {
    console.error(`Error deleting bill ${req.params.id}:`, error.message);
    res.status(error.response?.status || 500).json(
      error.response?.data || { error: 'Failed to delete bill' }
    );
  }
});

module.exports = { billsRouter };
