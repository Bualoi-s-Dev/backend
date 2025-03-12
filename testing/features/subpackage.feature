Feature: EPIC6 Package system
    As a photographer,
    I can create, update, and delete subpackages under a specific package,
    that I can customize my services.

    Background: Server is running
        Given the server is running
        And the photographer is logged in

    Scenario: Photographer creates a subpackage
        Given the photographer has a package
        When the photographer creates a subpackage
        Then the subpackage is created and added to the package

    Scenario: Photographer updates a subpackage
        When the photographer updates a subpackage
        Then the subpackage is updated

    Scenario: Photographer deletes a subpackage
        When the photographer deletes a subpackage
        Then the subpackage is deleted

    Scenario: Photographer cannot creates a subpackage
        When the photographer creates a subpackage with wrong format
        Then the subpackage is not created and not added to the package

    Scenario: Photographer cannot updates a subpackage
        When the photographer updates a subpackage with wrong format
        Then the subpackage is not updated

    Scenario: Photographer deletes a non-existent subpackage
        When the photographer deletes a non-existent subpackage
        Then the subpackage is not deleted