Feature: EPIC7-1 Appointment
    As a customer,
    I can create appointments based on sub package and selected time slots,
    so that my schedule is organized.

    Background: Server is running
        Given the server is running

    Scenario: Customer creates the appointment
        Given a photographer has a package and sub package
        And a customer is logged in
        When a customer creates an appointment
        Then the appointment is created