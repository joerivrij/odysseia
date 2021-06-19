Feature: sokrates
  In order to test odysseia
  As a developer
  We need to be able to validate the functioning of the Sokrates api

  @sokrates
  Scenario Outline: A user can create a new question
    Given the "<service>" is running
    When a new question is requested with category "<category>" and chapter "<chapter>"
    Then the responseCode should be "<response>"
    Examples:
      | service  | response | category | chapter |
      | sokrates | 200      | nomina   | 1       |
      | sokrates | 200      | verba    | 2       |
      | sokrates | 200      | misc     | 5       |
