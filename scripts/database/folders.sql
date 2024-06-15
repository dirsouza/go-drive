CREATE TABLE folders {
    id SERIAL,
    parent_id INT,
    name VARCHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    deleted BOOL DEFAULT FALSE,
    PRIMARY KEY (id),
    CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES folders(id)
}