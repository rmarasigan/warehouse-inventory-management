package db

const (
	users string = `CREATE TABLE IF NOT EXISTS users (
						id INT NOT NULL AUTO_INCREMENT,
						role_id INT NOT NULL,
						first_name VARCHAR(30) NOT NULL,
						last_name VARCHAR(30) NOT NULL,
						email VARCHAR(50) UNIQUE,
						password VARCHAR(255) NOT NULL,
						last_login TIMESTAMP NULL,
						date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
						is_active BOOLEAN DEFAULT TRUE,
						PRIMARY KEY (id),
						INDEX idx_last_name (last_name),
						CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES role(id)
					);`

	role string = `CREATE TABLE IF NOT EXISTS role (
						id INT NOT NULL AUTO_INCREMENT,
						name VARCHAR(30) NOT NULL,
						PRIMARY KEY (id),
						UNIQUE KEY idx_name (name)
					);`

	uom string = `CREATE TABLE IF NOT EXISTS unit_of_measurement (
						id INT NOT NULL AUTO_INCREMENT,
						code VARCHAR(10) UNIQUE,
						name VARCHAR(30) NOT NULL,
						PRIMARY KEY (id)
					);`

	storage string = `CREATE TABLE IF NOT EXISTS storage (
								id INT NOT NULL AUTO_INCREMENT,
								code VARCHAR(10) NOT NULL UNIQUE,
								name VARCHAR(50) NOT NULL,
								description VARCHAR(50),
								PRIMARY KEY (id),
								INDEX idx_name (name)
							);`

	currency string = `CREATE TABLE IF NOT EXISTS currency (
								id INT NOT NULL AUTO_INCREMENT,
								code VARCHAR(10) NOT NULL UNIQUE,
								symbol VARCHAR(10) NOT NULL,
								is_active BOOLEAN DEFAULT TRUE,
								PRIMARY KEY (id)
							);`

	item string = `CREATE TABLE IF NOT EXISTS item (
						id INT NOT NULL AUTO_INCREMENT,
						name VARCHAR(70) NOT NULL,
						description VARCHAR(100),
						quantity INT NOT NULL,
						unit_price DECIMAL(10,2) NOT NULL,
						uom_id INT NOT NULL,
						stock_status VARCHAR(20) NOT NULL,
						storage_id INT NOT NULL,
						created_by INT NOT NULL,
						date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
						PRIMARY KEY (id),
						INDEX idx_uom_id (uom_id),
						INDEX idx_stock_status (stock_status),
						CONSTRAINT fk_item_uom FOREIGN KEY (uom_id) REFERENCES unit_of_measurement(id),
						CONSTRAINT fk_item_storage FOREIGN KEY (storage_id) REFERENCES storage(id),
						CONSTRAINT fk_item_creator FOREIGN KEY (created_by) REFERENCES users(id)
					);`

	transactions string = `CREATE TABLE IF NOT EXISTS transactions (
										id INT NOT NULL AUTO_INCREMENT,
										reference VARCHAR(70) NOT NULL,
										amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
										type VARCHAR(20) NOT NULL,
										is_cancelled BOOLEAN DEFAULT FALSE,
										note VARCHAR(255),
										created_by INT NOT NULL,
										updated_by INT,
										date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
										date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
										PRIMARY KEY (id),
										UNIQUE KEY idx_reference (reference),
										INDEX idx_created_by (created_by),
										INDEX id_updated_by (updated_by),
										CONSTRAINT fk_transaction_creator FOREIGN KEY (created_by) REFERENCES users(id),
										CONSTRAINT fk_transaction_updater FOREIGN KEY (updated_by) REFERENCES users(id)
									);`

	orderline string = `CREATE TABLE IF NOT EXISTS orderline (
									id INT NOT NULL AUTO_INCREMENT,
									transaction_id INT NOT NULL,
									item_id INT NOT NULL,
									quantity INT NOT NULL,
									unit_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
									total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
									note VARCHAR(255),
									is_voided BOOLEAN DEFAULT FALSE,
									created_by INT NOT NULL,
									updated_by INT,
									date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
									date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
									PRIMARY KEY (id),
									INDEX idx_item_id (item_id),
									INDEX idx_created_by (created_by),
									INDEX id_updated_by (updated_by),
									CONSTRAINT fk_orderline_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id),
									CONSTRAINT fk_orderline_item FOREIGN KEY (item_id) REFERENCES item(id),
									CONSTRAINT fk_orderline_creator FOREIGN KEY (created_by) REFERENCES users(id),
									CONSTRAINT fk_orderline_updater FOREIGN KEY (updated_by) REFERENCES users(id)
								);`

	roleInsert string = `INSERT INTO role (name)
								SELECT :name FROM DUAL
								WHERE NOT EXISTS (SELECT 1 FROM role WHERE name = :name);`

	uomInsert string = `INSERT INTO unit_of_measurement (code, name)
							SELECT :code, :name FROM DUAL
							WHERE NOT EXISTS (SELECT 1 FROM unit_of_measurement WHERE code = :code AND name = :name);`

	currencyInsert string = `INSERT INTO currency (code, symbol, is_active)
										SELECT :code, :symbol, :is_active FROM DUAL
										WHERE NOT EXISTS (SELECT 1 FROM currency WHERE code = :code AND symbol = :symbol);`
)

// tablesOrder defines the order to create tables (parents first).
var tablesOrder = []string{
	"role",
	"users",
	"unit_of_measurement",
	"storage",
	"currency",
	"item",
	"transactions",
	"orderline",
}

// databaseTables contains the CREATE TABLE queries.
var databaseTables = map[string]string{
	"role":                role,
	"users":               users,
	"unit_of_measurement": uom,
	"storage":             storage,
	"currency":            currency,
	"item":                item,
	"orderline":           orderline,
	"transactions":        transactions,
}
