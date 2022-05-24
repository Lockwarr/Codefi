#TODO implementation
Feature: Internal External Links Extraction
    
    Scenario: Successful extraction
        Given the links API is up and running
        When I send a "POST" request to "/api/v1/links" endpoint
        And I have attached a correct expected file to the request
        Then I receive "statusOK"
        And I receive "expected results"

    Scenario: Successful retrieval of batch
        Given the links API is up and running
        And I send a "POST" request to "/api/v1/links" endpoint
        And I have attached a correct expected file to the request
        And I receive "statusOK"
        And I receive "expected results"
        When I send a "GET" request to "/api/v1/links/batchID" endpoint
        Then I receive "statusOK"
        And I receive "expected batch results"