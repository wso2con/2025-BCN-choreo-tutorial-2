const express = require('express');
const cors = require('cors');
const dotenv = require('dotenv');
const { billsRouter } = require('./src/routes/bills');
const { parserRouter } = require('./src/routes/billParser');

// Load environment variables
dotenv.config();

const app = express();
const PORT = 9090;

// Middleware
app.use(cors());
app.use(express.json());

// Routes
app.use('/api/bills', billsRouter);
app.use('/api/parser', parserRouter);

// Base route for health check
app.get('/api/health', (req, res) => {
  res.json({ status: 'UP', message: 'BFF API is running' });
});

// Start server
app.listen(PORT, () => {
  console.log(`BFF API server listening on port ${PORT}`);
});
