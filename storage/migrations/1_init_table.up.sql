CREATE TABLE IF NOT EXISTS cryptocurrency (
    code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    category VARCHAR(20) NOT NULL,
    algorithm VARCHAR(20) NOT NULL,
    platform VARCHAR(20) NOT NULL,
    industry VARCHAR(20) NOT NULL,
    types VARCHAR(10) NOT NULL,
    mineable boolean NOT NULL,
    audited boolean NOT NULL,
    price float NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_c1 (name,category,algorithm,platform,industry,types,mineable,audited,price),
    INDEX idx_c2 (category,algorithm,platform,industry,types,mineable,audited,price),
    INDEX idx_c3 (algorithm,platform,industry,types,mineable,audited,price),
    INDEX idx_c4 (platform,industry,types,mineable,audited,price),
    INDEX idx_c5 (industry,types,mineable,audited,price),
    INDEX idx_c6 (types,mineable,audited,price),
    INDEX idx_c7 (mineable,audited,price),
    INDEX idx_c8 (audited,price),
    INDEX idx_c9 (price)
)  ENGINE=INNODB;

