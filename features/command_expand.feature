Feature: command expansion at command line

  Background:
    Given I override the environment variables to:
      | variable | value |
      | e1       | a.txt |

  Scenario: Dont interfere with CLI "options" that are passed along after `--`
    When I successfully run `scmpuff expand -- git foo -x 1`
    Then the output should match /git\tfoo\t-x\ta.txt/

  Scenario: Semicolons in commit messages
    Given a git repository named "whatever"
      And I cd to "whatever"
      And a 4 byte file named "a.txt"
      And I successfully run the following commands:
        | git add a.txt                               |
      When I successfully run `scmpuff expand -- git commit -m "foo; bar"`
      Then the stderr should not contain anything
        And the output should match /git\tcommit\t-m\tfoo\\;\\ bar/

  Scenario: Dont expand files or directories with numeric names
    Given an empty file named "1"
    Given a directory named "2"
    When I successfully run `scmpuff expand 3 2 1`
      Then the output should contain "c.txt"
      But the output should not contain "b.txt"
      And the output should not contain "a.txt"
