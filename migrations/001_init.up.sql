-- User Permission Tables

CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=inactive',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_username (username),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Users table';

CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    description VARCHAR(255),
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=inactive',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_name (name),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Roles table';

CREATE TABLE permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(64) NOT NULL COMMENT 'Permission code, e.g. product:list',
    name VARCHAR(128) NOT NULL COMMENT 'Permission display name',
    resource VARCHAR(64) NOT NULL COMMENT 'Resource name',
    action VARCHAR(16) NOT NULL COMMENT 'Action: list, create, update, delete, export',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_code (code),
    KEY idx_resource (resource),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Permissions table';

CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_user_role (user_id, role_id),
    KEY idx_role_id (role_id),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User-Role association table';

CREATE TABLE role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    KEY idx_permission_id (permission_id),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Role-Permission association table';

-- Core Business Tables

CREATE TABLE warehouses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    address VARCHAR(255),
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=inactive',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Warehouses table';

CREATE TABLE locations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT NOT NULL,
    zone VARCHAR(32) COMMENT 'Zone identifier',
    shelf VARCHAR(32) COMMENT 'Shelf identifier',
    level VARCHAR(32) COMMENT 'Level identifier',
    position VARCHAR(32) COMMENT 'Position identifier',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1=active, 0=inactive',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_warehouse_id (warehouse_id),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at),
    KEY idx_location_code (warehouse_id, zone, shelf, level, position)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Locations table';

CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    parent_id BIGINT NULL COMMENT 'Parent category ID',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_parent_id (parent_id),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Categories table';

CREATE TABLE products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sku VARCHAR(64) NOT NULL COMMENT 'Stock Keeping Unit',
    name VARCHAR(255) NOT NULL,
    category_id BIGINT,
    spec VARCHAR(255) COMMENT 'Product specification',
    unit VARCHAR(32) COMMENT 'Unit of measurement',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_sku (sku),
    KEY idx_category_id (category_id),
    KEY idx_name (name),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Products table';

CREATE TABLE inventory (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    location_id BIGINT,
    quantity INT NOT NULL DEFAULT 0,
    batch_no VARCHAR(64) COMMENT 'Batch number',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_warehouse_product (warehouse_id, product_id),
    KEY idx_product_id (product_id),
    KEY idx_location_id (location_id),
    KEY idx_batch_no (batch_no),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Inventory table';

CREATE TABLE suppliers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    contact VARCHAR(64) COMMENT 'Contact person',
    phone VARCHAR(32),
    address VARCHAR(255),
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_name (name),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Suppliers table';

CREATE TABLE customers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    contact VARCHAR(64) COMMENT 'Contact person',
    phone VARCHAR(32),
    address VARCHAR(255),
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_name (name),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Customers table';

-- Order Tables

CREATE TABLE inbound_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(64) NOT NULL COMMENT 'Inbound order number',
    supplier_id BIGINT,
    warehouse_id BIGINT NOT NULL,
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending, 1=confirmed, 2=cancelled',
    total_quantity INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_supplier_id (supplier_id),
    KEY idx_warehouse_id (warehouse_id),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Inbound orders table';

CREATE TABLE inbound_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    inbound_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    location_id BIGINT,
    quantity INT NOT NULL DEFAULT 0,
    batch_no VARCHAR(64) COMMENT 'Batch number',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_inbound_order_id (inbound_order_id),
    KEY idx_product_id (product_id),
    KEY idx_location_id (location_id),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Inbound items table';

CREATE TABLE outbound_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(64) NOT NULL COMMENT 'Outbound order number',
    customer_id BIGINT,
    warehouse_id BIGINT NOT NULL,
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending, 1=confirmed, 2=cancelled',
    total_quantity INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_customer_id (customer_id),
    KEY idx_warehouse_id (warehouse_id),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Outbound orders table';

CREATE TABLE outbound_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    outbound_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    location_id BIGINT,
    quantity INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    KEY idx_outbound_order_id (outbound_order_id),
    KEY idx_product_id (product_id),
    KEY idx_location_id (location_id),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Outbound items table';

CREATE TABLE stock_transfers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(64) NOT NULL COMMENT 'Transfer order number',
    source_warehouse_id BIGINT NOT NULL,
    target_warehouse_id BIGINT NOT NULL,
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending, 1=in_transit, 2=completed, 3=cancelled',
    created_at DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_source_warehouse_id (source_warehouse_id),
    KEY idx_target_warehouse_id (target_warehouse_id),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Stock transfers table';

-- Audit Table

CREATE TABLE audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    table_name VARCHAR(64) NOT NULL COMMENT 'Table name being audited',
    record_id BIGINT NOT NULL COMMENT 'Record ID',
    action VARCHAR(16) NOT NULL COMMENT 'Action type: create, update, delete',
    old_value JSON COMMENT 'Value before change',
    new_value JSON COMMENT 'Value after change',
    operated_by BIGINT NOT NULL COMMENT 'User who performed the operation',
    operated_at DATETIME NOT NULL COMMENT 'Operation timestamp',
    ip_address VARCHAR(45) COMMENT 'IP address of the operator',
    KEY idx_table_record (table_name, record_id),
    KEY idx_operated_by (operated_by),
    KEY idx_operated_at (operated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Audit logs table';
