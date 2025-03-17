import React, { useState, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  List, 
  ListItem, 
  ListItemText, 
  Paper, 
  TextField, 
  Button, 
  Grid,
  AppBar,
  Toolbar,
  Fab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tabs,
  Tab,
  IconButton,
  useMediaQuery,
  CssBaseline,
  ThemeProvider,
  createTheme,
  Snackbar,
  Alert,
  CircularProgress
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PhotoCameraIcon from '@mui/icons-material/PhotoCamera';
import ImageIcon from '@mui/icons-material/Image';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import './App.css';
import { billsApi, billParserApi, Bill, BillInput } from './services/api';

// Define green money theme
const theme = createTheme({
  palette: {
    primary: {
      main: '#2e7d32', // Money green
      light: '#60ad5e',
      dark: '#005005',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#66bb6a',
      light: '#98ee99',
      dark: '#338a3e',
      contrastText: '#000000',
    },
    background: {
      default: '#f5f5f5',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h6: {
      fontWeight: 600,
    },
  },
});

// Interface for transformed bill data for display
interface ExpenseBill {
  id: number;
  title: string;
  amount: number;
  date: string;
  paid: boolean;
}

function App() {
  const [bills, setBills] = useState<Bill[]>([]);
  const [expenses, setExpenses] = useState<ExpenseBill[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [amount, setAmount] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [currentTab, setCurrentTab] = useState(0);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');
  const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>('success');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const isMobile = useMediaQuery('(max-width:600px)');

  // Fetch bills from API
  const fetchBills = async () => {
    setIsLoading(true);
    try {
      const data = await billsApi.getAllBills();
      setBills(data);
      
      // Transform bills to expenses for display
      const transformedExpenses: ExpenseBill[] = data.map(bill => ({
        id: bill.id,
        title: bill.title,
        amount: bill.total,
        date: new Date(bill.due_date).toISOString().split('T')[0],
        paid: bill.paid
      }));
      
      setExpenses(transformedExpenses);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch bills:', err);
      setError('Failed to load expenses. Please try again later.');
      showSnackbar('Failed to load expenses', 'error');
    } finally {
      setIsLoading(false);
    }
  };

  // Load bills on component mount
  useEffect(() => {
    fetchBills();
  }, []);

  const handleOpenDialog = () => {
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    setCurrentTab(0);
    setTitle('');
    setAmount('');
    setSelectedFile(null);
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  // Show snackbar message
  const showSnackbar = (message: string, severity: 'success' | 'error') => {
    setSnackbarMessage(message);
    setSnackbarSeverity(severity);
    setSnackbarOpen(true);
  };

  // Handle manual expense addition
  const handleAddExpense = async () => {
    if (title && amount) {
      setIsSubmitting(true);
      
      try {
        const billData: BillInput = {
          title,
          description: 'Manually added expense',
          due_date: new Date().toISOString().split('T')[0],
          paid: false,
          items: [
            {
              name: title,
              description: '',
              amount: parseFloat(amount),
              quantity: 1
            }
          ]
        };
        
        await billsApi.createBill(billData);
        showSnackbar('Expense added successfully', 'success');
        fetchBills();
        handleCloseDialog();
      } catch (err) {
        console.error('Failed to add expense:', err);
        showSnackbar('Failed to add expense', 'error');
      } finally {
        setIsSubmitting(false);
      }
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  // Handle bill upload and processing
  const handleFileUpload = async () => {
    if (selectedFile) {
      setIsSubmitting(true);
      
      try {
        // Upload the image and create a bill
        const result = await billParserApi.createBillFromImage(
          selectedFile,
          `Bill from ${selectedFile.name}`
        );
        
        showSnackbar('Bill processed successfully', 'success');
        fetchBills();
        handleCloseDialog();
      } catch (err) {
        console.error('Failed to process bill image:', err);
        showSnackbar('Failed to process bill image', 'error');
      } finally {
        setIsSubmitting(false);
      }
    }
  };

  const calculateTotal = () => {
    return expenses.reduce((sum, expense) => sum + expense.amount, 0).toFixed(2);
  };

  const handleCloseSnackbar = () => {
    setSnackbarOpen(false);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div className="App" style={{
        backgroundImage: `url("data:image/svg+xml,%3Csvg width='80' height='80' viewBox='0 0 80 80' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='%232e7d32' fill-opacity='0.05' fill-rule='evenodd'%3E%3Cpath d='M50 50c0-5.523 4.477-10 10-10s10 4.477 10 10-4.477 10-10 10c0 5.523-4.477 10-10 10s-10-4.477-10-10 4.477-10 10-10zM10 10c0-5.523 4.477-10 10-10s10 4.477 10 10-4.477 10-10 10c0 5.523-4.477 10-10 10S0 25.523 0 20s4.477-10 10-10zm10 8c4.418 0 8-3.582 8-8s-3.582-8-8-8-8 3.582-8 8 3.582 8 8 8zm40 40c4.418 0 8-3.582 8-8s-3.582-8-8-8-8 3.582-8 8 3.582 8 8 8z' /%3E%3Cpath d='M30 10c0-2.21-1.79-4-4-4s-4 1.79-4 4 1.79 4 4 4 4-1.79 4-4zm-4 2c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm-4 14c0-2.21-1.79-4-4-4s-4 1.79-4 4 1.79 4 4 4 4-1.79 4-4zm-4 2c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm12-8c0-2.21-1.79-4-4-4s-4 1.79-4 4 1.79 4 4 4 4-1.79 4-4zm-4 2c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z' fill='%23338a3e' fill-opacity='0.03' /%3E%3C/g%3E%3C/svg%3E")`,
        minHeight: '100vh',
        paddingBottom: '80px',
        position: 'relative',
        overflow: 'hidden'
      }}>
        {/* Animated dollar signs in background */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          overflow: 'hidden',
          pointerEvents: 'none',
          zIndex: 0
        }}>
          {[...Array(15)].map((_, i) => (
            <div
              key={i}
              style={{
                position: 'absolute',
                left: `${Math.random() * 100}%`,
                top: `${Math.random() * 100}%`,
                color: '#2e7d32',
                opacity: 0.05,
                fontSize: `${Math.random() * 20 + 10}px`,
                animation: `float ${Math.random() * 5 + 5}s infinite`,
                animationDelay: `${Math.random() * 5}s`
              }}
            >
              $
            </div>
          ))}
        </div>
        <AppBar position="static" color="primary">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Monthly Expense Tracker
            </Typography>
          </Toolbar>
        </AppBar>
        
        <Container maxWidth="md" sx={{ mt: 3 }}>
          <Paper 
            elevation={3} 
            sx={{ 
              p: 2, 
              mb: 2, 
              borderRadius: '12px',
              backgroundColor: 'rgba(255, 255, 255, 0.95)',
            }}
          >
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" color="primary" gutterBottom sx={{ mb: 0 }}>
                Expense Summary
              </Typography>
              <Typography variant="h6" color="primary" sx={{ fontWeight: 'bold' }}>
                ${calculateTotal()}
              </Typography>
            </Box>
            
            {isLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                <CircularProgress color="primary" />
              </Box>
            ) : error ? (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <Typography variant="body1" color="error">
                  {error}
                </Typography>
                <Button 
                  variant="outlined" 
                  color="primary" 
                  sx={{ mt: 2 }}
                  onClick={() => fetchBills()}
                >
                  Retry
                </Button>
              </Box>
            ) : (
              <List sx={{ 
                maxHeight: isMobile ? 'calc(100vh - 250px)' : 'none',
                overflow: isMobile ? 'auto' : 'visible'
              }}>
                {expenses.length === 0 ? (
                  <Box sx={{ textAlign: 'center', py: 4 }}>
                    <Typography variant="body1" color="text.secondary">
                      No expenses yet. Add your first expense!
                    </Typography>
                  </Box>
                ) : (
                  expenses.map((expense) => (
                    <ListItem 
                      key={expense.id} 
                      divider 
                      sx={{ 
                        borderRadius: '8px', 
                        mb: 1,
                        '&:hover': { 
                          backgroundColor: 'rgba(102, 187, 106, 0.1)' 
                        }
                      }}
                    >
                      <ListItemText 
                        primary={expense.title} 
                        secondary={`Date: ${expense.date} ${expense.paid ? '• Paid' : '• Unpaid'}`} 
                        primaryTypographyProps={{ fontWeight: 500 }}
                      />
                      <Typography variant="body1" sx={{ 
                        fontWeight: 'bold', 
                        color: theme.palette.primary.main
                      }}>
                        ${expense.amount.toFixed(2)}
                      </Typography>
                    </ListItem>
                  ))
                )}
              </List>
            )}
          </Paper>
        </Container>

        {/* Floating Action Button */}
        <Fab 
          color="primary" 
          aria-label="add" 
          sx={{ 
            position: 'fixed', 
            bottom: 20, 
            right: 20,
            boxShadow: '0 8px 16px rgba(46, 125, 50, 0.3)'
          }}
          onClick={handleOpenDialog}
        >
          <AttachMoneyIcon />
        </Fab>

        {/* Add Expense Dialog */}
        <Dialog 
          open={openDialog} 
          onClose={handleCloseDialog} 
          fullScreen={isMobile}
          fullWidth
          maxWidth="sm"
          PaperProps={{
            sx: {
              borderRadius: isMobile ? 0 : '12px',
              overflow: 'hidden'
            }
          }}
        >
          <DialogTitle sx={{ 
            bgcolor: 'primary.main', 
            color: 'white',
            m: 0,
            p: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between'
          }}>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              {isMobile && (
                <IconButton
                  edge="start"
                  color="inherit"
                  onClick={handleCloseDialog}
                  aria-label="close"
                  sx={{ mr: 1 }}
                >
                  <ArrowBackIcon />
                </IconButton>
              )}
              <Typography variant="h6">
                {currentTab === 0 ? "Add New Expense" : "Upload Bill"}
              </Typography>
            </Box>
            {!isMobile && (
              <IconButton
                edge="end"
                color="inherit"
                onClick={handleCloseDialog}
                aria-label="close"
              >
                <CloseIcon />
              </IconButton>
            )}
          </DialogTitle>
          
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs 
              value={currentTab} 
              onChange={handleTabChange} 
              variant="fullWidth"
              textColor="primary"
              indicatorColor="primary"
            >
              <Tab label="Manual Entry" />
              <Tab label="Upload Bill" />
            </Tabs>
          </Box>
          
          <DialogContent>
            {currentTab === 0 ? (
              <Box component="form" noValidate sx={{ mt: 1 }}>
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  id="title"
                  label="Title"
                  name="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  autoFocus
                />
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  name="amount"
                  label="Amount"
                  type="number"
                  id="amount"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
              </Box>
            ) : (
              <Box sx={{ mt: 2, textAlign: 'center' }}>
                <Grid container spacing={2} justifyContent="center">
                  <Grid item>
                    <input
                      accept="image/*"
                      style={{ display: 'none' }}
                      id="gallery-file"
                      type="file"
                      onChange={handleFileChange}
                    />
                    <label htmlFor="gallery-file">
                      <Button 
                        variant="outlined" 
                        component="span" 
                        startIcon={<ImageIcon />}
                        sx={{ mb: 2 }}
                      >
                        From Gallery
                      </Button>
                    </label>
                  </Grid>
                  
                  <Grid item>
                    <input
                      accept="image/*"
                      style={{ display: 'none' }}
                      id="camera-file"
                      type="file"
                      capture="environment"
                      onChange={handleFileChange}
                    />
                    <label htmlFor="camera-file">
                      <Button 
                        variant="contained" 
                        component="span"
                        color="primary"
                        startIcon={<PhotoCameraIcon />}
                        sx={{ mb: 2 }}
                      >
                        Take Photo
                      </Button>
                    </label>
                  </Grid>
                </Grid>
                
                {selectedFile && (
                  <Box sx={{ mt: 3 }}>
                    <Paper 
                      elevation={1} 
                      sx={{ 
                        p: 2, 
                        mb: 2, 
                        maxWidth: '100%',
                        borderRadius: '8px',
                        backgroundColor: 'rgba(102, 187, 106, 0.1)',
                        border: '1px dashed #2e7d32'
                      }}
                    >
                      <Typography variant="body2" sx={{ mb: 1 }}>
                        Selected: {selectedFile.name}
                      </Typography>
                      {selectedFile && (
                        <Box sx={{ 
                          width: '100%', 
                          display: 'flex', 
                          justifyContent: 'center',
                          mt: 2 
                        }}>
                          <img 
                            src={URL.createObjectURL(selectedFile)} 
                            alt="Bill preview" 
                            style={{ 
                              maxWidth: '100%',
                              maxHeight: '200px',
                              borderRadius: '4px',
                              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                            }} 
                          />
                        </Box>
                      )}
                    </Paper>
                  </Box>
                )}
              </Box>
            )}
          </DialogContent>
          
          <DialogActions sx={{ p: 3, pt: 0 }}>
            <Button 
              onClick={handleCloseDialog} 
              color="primary"
              variant="outlined"
              sx={{ borderRadius: '8px' }}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button 
              onClick={currentTab === 0 ? handleAddExpense : handleFileUpload}
              color="primary"
              variant="contained"
              disabled={currentTab === 0 ? !title || !amount || isSubmitting : !selectedFile || isSubmitting}
              sx={{ borderRadius: '8px' }}
              startIcon={isSubmitting ? <CircularProgress size={20} color="inherit" /> : null}
            >
              {isSubmitting ? 'Processing...' : currentTab === 0 ? 'Add Expense' : 'Process Bill'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Snackbar for notifications */}
        <Snackbar 
          open={snackbarOpen} 
          autoHideDuration={6000} 
          onClose={handleCloseSnackbar}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
          <Alert 
            onClose={handleCloseSnackbar} 
            severity={snackbarSeverity} 
            variant="filled"
            sx={{ width: '100%' }}
          >
            {snackbarMessage}
          </Alert>
        </Snackbar>
      </div>
    </ThemeProvider>
  );
}

export default App;
