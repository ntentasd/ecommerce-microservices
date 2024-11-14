-- Drop the customers table
DROP TABLE IF EXISTS customers CASCADE;

-- Drop the foreign key constraint before dropping the orders table
ALTER TABLE order_items DROP CONSTRAINT order_items_order_id_fkey;

-- Drop the products table and partitions
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS products_electronics;
DROP TABLE IF EXISTS products_apparel;

-- Drop the orders table
DROP TABLE IF EXISTS orders;

-- Drop the order_items table
DROP TABLE IF EXISTS order_items;