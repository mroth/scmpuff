Feature: init command

  Background:
    Given a mocked home directory

  Scenario Outline: Evaling init -s defines status shortcuts in environment
    When I run `<shell>` interactively
      And I initialize scmpuff in `<shell>`
      And I type "type scmpuff_status"
      And I type "type scmpuff_clear_vars"
      And I close the shell `<shell>`
    Then the output should not contain "not found"
    Examples:
      | shell |
      | bash  |
      | zsh   |
      | fish  |

