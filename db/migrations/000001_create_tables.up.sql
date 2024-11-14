-- Create the products table (Product Catalog) as a partitioned table with primary key (id, category)
CREATE TABLE products (
    id SERIAL,                    -- Unique identifier for each product (NOT the primary key)
    name VARCHAR(255) NOT NULL,    -- Product name
    description TEXT,             -- Detailed description of the product
    category VARCHAR(100),        -- Product category (e.g., Electronics, Apparel)
    price DECIMAL(10, 2) NOT NULL,-- Price of the product (e.g., 99.99)
    stock_quantity INT DEFAULT 0, -- Available stock
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the product was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the product details were last updated
    PRIMARY KEY (id, category)    -- Modified to include 'category' in the primary key
) PARTITION BY LIST (category);    -- Partition by the 'category' column

-- Create partition for Electronics products
CREATE TABLE products_electronics PARTITION OF products
    FOR VALUES IN ('Electronics');

-- Create partition for Apparel products
CREATE TABLE products_apparel PARTITION OF products
    FOR VALUES IN ('Apparel');

-- Create the orders table (for managing customer orders)
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,               -- Unique identifier for each order
    customer_id INT NOT NULL,            -- Reference to the customer placing the order
    status VARCHAR(50) DEFAULT 'pending',-- Order status (e.g., pending, completed, canceled)
    total_amount DECIMAL(10, 2),         -- Total amount for the order
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Order creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Order status update timestamp
);

-- Create the order_items table (for storing the items in each order)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,               -- Unique identifier for the order item
    order_id INT REFERENCES orders(id),  -- Reference to the order
    product_id INT,                      -- Reference to the product (will be foreign key)
    category VARCHAR(100),               -- Ensure category is referenced for sharding
    quantity INT NOT NULL,               -- Quantity of the product in the order
    price DECIMAL(10, 2) NOT NULL,        -- Price at the time of ordering (in case of price updates)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Item creation timestamp
    FOREIGN KEY (product_id, category) REFERENCES products(id, category) -- Foreign key referencing the partitioned products table
);

-- Create the customers table (optional for storing customer info)
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,      -- Unique identifier for each customer
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,  -- Ensure the email is unique
    phone VARCHAR(20),          -- Optional phone number
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Customer creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Customer details update timestamp
);

-- Indexing for better search performance (for Elasticsearch integration)
CREATE INDEX idx_product_name ON products (name);
CREATE INDEX idx_product_category ON products (category);

-- Optional: Create triggers or functions for automatic updates (for example, updating timestamps)
-- Trigger to update the "updated_at" timestamp whenever a product is updated
CREATE OR REPLACE FUNCTION update_timestamp() 
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();  -- Update the "updated_at" column
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to "products" table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Insert 20 random products
DO $$ 
DECLARE 
    product_names TEXT[] := ARRAY['Wireless Mouse', 'Bluetooth Speaker', 'Laptop', 'Smartphone', 'Keyboard', 'Tablet', 'Headphones', 'Monitor', 'Smartwatch', 'External Hard Drive', 'Game Console', 'TV', 'Camera', 'Projector', 'Drone', 'Action Camera', 'Smart Light', 'VR Headset', 'Gaming Chair', 'Power Bank'];
    product_categories TEXT[] := ARRAY['Electronics', 'Apparel', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Electronics', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel', 'Apparel'];
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
