Feature: Command with no arguments

  Scenario: Output when command is run without args
    When I successfully run `scmpuff`
    Then the exit status should be 0
      And the output should contain:
        """
        scmpuff extends common git commands with numeric filename shortcuts.
        """
      And the output should contain "Usage:"
      And the output should contain "Available Commands:"
