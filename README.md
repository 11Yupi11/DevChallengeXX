# DEV-CHALLENGE 2023


## Docker instructions
#### to run the service:
```bash
docker-compose build
```
```bash
docker-compose up -d app
```

#### to run tests:
```bash
docker-compose exec app make test
```

## Local run instructions
#### to run the service:
```bash
make run
```
#### to run tests:
```bash
make test
```

## Not covered cases
```
From the task, {sheet_id} and {cell_id} must be URL compatible.
Nevertheless, numbers are URL compatible I do not recommend use just a number as a {cell_id}.

Example:
first request POST /api/v1/123/123 with {"value":"1"}  - will return result = 1
second request POST /api/v1/123/124 with {"value":"=123+1"} - will return result = 124, but not result = 2

The same situation with operators "+", "-", "*", "/"
Example:
first request POST /api/v1/123/a+b with {"value":"1"}  - will return result = 1
second request POST /api/v1/123/124 with {"value":"=a+b+1"} - will return result = ERROR, as it won't found "a" and "b"

So to prevent inconsistency do not use as a {cell_id} parameter any kind of operators, numbers or concise numbers like 2e3
```
```
In current implementation, {sheet_id} and {cell_id} are restricted to be no longer than 255 signs long
```

## Short choice description

```
As a persistent storage, SQLite has been chosen.

Advantages:
* no need to create additional instance of DB and connection to outer application
* supports indexing of data to make GET requests work faster
* supports transactions to save consistency 

Disadvantages:
* may make docker-container long time to up in case lots of data stored
```

## Covered cases

```
* all calculations are done in range from min(float64) to max(float64)
* if the result is over the max or min amount, you will receive +Inf or -Inf
* devision by 0, doubled operations like ("++", "--", "**", "//"), will return an error
* if you have expression, like "=-b" and b=-2 it will be counted correctly
* an error will be returned if cell linked for calculations to itself to prevent endless loop
* normal flow if {cell_id} parameters one contain part of another, i.e. "par" and "param"
* all params are saved and shown in lowercase, though you can use uppercase
```

## Ways to improve 
```
Case Handling: Add functionality to handle {cell_id} passed like a number and containing operators "+", "-", "/", "*".

Decrease Complexity Level: The complexity of 'eval' function for now is 16. It is not critical, but it would be good to 
refactor this function to make comlexity less.

Enhance Storage Options: While SQLite is excellent for lightweight and standalone applications, scalability might be
a concern for larger datasets or high concurrent access. Evaluating and offering support for more scalable database 
systems could cater to a wider range of use-cases.

Improve Error Messaging: For all error scenarios, including division by zero or malformed operations, ensure that the 
API returns descriptive error messages to guide the user on the correct path, rather than just indicating an error occurred.

Documentation and Examples: As with all APIs, continuous improvement in documentation, including more diverse examples, 
troubleshooting guides, and best practices.
```