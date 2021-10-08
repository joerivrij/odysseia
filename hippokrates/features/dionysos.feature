Feature: Dionysos
  In order to test odysseia
  As a developer
  We need to be able to validate the functioning of the Dionysos api

  @dionysos
  Scenario Outline: A user gets a response when searching for a declension
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the responseCode should be "<response>"
    Examples:
      | service  | response |  word |
      | dionysos | 200      |  μάχη |
