The directory structure is as follows
```
.
├── DIRECTORY.md # Current document
├── LICENSE # License statement
├── Makefile # Make configuration file
├── README.md # Readme document
├── builds # Generated executable files
│   ├── debug # Local debug
│   ├── docker # Docker build file
│   │   └── Dockerfile # Docker build file
│   └── release # Release directory
├── charts # Helm charts for deploying KPaaS components
│   └── calico # chart for deploying calico networking
├── cli # Commands line interface
│   ├── console # Provide RESTful APIs and manage cluster services
│   │   └── main.go
│   └── docker # Files that need to be used in Docker, such as: startup files
│       └── entrypoint.sh
├── docs # Design documents
└── pkg
    ├── api # RESTful API Controller
    ├── application # Software life cycle code
    ├── config # configuration structure
    ├── model # data model structure
    ├── swaggerdocs # swagger docs
    └── utils # Util codes
```
