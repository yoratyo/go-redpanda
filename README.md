# Go-RedPanda

This is project to demonstrate 2 binary as HTTP REST server and Kafka event consumer that communicate via Redpanda and MariaDB. The consumer will do upsert cryptocurrency data to database.

## Step to run
1. Running docker-compose to up Redpanda and MariaDB container 
    ```
    docker compose up --build -d
    ```

2. Run DB migration, to initialize database 
    ```
    make migrate
    ```

3. Build project to be executeable
    ```
    make build
    ```

4. Run Consumer instance
    ```
    make run-consumer
    ```

5. Run REST server instance
    ```
    make run-server
    ```

## API Docs

#### Manage cryptocurrency data

<details>
 <summary><code>POST</code> <code><b>/crypto</b></code> <code>(post cryptocurrency data to Redpanda)</code></summary>

##### Payload (JSON)

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | code      |  required | string   |   |
> | name      |  required | string   |   |
> | category  |  optional | string   |   |
> | algorithm |  optional | string   |   |
> | platform  |  optional | string   |   |
> | industry  |  optional | string   |   |
> | types     |  required | string   |   |
> | mineable  |  required | boolean  |   |
> | audited   |  required | boolean  |   |
> | price     |  required | float    |   |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`        | `{"code":"ETH","name":"Ethereum","category":"business","algorithm":"SHA","platform":"","industry":"","types":"Coin","mineable":true,"audited":true,"price":12456}`                                |
> | `400`         | `application/json`                | `{"error":"Error validate payload: Key: 'CryptoDTO.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`                            |


##### Example cURL

> ```javascript
>  curl --location 'localhost:8090/crypto' \
>--header 'Content-Type: application/json' \
>--data '{
>    "code":"ETH",
>    "name":"Ethereum",
>    "category": "business",
>    "types": "Coin",
>    "algorithm": "SHA",
>    "mineable": true,
>    "audited": true,
>    "price": 12456
>}'
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/crypto</b></code> <code>(get list cryptocurrency data from database)</code></summary>

##### Query Parameter

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | page      |  optional | integer  |   |
> | pageSize  |  optional | integer  |   |
> | code      |  optional | string   |   |
> | name      |  optional | string   |   |
> | category  |  optional | string   |   |
> | algorithm |  optional | string   |   |
> | platform  |  optional | string   |   |
> | industry  |  optional | string   |   |
> | types     |  optional | string   |   |
> | mineable  |  optional | boolean  |   |
> | audited   |  optional | boolean  |   |
> | priceMin  |  optional | float    |   |
> | priceMax  |  optional | float    |   |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`        | `{"page":1,"pageSize":10,"totalPages":1,"totalItems":1,"cryptos":[{"code":"ETH","name":"Ethereum","category":"business","algorithm":"SHA","platform":"","industry":"","types":"Coin","mineable":true,"audited":true,"price":12456}]}`                                |
> | `400`         | `application/json`                | `{"error":"Error decode query:"}`                            |


##### Example cURL

> ```javascript
>  curl --location 'localhost:8090/crypto?mineable=true'
> ```

</details>

------------------------------------------------------------------------------------------