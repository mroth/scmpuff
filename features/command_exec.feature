Feature: command expansion at command line

  Background:
    Given I am in a git repository
    And an empty file named "a.txt"
    And an empty file named "b.txt"
    And I override the environment variables to:
      | variable | value |
      | e1       | a.txt |
      | e2       | b.txt |

  Scenario: Expand single digit case
    When I successfully run `scmpuff exec -- git add 2`
    And I successfully run `git status -s a.txt b.txt`
    Then the stdout should contain exactly "A  b.txt\n?? a.txt\n"

  Scenario: Expand multiple digit case
    When I successfully run `scmpuff exec -- git add 1 2`
    And I successfully run `git status -s a.txt b.txt`
    Then the stdout should contain exactly "A  a.txt\nA  b.txt\n"

  Scenario: Expand ranged digit case
    When I successfully run `scmpuff exec -- git add 1-2`
    And I successfully run `git status -s a.txt b.txt`
    Then the stdout should contain exactly "A  a.txt\nA  b.txt\n"
