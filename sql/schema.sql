CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY users_uc_email (email)
) ENGINE=InnoDB;

CREATE TABLE snippets (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires DATETIME NOT NULL,
    PRIMARY KEY (id),
    KEY idx_snippets_created (created)
) ENGINE=InnoDB;

CREATE TABLE sessions (
    token CHAR(43) NOT NULL,
    data BLOB NOT NULL,
    expiry DATETIME NOT NULL,
    PRIMARY KEY (token),
    KEY sessions_expiry_idx (expiry)
) ENGINE=InnoDB;
