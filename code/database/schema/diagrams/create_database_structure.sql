    CREATE TABLE [IF NOT EXISTS] domain (
        domain_name TEXT NOT NULL,
        domain_extension TEXT NOT NULL,
        
        PRIMARY KEY (domain_name, domain_extension),
    );

CREATE TABLE [IF NOT EXISTS] domain_metadata (
    domain_name TEXT NOT NULL,
    domain_extension TEXT NOT NULL,

    ipv4 VARCHAR(16),
    geolocation TEXT,
    language TEXT, 
    accessible BOOLEAN,

    PRIMARY KEY (domain_name, domain_extension),

    FOREIGN KEY (domain_name, domain_extension) REFERENCES domain (name, extension),
);

CREATE TABLE [IF NOT EXISTS] context_sample (
    domain_name TEXT NOT NULL,
    domain_extension TEXT NOT NULL,
    
    PRIMARY KEY (domain_name, domain_extension),

    FOREIGN KEY (domain_name, domain_extension) REFERENCES domain (name, extension),
);

CREATE TABLE [IF NOT EXISTS] incoming_link (
    protocol TEXT NOT NULL,
    subdomain TEXT NOT NULL,
    domain_name TEXT NOT NULL,
    domain_extension TEXT NOT NULL,
    path TEXT NOT NULL, 


    PRIMARY KEY (protocol, subdomain, domain_name, domain_extension, path),

    FOREIGN KEY (domain_name, domain_extension) REFERENCES domain (name, extension),
    FOREIGN KEY (protocol, subdomain, domain_name, domain_extension, path) REFERENCES url (protocol, subdomain, domain_name, domain_extension, path)
);

CREATE TABLE [IF NOT EXISTS] outgoing_link (
    protocol TEXT NOT NULL,
    subdomain TEXT NOT NULL,
    domain_name TEXT NOT NULL,
    domain_extension TEXT NOT NULL,
    path TEXT NOT NULL, 


    PRIMARY KEY (protocol, subdomain, domain_name, domain_extension, path),

    FOREIGN KEY (domain_name, domain_extension) REFERENCES domain (name, extension),
    FOREIGN KEY (protocol, subdomain, domain_name, domain_extension, path) REFERENCES url (protocol, subdomain, domain_name, domain_extension, path)
);

CREATE TABLE [IF NOT EXISTS] url (
    protocol TEXT NOT NULL,
    subdomain TEXT NOT NULL,
    domain_name TEXT NOT NULL,
    domain_extension TEXT NOT NULL,
    path TEXT NOT NULL, 

    PRIMARY KEY (protocol, subdomain, domain_name, domain_extension, path),
);