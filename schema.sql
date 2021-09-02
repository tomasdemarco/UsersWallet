CREATE TABLE users(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
	name VARCHAR(255) NOT NULL,
	lastname VARCHAR(255) NOT NULL,
	documento INTEGER NOT NULL,
	birthday VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
	phone INTEGER NOT NULL,
    password VARCHAR(255) NOT NULL
)