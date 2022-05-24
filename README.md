# Codefi

## Windows
The code has been tested with Windows OS only

## Starting the rest api
1. Install Chocolatey so you can use `make` on Windows OS - https://community.chocolatey.org/courses/installation/installing?method=installing-chocolatey

2. You can start this service with the command:

make build_and_run

## Testing
1. Run all tests:

make test-with-component

2. Run only unit tests(without integration tests):

make test

## Using the rest api
This application has one service.
There are 2 REST API endpoints for this service:

1. `/api/v1/links`
POST endpoint expecting content-type set to form-data with key name `urlsFile` and value the attached file. The file should be consisting of multi-line text, a valid url on each line
### Example:
there is an file in the `services\links\component-tests\testdata` folder

### Results example:
```json
{
    "data": {
        "Results": [
            {
                "id": "241edc85-6221-42c4-abd4-24eb1fb3261d",
                "batch_id": "b2fe8be7-902d-4211-bf55-f3119a282986",
                "page_url": "https://www.google.com/",
                "internal_links_num": 6,
                "external_links_num": 13,
                "success": true,
                "error": null,
                "created_at": "2022-05-23T10:51:01.5371587Z",
                "updated_at": "2022-05-23T10:51:01.5371587Z"
            },
            {
                "id": "428c2cea-30a4-43eb-8230-de7efcb82132",
                "batch_id": "b2fe8be7-902d-4211-bf55-f3119a282986",
                "page_url": "https://www.facebook.com",
                "internal_links_num": 27,
                "external_links_num": 20,
                "success": true,
                "error": null,
                "created_at": "2022-05-23T10:51:01.5371587Z",
                "updated_at": "2022-05-23T10:51:01.5371587Z"
            }
        ]
    }
}
```


2. `/api/v1/links/{batch_id}`
GET endpoint for listing all links for given batch_id where batch_id is id shared between urls which were processed at once

### Results example:
```json
{
    "data": {
        "Results": [
            {
                "id": "241edc85-6221-42c4-abd4-24eb1fb3261d",
                "batch_id": "b2fe8be7-902d-4211-bf55-f3119a282986",
                "page_url": "https://www.google.com/",
                "internal_links_num": 6,
                "external_links_num": 13,
                "success": true,
                "error": null,
                "created_at": "2022-05-23T10:51:01.5371587Z",
                "updated_at": "2022-05-23T10:51:01.5371587Z"
            },
            {
                "id": "428c2cea-30a4-43eb-8230-de7efcb82132",
                "batch_id": "b2fe8be7-902d-4211-bf55-f3119a282986",
                "page_url": "https://www.facebook.com",
                "internal_links_num": 27,
                "external_links_num": 20,
                "success": true,
                "error": null,
                "created_at": "2022-05-23T10:51:01.5371587Z",
                "updated_at": "2022-05-23T10:51:01.5371587Z"
            }
        ]
    }
}
```
