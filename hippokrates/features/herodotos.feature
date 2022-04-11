Feature: Herodotos
  In order to work with sentences
  As a greek enthusiast
  We need to be able to validate the functioning of the Herodotos api

  @herodotos
  Scenario Outline: Querying authors should return a list of authors
    Given the "<service>" is running
    When a query is made for all authors
    Then the author "<author>" should be included
    And the number of authors should exceed "<results>"
    Examples:
      | service    | author     | results |
      | herodotos | herodotos   |   4      |
      | herodotos | plato   |   4      |

  @herodotos
  Scenario Outline: Querying books should return a list of books
    Given the "<service>" is running
    When a query is made for all books by author "<author>"
    Then the book "<book>" should be included
    Examples:
      | service    | author     | book |
      | herodotos | thucydides   | 1    |
      | herodotos | ploutarchos   | 1        |

  @herodotos
  Scenario Outline: A client can create a question with a new sentence
    Given the "<service>" is running
    When an author and book combination is queried
    Then the sentenceId should be longer than "<length>"
    And the sentence should include non-ASCII (Greek) characters
    Examples:
      | service    | length |
      | herodotos | 12 |

  @herodotos
  Scenario Outline: A client can return a question with an answer
    Given the "<service>" is running
    When an author and book combination is queried
    Then a translation is returned
    And a correctness percentage
    And a sentence with a translation
    Examples:
      | service    |
      | herodotos |