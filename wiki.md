# System Design

- [System Design](#system-design)
  - [Summary](#summary)
  - [Assumption](#assumption)
  - [Feature expectations](#feature-expectations)
  - [Estimations](#estimations)
  - [High Level Design](#high-level-design)
    - [Solution A - Brute force (No additional entry)](#solution-a---brute-force-no-additional-entry)
      - [Steps](#steps)
      - [Analysis](#analysis)
    - [Solution B - Additional entries](#solution-b---additional-entries)
      - [Steps](#steps-1)
      - [Analysis](#analysis-1)
    - [Solution C - Pre-process](#solution-c---pre-process)
      - [Steps](#steps-2)
      - [Analysis](#analysis-2)
  - [Architecture / Detail Design](#architecture--detail-design)
    - [Why serverless?](#why-serverless)
    - [Why golang?](#why-golang)
    - [Database Design](#database-design)
    - [Error handling](#error-handling)
    - [Testing](#testing)

## Summary

This is the system design of Line interview project. For development / deployment / testing details, check the README.md.

All estimations are made based on real product. Any assumptions that made in order to simplify the problem but contradicts with the estimation will be discussed.

I choose [Solution B](#solution-b---additional-entries) as the final design.

## Assumption

1. The number of records that fits the criteria in DB won't break the size limit of RESTful API response.
2. The number of records that fits the criteria in DB doesn't require batch operation (In other word, all required information can be completed within one transaction/query/scan).

## Feature expectations

The api is expected to be called once per day. Based on the current date, it will return a list of people who is having a birthday.

As required, we have the birthday-greeting API with 6 different versions:
* Version 1: Simple Message
* Version 2: Tailer-made Message for different gender
* Version 3: Message with an Elder Picture for those whose age is over 49.
* Version 4: Simple Message with full name
* Version 5: Simple Message but database changes. Please choose another way to persist your member data
* Version 6: Simple Message but different output data format. Please choose another output data format

## Estimations

1. Throughput: pretty low, once per day
2. Latency expectation: the lower the better
3. Read/Write ratio: NO writing to DB. **However, depeneding on the design, we'll do either Query of Scan of whole DB.**
4. Storage estimation: 
   * Line has 84m active users in 2020. 
   * Divided by 365, we has 230k users per day.
   * Estimated size of each record: 200 byte
   * DB size: 16.8G
   * Output size: 63.39MB

Here we'll have few **callouts** for issues that may have on real product:
1. BatchQuery / BatchScan will be needed: DynamoDB has 1MB limit on Query/Scan response size. Based on the above estimation, we'll hit the limit for both Query/Scan request. Last evaluated key will be needed. See more at [AWS DynamoDB documentation](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html#DDB-Query-request-Limit).
2. Size of API response may break limit: AWS API Gateway has only 1MB size limit on API response. [The feasible approach](solution-c---pre-process) is to store the json file (can be partitioned into several files) in cloud storage and return the URL.

## High Level Design

Given examples are shown below.  

```json
{
    "dateOfBirth": "1985/8/8",
    "email": "robert.yen@linecorp.com",
    "firstName": "Robert",
    "gender": "M",
    "id": "1",
    "lastName": "Yen"
}
```

To meet the design goals, **we'll need to parse the dateOfBirth entry and find those people having birthday today.**

Here's some options, which involves with different **database schema**:
### Solution A - Brute force (No additional entry)
#### Steps
   1. Scan the DB when the api is called.
   2. Parse the dateOfBirth.
   3. Filter those people who has birthday today.
   4. Return the result.
#### Analysis
* Pros
  * No additional change on DB records
* Cons
  * Expensive Scan operation
  * Complicate logic
  * Longer response time
### Solution B - Additional entries
#### Steps
   1. Add additional entry as birthYear, birthMonth, and birthDate, which includes only year, month, and date. 
   2. When the api is called, make a DB query and DB should give us only those who has birthday today.
   3. Return the result.
#### Analysis
* Pros
  * Simple logic
* Cons
  * Additional entries per records
### Solution C - Pre-process
#### Steps
   1. Similar to solution 1 and 2, but the Scan/Query on DB will be scheduled and made before the api is called (e.g. Cron job). If use Solution B, additional entries can be added with Cron job.
   2. Store the list in cloud storage.
   3. When the api is called, get the list from cloud storage and return the result.

#### Analysis
* Pros
  * Shortest response time
  * Highly scalable
* Cons
  * Require additional storage/cache for each day.
  * Require mroe json parsing, which will be needed and tailored for six different versions each.

Considering the real-world product, Solution C is optimal. However, we choose Solution B to save some works that needed for six different versions. It can also act as the preliminary stage of Solution C.

---

## Architecture / Detail Design

The project is using infrastructure-as-code philosophy and its architecture is defined in CloudFormation (template.yaml). A custom domain name can be added in front of API gateway.

### Why serverless?
- infrastructure as code
- scalability
- cost
- dev speed

### Why golang?
- performance
- other candidates (Python, NodeJS, Java)

### Database Design

User table: We choose userId as hash key, and use birthYear, birthMonth, and birthDate as secondary index.

### Error handling

We can simply check the http status to know the result of API call. For example, when we're calling CreateUser api:

```
200: success
400: invalid input, which is rarely happening in our project
500: internal error 
```

See more at unit tests so that all cases for each APIs are covered.

### Testing

To follow the CI/CD, a beta/gamma stage will be needed. In this project since we don't have other service as dependencys, simply implement the beta stage will be enough. From service perspective, the beta stage can have:

- testing events: as integration test, a failed event then the deployment will stop the deployment
- different configuration: add stage-related logic within cloudformation template or on new account
- different dependencies: testing against either on different dynamodb table or local dynamodb instance