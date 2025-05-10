import axios from 'axios';

/**
 * API service for interacting with the BFF API
 */

// Create axios instance with base URL from environment variables
const api = axios.create({
  // @ts-ignore
  baseURL: window.configs.apiUrl
});

// Add response interceptor to handle 401 errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 || error.response?.status === 403) {
      window.location.href = '/auth/login';
    }
    return Promise.reject(error);
  }
);

// Define interfaces matching API responses
export interface BillItem {
  id?: number;
  bill_id?: number;
  name: string;
  description?: string;
  amount: number;
  quantity: number;
  created_at?: string;
  updated_at?: string;
}

export interface Bill {
  id: number;
  title: string;
  description?: string;
  total: number;
  due_date: string;
  paid: boolean;
  items?: BillItem[];
  item_count?: number;
  created_at: string;
  updated_at: string;
}

export interface BillInput {
  title: string;
  description?: string;
  due_date?: string;
  paid?: boolean;
  items?: {
    name: string;
    description?: string;
    amount: number;
    quantity?: number;
  }[];
}

export interface ParsedReceipt {
  items: {
    name: string;
    quantity: number;
    price: number;
  }[];
  total: number;
  currency?: string;
  date?: string;
  merchant?: string;
}

// Bills API methods
export const billsApi = {
  // Get all bills
  getAllBills: async (): Promise<Bill[]> => {
    const response = await api.get('/bills');
    return response.data;
  },

  // Get a bill by ID
  getBillById: async (id: number): Promise<Bill> => {
    const response = await api.get(`/bills/${id}`);
    return response.data;
  },

  // Create a new bill
  createBill: async (billData: BillInput): Promise<{ id: number }> => {
    const response = await api.post('/bills', billData);
    return response.data;
  },

  // Update a bill
  updateBill: async (id: number, billData: BillInput): Promise<{ message: string }> => {
    const response = await api.put(`/bills/${id}`, billData);
    return response.data;
  },

  // Delete a bill
  deleteBill: async (id: number): Promise<{ message: string }> => {
    const response = await api.delete(`/bills/${id}`);
    return response.data;
  }
};

// Bill Parser API methods
export const billParserApi = {
  // Parse a bill image
  parseBill: async (imageFile: File): Promise<ParsedReceipt> => {
    const formData = new FormData();
    formData.append('image', imageFile);
    formData.append('create_bill', 'false');

    const response = await api.post('/parser/parse', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  },

  // Parse an image and create a bill
  createBillFromImage: async (imageFile: File, title?: string): Promise<{
    message: string;
    billId: number;
    parsedData: ParsedReceipt;
  }> => {
    const formData = new FormData();
    formData.append('image', imageFile);
    formData.append('create_bill', 'true');
    if (title) {
      formData.append('title', title);
    }

    const response = await api.post('/parser/parse', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  }
};

export default { billsApi, billParserApi };
