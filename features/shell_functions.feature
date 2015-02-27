Feature: scmpuff_status_shortcuts() function
  In order to verify the scmpuff_status_shortcuts function work correctly
  I want to make sure they wrap the underlying scmpuff commands properly

  Background:
    Given a mocked home directory

  # we don't want to swallow up error conditions, so make sure we pass them
  # along properly
  @outside-repo
  Scenario Outline: Handles error conditions from underlying scmpuff status cmd
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -s)"`
      And I type "scmpuff_status_shortcuts"
      And I type "exit $?"
    Then the exit status should be 128
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |


  Scenario Outline: Handles normal conditions from underlying scmpuff status cmd
    Given I am in a git repository
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -s)"`
      And I type "scmpuff_status_shortcuts"
      And I type "exit $?"
    Then the exit status should be 0
    And the output should contain "No changes (working directory clean)"
    Examples:
      | shell |
      | bash  |
      | zsh   |
