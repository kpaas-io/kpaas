The directory structure is as follows
```
.
├── DIRECTORY.md # Current document
├── LICENSE # License statement
├── README.md # Readme document
├── docs # design documents
├── cli # commands line interface
│   ├── deploy # use cli command to deploy kubernetes cluster
│       └── main.go
│   └── portal # use web interface to deploy kubernetes cluster
│       └── main.go
└── pkg
│   ├── deploy # deploy command code
│   ├── restful # web restful interface code
│   ├── result # deploy result and logs
│   ├── task # task manager
│   └── utils # utils
└── sites # web interface frontend
    ├── portal # portal frontend codes
    └── docs # manual
```