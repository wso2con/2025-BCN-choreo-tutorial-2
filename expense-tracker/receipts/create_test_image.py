from PIL import Image, ImageDraw, ImageFont
import os

def create_sample_receipt():
    """
    Create a synthetic receipt image for testing purposes.
    """
    # Create a blank white image
    img_width = 400
    img_height = 800
    image = Image.new('RGB', (img_width, img_height), color='white')
    draw = ImageDraw.Draw(image)
    
    # Try to get a font, fallback to default if not available
    try:
        font = ImageFont.truetype("arial.ttf", 16)
        small_font = ImageFont.truetype("arial.ttf", 12)
    except IOError:
        font = ImageFont.load_default()
        small_font = ImageFont.load_default()
    
    # Draw store name and header
    draw.text((img_width//2 - 80, 20), "SAMPLE GROCERY", fill="black", font=font)
    draw.text((img_width//2 - 100, 50), "123 Main Street, Anytown", fill="black", font=small_font)
    draw.text((img_width//2 - 50, 70), "Tel: 555-1234", fill="black", font=small_font)
    draw.text((20, 100), "Receipt #: 45678", fill="black", font=small_font)
    draw.text((20, 120), "Date: 2023-05-15 14:30", fill="black", font=small_font)
    draw.text((20, 140), "Cashier: John", fill="black", font=small_font)
    
    # Draw horizontal line
    draw.line([(20, 160), (img_width-20, 160)], fill="black", width=1)
    
    # Draw column headers
    draw.text((20, 170), "Item", fill="black", font=small_font)
    draw.text((200, 170), "Qty", fill="black", font=small_font)
    draw.text((250, 170), "Price", fill="black", font=small_font)
    draw.text((320, 170), "Total", fill="black", font=small_font)
    
    # Draw horizontal line
    draw.line([(20, 190), (img_width-20, 190)], fill="black", width=1)
    
    # Sample items
    items = [
        {"name": "Milk", "qty": 1, "price": 3.99},
        {"name": "Bread", "qty": 2, "price": 2.49},
        {"name": "Eggs (dozen)", "qty": 1, "price": 4.99},
        {"name": "Bananas", "qty": 3, "price": 0.59},
        {"name": "Cereal", "qty": 1, "price": 4.29},
        {"name": "Orange Juice", "qty": 1, "price": 3.49},
        {"name": "Cheese", "qty": 1, "price": 5.99},
    ]
    
    # Draw items
    y_pos = 200
    total = 0
    for item in items:
        item_total = item["qty"] * item["price"]
        total += item_total
        
        draw.text((20, y_pos), item["name"], fill="black", font=small_font)
        draw.text((200, y_pos), str(item["qty"]), fill="black", font=small_font)
        draw.text((250, y_pos), f"${item['price']:.2f}", fill="black", font=small_font)
        draw.text((320, y_pos), f"${item_total:.2f}", fill="black", font=small_font)
        
        y_pos += 25
    
    # Draw horizontal line
    draw.line([(20, y_pos), (img_width-20, y_pos)], fill="black", width=1)
    y_pos += 20
    
    # Draw subtotal, tax, and total
    tax = total * 0.08  # 8% tax
    draw.text((220, y_pos), "Subtotal:", fill="black", font=small_font)
    draw.text((320, y_pos), f"${total:.2f}", fill="black", font=small_font)
    y_pos += 25
    
    draw.text((220, y_pos), "Tax (8%):", fill="black", font=small_font)
    draw.text((320, y_pos), f"${tax:.2f}", fill="black", font=small_font)
    y_pos += 25
    
    draw.text((220, y_pos), "Total:", fill="black", font=font)
    draw.text((320, y_pos), f"${(total + tax):.2f}", fill="black", font=font)
    y_pos += 40
    
    # Draw payment method
    draw.text((20, y_pos), "Payment Method: Credit Card", fill="black", font=small_font)
    y_pos += 25
    draw.text((20, y_pos), "Card #: XXXX-XXXX-XXXX-1234", fill="black", font=small_font)
    y_pos += 40
    
    # Draw thank you message
    draw.text((img_width//2 - 80, y_pos), "Thank you for shopping!", fill="black", font=small_font)
    y_pos += 25
    draw.text((img_width//2 - 50, y_pos), "Please come again", fill="black", font=small_font)
    
    # Save the image
    save_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "sample_receipt.jpg")
    image.save(save_path, "JPEG", quality=95)
    print(f"Sample receipt created at: {save_path}")
    return save_path

if __name__ == "__main__":
    create_sample_receipt()