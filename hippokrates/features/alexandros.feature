Feature: Alexandros
  In order to use the dictionary
  As a greek enthusiast
  We need to be able to validate the functioning of the Alexandros api

  @alexandros
  Scenario Outline: Searching for a word in the dictionary that word should be included in the response
    Given the "<service>" is running
    When the word "<word>" is queried
    Then the word "<word>" should be included in the response
    Examples:
      | service    | word     |
      | alexandros | ἀγαθός   |
      | alexandros | ἡσσάομαι |

  @alexandros
  Scenario Outline: Searching for a word stripped of accents the result should contain an original version of that word
    Given the "<service>" is running
    When the word "<word>" is stripped of accents
    Then the word "<word>" should be included in the response
    Examples:
      | service    | word    |
      | alexandros | ὕδατος  |
      | alexandros | ἰδιώτης |

  @alexandros
  Scenario Outline: Searching for the beginning of a word a response with a full set of words should be returned
    Given the "<service>" is running
    When the partial "<partial>" is queried
    Then the word "<word>" should be included in the response
    Examples:
      | service    | partial | word    |
      | alexandros | αγα     | ἀγαθός  |
      | alexandros | ἱστ     | ἱστορία |

  @alexandros
  Scenario Outline: The maximum number of results is set
    Given the "<service>" is running
    When the word "<word>" is queried
    Then the number of results should not exceed "<results>"
    Examples:
      | service    | word | results |
      | alexandros | α    | 15      |
      | alexandros | ν    | 15      |

  @alexandros
  Scenario Outline: Queries with no results return an error
    Given the "<service>" is running
    When the word "<word>" is queried with an error
    Then an error containing "<message>" is returned
    Examples:
      | service    | word             | message         |
      | alexandros | ναυμαχίαναυμαχία | 200 but got 404 |
    