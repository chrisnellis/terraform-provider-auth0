---
page_title: "Resource: auth0_prompt_screen_partial"
description: |-
  With this resource, you can manage a customized sign up and login experience by adding custom content, form elements and css/javascript. You can read more about this here https://auth0.com/docs/customize/universal-login-pages/customize-signup-and-login-prompts.
---

# Resource: auth0_prompt_screen_partial

With this resource, you can manage a customized sign up and login experience by adding custom content, form elements and css/javascript. You can read more about this [here](https://auth0.com/docs/customize/universal-login-pages/customize-signup-and-login-prompts).

!> This resource appends a specific prompt screen to the list of prompt screens displayed to the user during the authentication flow.
 In contrast, the `auth0_prompt_screen_partials` resource manages the complete set of prompt screens that are displayed during the
 authentication flow. To avoid potential issues, it is recommended not to use this resource in conjunction with the
 `auth0_prompt_screen_partials` resource when managing prompt screens for the same prompt.

## Example Usage

```terraform
resource "auth0_prompt_screen_partial" "login" {
  prompt_type = "login"
  screen_name = "login"
  insertion_points {
    form_content_start = "<div>Form Content Start</div>"
    form_content_end   = "<div>Form Content End</div>"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `prompt_type` (String) The prompt that you are adding partials for. Options are: `login-id`, `login`, `login-password`, `signup`, `signup-id`, `signup-password`, `login-passwordless`.
- `screen_name` (String) The name of the screen associated with the partials

### Optional

- `insertion_points` (Block List, Max: 1) The insertion points for the partials. (see [below for nested schema](#nestedblock--insertion_points))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--insertion_points"></a>
### Nested Schema for `insertion_points`

Optional:

- `form_content_end` (String) Content that goes at the end of the form.
- `form_content_start` (String) Content that goes at the start of the form.
- `form_footer_end` (String) Footer content for the end of the footer.
- `form_footer_start` (String) Footer content for the start of the footer.
- `secondary_actions_end` (String) Actions that go at the end of secondary actions.
- `secondary_actions_start` (String) Actions that go at the start of secondary actions.

## Import

Import is supported using the following syntax:

```shell
# This resource can be imported using the prompt name and screen_name.
#
# As this is not a resource identifiable by an ID within the Auth0 Management API,
# login can be imported using the prompt name and screen name using the format:
# prompt_name:screen_name
#
# Example:
terraform import auth0_prompt_screen_partial.login "login:login"
```