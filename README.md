##  A simple load balancer server written in GO 

### Features
- [ ] Protocol type
  - [x] TCP (HTTP/ HTTPS/ TCP over SSL)
  - [ ] UDP
- [x] Load balance algorithms
  - [x] Random 
  - [x] Round-robin
  - [x] IP-hash 

### Config
#### listeners
| Field | Type | Requirement | Description |
| -- | -- | -- | -- |
| port | number | required | |
| upstream | number| required | upstream port |
| health_check_interval | number| optional | seconds ; default : 30 seconds|
| ssl | boolean| optional | enable ssl ; default : false  |
| ssl_certificate | string | optional| |
| ssl_certificate_key | string | optional| |
| algo| string | optional | rand, round-robin, ip-hash ;default : rand |
| nagle | boolean | optional | https://networkencyclopedia.com/nagles-algorithm/ default : false |

#### instances

| Field | Type | Requirement | Description |
| -- | -- | -- | -- |
| addr | string | required | address of upstream server |

gobalancer.json
```json
{
    "listeners":[
        {
            "nagle":false,
            "healthcheck":10,
            "algo":"round-robin",
            "port":3001,
            "ssl":true,
            "ssl_certificate":"./server.crt",
            "ssl_certificate_key":"./server.key",
            "upstream":4001   
        },
        {
            "port":3002,
            "upstream":4002   
        },
        {
            "port":3003,
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