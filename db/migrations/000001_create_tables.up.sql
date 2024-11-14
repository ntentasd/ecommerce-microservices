-- Create partitioned products table with primary key (id, category)
CREATE TABLE products (
    id SERIAL,                    
    name VARCHAR(255) NOT NULL,    
    description TEXT,             
    category VARCHAR(100),        
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INT DEFAULT 0, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    PRIMARY KEY (id, category)    
) PARTITION BY LIST (category);    

-- Create partitions for Electronics and Apparel
CREATE TABLE products_electronics PARTITION OF products
    FOR VALUES IN ('Electronics');

CREATE TABLE products_apparel PARTITION OF products
    FOR VALUES IN ('Apparel');

-- Create orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,               
    customer_id INT NOT NULL,            
    status VARCHAR(50) DEFAULT 'pending',
    total_amount DECIMAL(10, 2),         
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   
);

-- Create order_items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,               
    order_id INT REFERENCES orders(id),  
    product_id INT,                      
    category VARCHAR(100),               
    quantity INT NOT NULL,               
    price DECIMAL(10, 2) NOT NULL,        
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (product_id, category) REFERENCES products(id, category) 
);

-- Create customers table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,      
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,  
    phone VARCHAR(20),          
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  
);

-- Add indexes for improved search performance
CREATE INDEX idx_product_name ON products (name);
CREATE INDEX idx_product_category ON products (category);

-- Function to automatically update timestamps
CREATE OR REPLACE FUNCTION update_timestamp() 
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();  
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update product "updated_at" timestamp
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Insert sample product data
DO $$ 
DECLARE 
    product_names TEXT[] := ARRAY['Wireless Mouse', 'Bluetooth Speaker', 'Laptop', 'Smartphone', 'Keyboard', 'Tablet', 'Headphones', 'Monitor', 'Smartwatch', 'External Hard Drive', 'Game Console', 'TV', 'Camera', 'Projector', 'Drone', 'Action Camera', 'Smart Light', 'VR Headset', 'Gaming Chair', 'Power Bank'];
    product_categories TEXT[] := ARRAY['Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel'];
    product_prices DECIMAL[] := ARRAY[29.99, 49.99, 799.99, 999.99, 59.99, 149.99, 199.99, 249.99, 299.99, 99.99, 199.99, 399.99, 499.99, 79.99, 149.99, 89.99, 69.99, 399.99, 129.99, 59.99];
    product_stock INT[] := ARRAY[100, 50, 20, 15, 200, 75, 80, 40, 60, 150, 30, 25, 18, 90, 85, 100, 95, 110, 70, 65];
    i INT := 1;
BEGIN
    FOR i IN 1..20 LOOP
        INSERT INTO products (name, description, category, price, stock_quantity)
        VALUES (
            product_names[i], 
            CONCAT('Description of ', product_names[i]), 
            product_categories[i], 
            product_prices[i], 
            product_stock[i]
        );
    END LOOP;
END $$;
