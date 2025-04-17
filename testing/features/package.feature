Feature: Package Management
    As a photographer,
    I can manage packages by updating or deleting them
    so that my package information is accurate and up-to-date.

    Background: Server is running
        Given the server is running
        And a photographer has a package and sub package
        And a photographer is logged in

    Scenario: Photographer updates a package
        When the photographer updates the package details
        Then the package information is updated

    Scenario: Photographer deletes a package
        When the photographer deletes the package
        Then the package is removed
