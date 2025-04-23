Feature: Package Management
    As a photographer,
    I can manage packages by updating or deleting them
    so that my package information is accurate and up-to-date.

    Background: Server is running
        Given the server is running
        And a photographer is logged in
        And a photographer has a package and sub package

    Scenario: Photographer updates a package
        When the photographer updates the package details with the following data:
            | title | type          | photos                                                                                      |
            | A     | WEDDING_BLISS | data:image/jpeg;base64,/9j/4AAQSkZJRg==, data:image/jpeg;base64,/9j/4AAQSkZJRg==             |
        Then the package information is updated with following data:
            | title | type          | photos                                                                                      |
            | A     | WEDDING_BLISS | data:image/jpeg;base64,/9j/4AAQSkZJRg==, data:image/jpeg;base64,/9j/4AAQSkZJRg==             |
    
    Scenario: Photographer updates a package
        When the photographer updates the package details with the following data:
            | title         | type | photos |
            | Birthday 2025 | __omit__ | __omit__   |
        Then the package information is updated with following data:
            | title         | type | photos |
            | Birthday 2025 | __omit__ | __omit__   |
    
    Scenario: Photographer updates a package
        When the photographer updates the package details with the following data:
            | title         | type | photos |
            | Maternity Glow - 2025! | __omit__ | __omit__   |
        Then the package information is updated with following data:
            | title         | type | photos |
            | Maternity Glow - 2025! | __omit__ | __omit__   |
    
    
    Scenario: Photographer deletes a package
        When the photographer deletes the package
        Then the package is removed
