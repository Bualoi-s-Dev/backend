Feature: Appointment Management US7
    As a customer or photographer,
    I can manage appointments by creating, completing, or canceling them
    so that my schedule is organized and up-to-date.

    Background: Server is running
        Given the server is running

    Scenario: Customer creates an appointment
        Given a customer is logged in
        And a sub package is selected
        And a time slot is selected
        When the customer submits the appointment request
        Then the appointment is created
        And the customer’s schedule is updated

    Scenario: Photographer marks an appointment as completed
        Given a photographer is logged in
        And an appointment is scheduled
        When the photographer marks the appointment as completed
        Then the appointment status is updated to completed
        And the photographer’s schedule is updated

    Scenario: Photographer cancels an appointment
        Given a photographer is logged in
        And an appointment is scheduled
        When the photographer cancels the appointment
        Then the appointment is canceled
        And the photographer’s schedule is updated

    Scenario: Customer cancels an appointment
        Given a customer is logged in
        And an appointment is scheduled
        When the customer cancels the appointment
        Then the appointment is canceled
        And the customer’s schedule is updated