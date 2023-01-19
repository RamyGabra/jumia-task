# jumia-task
This repository contains a solution to Jumia coding challenge that depends on [MySQL](https://www.mysql.com/). <br />
The service manages the stock of different products. 
it contains several functions like:
1. Get Product Stock by SKU
2. Consume stock from specific product
3. batch update stock for different products by sending csv file containing stock changes

## Clone The Project & Run
use [git](https://git-scm.com/) to clone the project
```bash
git clone https://github.com/RamyGabra/jumia-task.git
cd jumia-task
```
use _Make_ to start up the service
```bash
make run
```

OR

use _docker-compose_ to start up the service
```bash
	sudo docker-compose -f docker_compose.yml build
	sudo docker-compose -f docker_compose.yml up
```

## Service Endpoints
The service has two endpoints that require authentication to POST & GET events. 

### Fetch stock endoint `GET /stock`:
the endpoint accepts two params: `sku` & `country`
  ```bash
curl --location --request GET 'http://localhost:3333/stock?country=ci&sku=02ed82eaa783'
  ```
  
### Consume stock endpoint `PUT /stock`
the endpoint accepts two params: `sku` & `country`.
the endpoint will decrement the stock of the product by 1 if there is enough stock.
  ```bash
curl --location --request PUT 'http://localhost:3333/stock?country=ci&sku=02ed82eaa783'
  ```
### Batch update endpoint `POST /stock`
the endpoint accepts a csv file and updates the stock of all products with _stock changes_ in the file.
  ```bash
  curl --location --request POST 'http://localhost:3333/stock' \
--form 'file=@"<CSV file directory>"'
  ```
## Database Model

The `Product Stock` database model uses both `sku` & `country` fields to construct a composite primary key. 
```go
type ProductStock struct {
	Country string `gorm:"primaryKey; not null; autoIncrement:false" json:"country"`
	SKU     string `gorm:"primaryKey; not null; autoIncrement:false" json:"sku"`
	Name    string `json:"name"`
	Stock   int    `gorm:"default:0" json:"stock"`
}

```
<br/>
default value for stick is `0`

## Design Decisions

### Batch Update
**Bulk insertions** are much faster than executing a single **UPDATE** statement for each individual record. It reduces connection overhead and delays index writing until many records have been inserted. check this [link](https://dev.mysql.com/doc/refman/5.7/en/insert-optimization.html) for more details. 

### Leveraging concurrency of Go routines in Batch update 
In the [batch update](https://github.com/RamyGabra/jumia-task/blob/019d1792d450fad75cea6ae7e32c3e59cb2f698f/repo/repo.go#L58) in `repo.go` the function splits up all of the records from the csv into slices. This is done to update the records in the database concurrently speeding up the process. the `sliceSize` can be changed to the optimum value after experimenting and trying out different values.

### Leveraging database transactions in Batch update
For the same file all of the database queries are executed within the same transaction block to be able to rollback if an error occurs in one of the queries.
