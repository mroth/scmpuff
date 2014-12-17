Feature: init command

  Background:
    Given a mocked home directory

  @wip
  Scenario Outline: Evaling init -s sets up status shortcuts in environment
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -s)"`
      And I type "type scmpuff_status_shortcuts"
      And I type "type scmpuff_clear_vars"
      And I type "exit"
    Then the output should not contain "not found"
    Examples:
      | shell |
      | bash  |
      | zsh   |

  @wip
  Scenario: init -s should contain status shortcuts
    Then PENDING
  @wip
  Scenario: init -a -s should add aliases to output
    Then PENDING
  @wip
  Scenario: init -w -s should add git wrapping to output
    Then PENDING
