##  A simple tcp load balancer server written in GO 

### algo
- rand
- round-robin
- ip-hash

### listeners
| Field | Type | Requirement | Description |
| -- | -- | -- | -- |
| port | number | required | |
| ssl | boolean| required | enable ssl |
| upstream | number| required | upstream port |
| healthcheck | number| optional | seconds ; default : 30 seconds|
| ssl_certificate | string | optional| |
| ssl_certificate_key | string | optional| |
| algo| string | optional | rand, round-robin, ip-hash ;default : rand |

### instances

| Field | Type | Requirement | Description |
| -- | -- | -- | -- |
| addr | string | required | address of upstream server |

gobalancer.json
```json
{
    "listeners":[
        {
            "healthcheck":10,
            "algo":"round-robin",
            "port":3001,
            "ssl":false,
            "ssl_certificate":"./server.crt",
            "ssl_certificate_key":"./server.key",
            "upstream":4001   
        },
        {
            "port":3002,
            "ssl":false,
            "upstream":4002   
        },
        {
            "port":3003,
            "ssl":false,
            "upstream":4003,
            "algo":"ip-hash"
        }
    ],
    "instances":[
        {
            "addr":"127.0.0.1"
        },
        {
            "addr":"127.0.0.2"
        }
    ]
}
```