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
│   │   └── service # service docker build file
│   └── release # Release directory
├── charts # Helm charts for deploying KPaaS components
│   └── calico # chart for deploying calico networking
├── docs # Design documents
│   ├── manual # Manual files
│   └── ui-design # UI Design Files
├── pkg
│   ├── deploy # Deployment command codes
│   ├── service # Global control service codes
│   │   ├── api # RESTful API Controller
│   │   ├── application # Software life cycle codes
│   │   ├── config # configuration structure
│   │   ├── model # data model structure
│   │   └── swaggerdocs # swagger docs
│   └── utils # Util codes
├── run # application main entrypoints
│   ├── deploy # Kubernetes deployment service
│   │   └── main.go
│   ├── docker # Files that need to be used in Docker, such as: startup files
│   │   └── entrypoint.sh # docker entry point file
│   ├── portal # User web interface
│   │   └── main.go
│   └── service # Provide RESTful APIs and manage cluster services
│       └── main.go
└── sites # Web interface codes
```
