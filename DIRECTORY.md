The directory structure is as follows
```
.
├── DIRECTORY.md # Current document
├── LICENSE # License statement
├── README.md # Readme document
├── cmd # commands
│   ├── deploy # use cli command to deploy kubernetes cluster
│       └── main.go
│   └── portal # use web interface to deploy kubernetes cluster
│       └── main.go
└── pkg
    ├── deploy # deploy command code
    ├── restful # web restful interface code
    ├── result # deploy result and logs
    ├── task # task manager
    ├── utils # utils
    └── web # web interface frontend
```