CREATE TABLE files {
    id SERIAL,
    folder_id INT,
    owner_id INT NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,
    path VARCHAR(250) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    deleted BOOL DEFAULT FALSE,
    PRIMARY KEY (id),
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folders(id),
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id)
}