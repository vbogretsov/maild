# maild
Simple notification service for micro service architecture

## Installation

## Usage

## Why it was decided to use files and not database

Email microservice

Where to store templates

Scores:

1. Code size 3
2. Dependencies 5
3. Blue green deployment 2
4. New image building 1
5. Unsolved questions 5

Disk

    Pros:
        1. Less code +3
        2. Less dependencies +5

    Cons:
        1. Requires building image with templates -1
        2. Blue green deployment impossible if just templates added -2

    Sum: 5

Databse

    Pros:
        1. Not required building image with templates +1
        2. Blue green deployment not actual if just templates added +2

    Cons:
        1. More code -3
        2. More dependencies -5
        3. Who inserts templates
            3.1 email service
                3.1.1 More code -3
            3.2 site service
                3.1.2 More dependencies -5

    Sum: -8 or -10

## Licence

See the LICENCE file.
