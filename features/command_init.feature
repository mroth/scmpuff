Feature: init command

  Background:
    Given a mocked home directory

  Scenario: init -s should contain status shortcuts
    When I successfully run `scmpuff init -s`
    Then the output should contain "scmpuff_status_shortcuts()"

  Scenario Outline: --aliases should control whether short aliases are in output
    When I successfully run `scmpuff <subcommand>`
    Then the output <should?> contain "alias gs='scmpuff_status_shortcuts'"
    And  the output <should?> contain "alias ga='scmpuff_add_shortcuts'"
    Examples:
      | subcommand              | should?    |
      | init -s                 | should     |
      | init -as                | should     |
      | init -a -s              | should     |
      | init --aliases=true -s  | should     |
      | init --aliases=false -s | should not |

  Scenario: init -w -s should add git wrapping to output
    Then PENDING

  Scenario Outline: Evaling init -s defines status shortcuts in environment
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
